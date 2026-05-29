# Authentication Service API Specification

The **Authentication Service** is responsible for managing session lifecycles, issuing short-lived JSON Web Tokens (JWT) for stateless access, and rotating long-lived Refresh Tokens for persistent sessions. It registers with Consul under the name `auth-service` and operates on port `8081`. 

All external requests are routed via the API Gateway using the prefix `/api/auth/*`.

---

## 📡 Endpoints Specification

### 🟢 Login (Create Session)

#### Endpoint Description
Registers a new session in the system for a user using their credentials. This endpoint returns both an Access Token (for short-term API calls) and a Refresh Token (to securely obtain new access tokens later).

- **HTTP Method**: `POST`
- **Path**: `/api/auth/sessions` (routed to `POST /sessions` downstream)
- **Authentication**: None (Public)

#### Request Body
| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `email` | `string` | Yes | User email |
| `password` | `string` | Yes | User password |

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "id": "e9c15ad2-f673-455b-86d7-4632b50fe732",
  "user_id": "c3b99db1-d419-48e0-bb15-081079d38bb1",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "access_expired_at": "2026-05-19T11:45:00Z",
  "refresh_expired_at": "2026-05-26T11:30:00Z",
  "created_at": "2026-05-19T11:30:00Z"
}
```

#### Failure Response
- **Status Code**: `401 Unauthorized`
- **Body**:
```json
{
  "error": "Invalid email or password"
}
```

---

### 🟢 Refresh Session (Public)

#### Endpoint Description
Exchanges a valid refresh token for a **brand new** Access Token and a **brand new** Refresh Token. The previous refresh token is immediately invalidated (Token Rotation), preventing replay attacks and ensuring tight security control.

- **HTTP Method**: `POST`
- **Path**: `/api/auth/sessions/refresh` (routed to `POST /sessions/refresh` downstream)
- **Authentication**: None (Public)

#### Request Body
| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `refresh_token` | `string` | Yes | A valid, unexpired refresh token. |

#### Request Example
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "id": "1ab23cd4-e567-8901-f234-567890abcdef",
  "user_id": "c3b99db1-d419-48e0-bb15-081079d38bb1",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.new...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.new...",
  "access_expired_at": "2026-05-19T12:00:00Z",
  "refresh_expired_at": "2026-05-26T11:45:00Z",
  "created_at": "2026-05-19T11:45:00Z"
}
```

#### Failure Response
- **Status Code**: `401 Unauthorized`
- **Body**:
```json
{
  "error": "invalid refresh token: token has expired"
}
```

---

### 🟢 Validate Session (Internal)

#### Endpoint Description
Verifies if an **access token** is active, valid, and not expired. Used by other microservices or the gateway to perform session checks.

- **HTTP Method**: `POST`
- **Path**: `/api/auth/sessions/validate` (routed to `POST /sessions/validate` downstream)
- **Authentication**: None (Public / Internal)

#### Request Query Parameters
| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `token` | `string` | Yes | The active access token to validate. |

#### Request Example
`POST http://localhost:8080/api/auth/sessions/validate?token=eyJhbGciOiJIUzI1...`

#### Success Response
- **Status Code**: `200 OK`
- **Body**: (Returns the active session payload matching the token)

#### Failure Response
- **Status Code**: `401 Unauthorized`
- **Body**:
```json
{
  "error": "session expired or invalid"
}
```

---

### 🟢 Revoke Session (Authenticated)

#### Endpoint Description
Explicitly terminates a session. Used when a user logs out. You can pass either the Access Token or the Refresh Token.

- **HTTP Method**: `DELETE`
- **Path**: `/api/auth/sessions` (routed to `DELETE /sessions` downstream)
- **Authentication**: Bearer JWT

#### Request Headers
| Header | Value | Required | Description |
| :--- | :--- | :--- | :--- |
| `Authorization` | `Bearer <token>` | Yes | The JWT of the authenticated session owner. |

#### Request Query Parameters
| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `token` | `string` | Yes | The access or refresh token to revoke/delete. |

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "message": "Session revoked successfully"
}
```

#### Failure Response
- **Status Code**: `400 Bad Request`
- **Body**:
```json
{
  "error": "Query parameter token is required"
}
```
