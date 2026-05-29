package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"backend-gmao/pkg/auth"
	"github.com/gin-gonic/gin"
)

type ContextKey string

const (
	// Context keys for storing authenticated user information.
	ContextKeyUserID     ContextKey = "auth_user_id"
	ContextKeyEmail      ContextKey = "auth_email"
	ContextKeyFullName   ContextKey = "auth_full_name"
	ContextKeyRole       ContextKey = "auth_role"
	ContextKeyPrivileges ContextKey = "auth_privileges"
)

// RequireAuth returns a Gin middleware that validates the JWT access token
// from the Authorization header and injects claims into the Gin context.
func RequireAuth(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Authorization header is required",
				},
			})
			return
		}

		// Expect "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_TOKEN_FORMAT",
					"message": "Authorization header must be in the format: Bearer <token>",
				},
			})
			return
		}

		claims, err := jwtManager.ValidateAccessToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_TOKEN",
					"message": "Invalid or expired access token",
				},
			})
			return
		}

		// Inject user info into Gin context for downstream Gin middlewares
		c.Set(string(ContextKeyUserID), claims.UserID)
		c.Set(string(ContextKeyEmail), claims.Email)
		c.Set(string(ContextKeyFullName), claims.FullName)
		c.Set(string(ContextKeyRole), claims.Role)
		c.Set(string(ContextKeyPrivileges), claims.Privileges)

		// Inject user info into standard request context for downstream service layer
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyEmail, claims.Email)
		ctx = context.WithValue(ctx, ContextKeyFullName, claims.FullName)
		ctx = context.WithValue(ctx, ContextKeyRole, claims.Role)
		ctx = context.WithValue(ctx, ContextKeyPrivileges, claims.Privileges)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// RequirePrivilege returns a Gin middleware that checks if the authenticated user
// has the specified privilege. Must be used after RequireAuth.
func RequirePrivilege(privilege string) gin.HandlerFunc {
	return func(c *gin.Context) {
		privileges, exists := c.Get(string(ContextKeyPrivileges))
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "Access denied: no privileges found",
				},
			})
			return
		}

		userPrivileges, ok := privileges.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to parse user privileges",
				},
			})
			return
		}

		// SYSTEM_ADMIN has access to everything
		for _, p := range userPrivileges {
			if p == "SYSTEM_ADMIN" || p == privilege {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INSUFFICIENT_PRIVILEGES",
				"message": "You do not have the required privilege: " + privilege,
			},
		})
	}
}

// RequireAnyPrivilege returns a Gin middleware that checks if the authenticated user
// has AT LEAST ONE of the specified privileges. Must be used after RequireAuth.
func RequireAnyPrivilege(allowedPrivileges ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		privileges, exists := c.Get(string(ContextKeyPrivileges))
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "Access denied: no privileges found",
				},
			})
			return
		}

		userPrivileges, ok := privileges.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to parse user privileges",
				},
			})
			return
		}

		// SYSTEM_ADMIN has access to everything, or match any allowed privilege
		for _, up := range userPrivileges {
			if up == "SYSTEM_ADMIN" {
				c.Next()
				return
			}
			for _, ap := range allowedPrivileges {
				if up == ap {
					c.Next()
					return
				}
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INSUFFICIENT_PRIVILEGES",
				"message": "You do not have any of the required privileges to access this resource",
			},
		})
	}
}


// Cors returns a middleware that handles Cross-Origin Resource Sharing (CORS).
// This is essential for allowing the Flutter Web frontend to communicate with the backend.
func Cors() gin.HandlerFunc {
	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsEnv != "" {
		allowedOrigins = strings.Split(allowedOriginsEnv, ",")
		for i, v := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(v)
		}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Determine the allowed origin to return
		allowedOrigin := ""
		if len(allowedOrigins) == 0 {
			// Fallback for development if not configured.
			// For defense in depth, we only allow localhost explicitly.
			// 🛡️ Security: Never reflect arbitrary origins, especially with credentials.
			if strings.HasPrefix(origin, "http://localhost:") || origin == "http://localhost" {
				allowedOrigin = origin
			}
		} else {
			for _, o := range allowedOrigins {
				// 🛡️ Security: Explicitly match origin. If a wildcard is requested, return the literal "*"
				// rather than reflecting the origin. Browsers will inherently block "*" with credentials,
				// which is the desired secure behavior, while still allowing non-credentialed public access.
				if origin == o {
					allowedOrigin = origin
					break
				} else if o == "*" {
					allowedOrigin = "*"
					break
				}
			}
		}

		if allowedOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Gateway-Service, X-Internal-Service")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequireInternalService returns a middleware that ensures the request comes from
// an internal service (via the API gateway or direct inter-service call).
func RequireInternalService() gin.HandlerFunc {
	return func(c *gin.Context) {
		gatewayHeader := c.GetHeader("X-Gateway-Service")
		internalHeader := c.GetHeader("X-Internal-Service")

		if gatewayHeader == "" && internalHeader == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "This endpoint is only accessible internally",
				},
			})
			return
		}

		c.Next()
	}
}
