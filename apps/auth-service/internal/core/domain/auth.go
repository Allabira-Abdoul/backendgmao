package domain

import "github.com/google/uuid"

// Credentials represents a user's login credentials.
type Credentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// TokenPair holds the generated JWT and the Refresh Token.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// InternalUserResponse represents the expected response from the user-service internal API.
type InternalUserResponse struct {
	ID                 uuid.UUID `json:"id"`
	FullName           string    `json:"full_name"`
	Email              string    `json:"email"`
	Password           string    `json:"password"`
	Status             string    `json:"status"`
	MustChangePassword bool      `json:"must_change_password"`
	RoleName           string    `json:"role_name"`
	Privileges         []string  `json:"privileges"`
}
