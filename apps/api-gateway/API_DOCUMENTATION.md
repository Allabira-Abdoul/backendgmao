# API Gateway Service API Specification

The **API Gateway** is the single entry point for all client requests. It provides dynamic reverse-proxying, whitelists allowed downstream services, prevents Server-Side Request Forgery (SSRF) and path-traversal attacks, manages rate-limiting, and strips internal-only headers.

---

## 🔒 Security Architectures & Protections

1. **Global Rate Limiting**: Limit of **100 requests per minute** per client IP. Returns `429 Too Many Requests` when exceeded.
2. **Security Headers**:
   - `X-Content-Type-Options: nosniff` (prevents MIME type sniffing)
   - `X-Frame-Options: DENY` (prevents clickjacking)
   - `X-XSS-Protection: 1; mode=block` (enforces XSS protection)
   - `Strict-Transport-Security` (enforces SSL maximum age of 1 year)
3. **SSRF & Internal Block**: Block attempts to access any endpoint containing `/internal` or `/internal/` in the path, returning a `403 Forbidden` response.
4. **Header Stripping**: Drops the sensitive header `X-Internal-Service` if present in incoming client requests to prevent header injection. Adds downstream audit tracking headers:
   - `X-Forwarded-For`: Original client IP address
   - `X-Forwarded-Host`: Original client request host
   - `X-Gateway-Service`: `api-gateway`

---

## ⚡ Routing Configurations

Requests matching `/api/{service}/*` are proxy-routed dynamically via **Consul Service Discovery**:

| Service Path Prefix | Target Service Registered ID | Allowed Downstream Ports |
| :--- | :--- | :--- |
| `/api/auth/*` | `auth-service` | `8081` |
| `/api/user/*` | `user-service` | `8082` |
| `/api/asset/*` | `asset-service` | `8083` |
| `/api/maintenance/*` | `maintenance-service` | `8084` |
| `/api/prediction/*` | `prediction-service` | `8085` |
| `/api/analytics/*` | `analytics-service` | `8086` |
| `/api/audit/*` | `audit-service` | `8087` |

---

## 📡 Endpoints Specification

### 🟢 Get Gateway Health Status

#### Endpoint Description
Retrieve the operational status of the API Gateway instance itself.

- **HTTP Method**: `GET`
- **Path**: `/health`
- **Authentication**: None (Public)

#### Request Format
- **Headers**: None

#### Success Response
- **Status Code**: `200 OK`
- **Headers**: `Content-Type: application/json`
- **Body**:
```json
{
  "service": "api-gateway",
  "status": "UP"
}
```

---

### 🟡 Proxy Service Routes

#### Endpoint Description
Dynamic catch-all path forwarding to Whitelisted microservices.

- **HTTP Method**: `ANY` (`GET` / `POST` / `PUT` / `DELETE` / `PATCH`)
- **Path**: `/api/:service/*path`
- **Authentication**: Inherited from target endpoints.

#### Request Format
- **Path Parameters**:
  - `service`: `user`, `auth`, `asset`, `maintenance`, `prediction`, `analytics`, `audit`
  - `path`: Downstream service resource path.

#### Failure Response (Service Not Whitelisted)
- **Status Code**: `403 Forbidden`
- **Body**:
```json
{
  "error": "forbidden",
  "message": "Access to this service is not allowed"
}
```

#### Failure Response (Path Traversal / Internal Endpoint Blocked)
- **Status Code**: `403 Forbidden`
- **Body**:
```json
{
  "error": "forbidden",
  "message": "Access to internal endpoints is not allowed"
}
```

#### Failure Response (Consul Service Down / Unreachable)
- **Status Code**: `502 Bad Gateway`
- **Body**:
```json
{
  "error": "service_unavailable",
  "message": "Could not discover service: user-service"
}
```
