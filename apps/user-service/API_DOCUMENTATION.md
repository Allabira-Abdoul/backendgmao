# User & RBAC Service API Specification

The **User & RBAC Service** manages user accounts, user profiles, system roles, teams, and associated capability privileges. It registers with Consul under the name `user-service` and operates on port `8082`.

All external client requests are routed through the API Gateway using the prefix `/api/user/*`.

---

## 🔒 Security & Privileges

This service enforces the privilege-based RBAC model. Users must present a Bearer JWT containing the required privilege for each endpoint:

- **User Read**: Requires `USER_VIEW` (Note: Users can always view their own profile without this privilege)
- **User Create**: Requires `USER_CREATE`
- **User Update**: Requires `USER_UPDATE`
- **User Delete**: Requires `USER_DELETE`
- **Role Read**: Requires `ROLE_VIEW`
- **Role Create**: Requires `ROLE_CREATE`
- **Role Update**: Requires `ROLE_UPDATE`
- **Role Delete**: Requires `ROLE_DELETE`
- **Team Read**: Requires `TEAM_VIEW`
- **Team Create**: Requires `TEAM_CREATE`
- **Team Update**: Requires `TEAM_UPDATE`
- **Team Delete**: Requires `TEAM_DELETE`

---

## 📡 Endpoints Specification

### 🟢 Get Current User Profile (Authenticated)

Retrieve details of the currently logged-in user context.

- **HTTP Method**: `GET`
- **Path**: `/api/user/users/me` (routed to `GET /users/me` downstream)
- **Authentication**: Bearer JWT (Any active logged-in user)

#### Request Format
- **Headers**:
  - `Authorization: Bearer <token>`

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "id": "c3b99db1-d419-48e0-bb15-081079d38bb1",
  "full_name": "Auditor User",
  "email": "auditor@gmao.com",
  "status": "ACTIVE",
  "must_change_password": true,
  "role": {
    "id": "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d",
    "name": "Auditor",
    "description": "System auditor with logs inspection privilege",
    "privileges": [
      "USER_VIEW",
      "ROLE_VIEW",
      "SYSTEM_AUDIT_VIEW"
    ],
    "created_at": "2026-05-19T10:00:00Z",
    "updated_at": "2026-05-19T10:00:00Z"
  },
  "team": {
    "id": "f8a71b2d-3c4e-5f6a-7b8c-9d0e1f2a3b4c",
    "name": "Maintenance Alpha",
    "manager_id": "c3b99db1-d419-48e0-bb15-081079d38bb1",
    "description": "Primary maintenance squad",
    "created_at": "2026-05-19T10:00:00Z",
    "updated_at": "2026-05-19T10:00:00Z"
  },
  "created_at": "2026-05-19T10:15:00Z",
  "updated_at": "2026-05-19T10:15:00Z"
}
```

---

### 🟢 Change Password (Authenticated)

Allows a user to change their own password. This is required if their account has been reset by an administrator and `must_change_password` is true.

- **HTTP Method**: `POST`
- **Path**: `/api/user/users/me/change-password`
- **Authentication**: Bearer JWT

#### Request Body
```json
{
  "new_password": "mySecureNewPassword123"
}
```

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "message": "Password changed successfully"
}
```

---

### 🟢 List Users (Privilege: `USER_VIEW`)

Paginated list of all users in the system.

- **HTTP Method**: `GET`
- **Path**: `/api/user/users`
- **Authentication**: Bearer JWT

#### Request Query Parameters
| Parameter | Type | Required | Default | Description |
| :--- | :--- | :--- | :--- | :--- |
| `page` | `integer` | No | `1` | The page number to fetch. |
| `per_page` | `integer` | No | `20` | The page size (limit). |

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
[
  {
    "id": "c3b99db1-d419-48e0-bb15-081079d38bb1",
    "full_name": "Auditor User",
    "email": "auditor@gmao.com",
    "status": "ACTIVE",
    "must_change_password": false,
    "role": {
      "id": "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d",
      "name": "Auditor",
      "description": "System auditor with logs inspection privilege"
    },
    "team": {
      "id": "f8a71b2d-3c4e-5f6a-7b8c-9d0e1f2a3b4c",
      "name": "Maintenance Alpha",
      "manager_id": "c3b99db1-d419-48e0-bb15-081079d38bb1",
      "description": "Primary maintenance squad"
    },
    "created_at": "2026-05-19T10:15:00Z",
    "updated_at": "2026-05-19T10:15:00Z"
  }
]
```

---

### 🟢 Get Specific User (Privilege: `USER_VIEW` or Self)

Retrieve user record by unique UUID. Users can always retrieve their own record without needing the `USER_VIEW` privilege.

- **HTTP Method**: `GET`
- **Path**: `/api/user/users/:id`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**: (Similar structure to `GET /users/me`)

#### Failure Response (Forbidden)
- **Status Code**: `403 Forbidden`
- **Body**:
```json
{
  "error": "Insufficient privileges to view other users"
}
```

---

### 🟢 Create User (Privilege: `USER_CREATE`)

Register a new user into the system. Note: Users must be assigned exactly one role, and optionally one team.

- **HTTP Method**: `POST`
- **Path**: `/api/user/users`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Validation | Description |
| :--- | :--- | :--- | :--- | :--- |
| `full_name` | `string` | Yes | Min 2, max 255 chars | Full name of the user. |
| `email` | `string` | Yes | Valid email syntax | System login email. |
| `password` | `string` | Yes | Min 8 chars | Plain-text password. |
| `role_id` | `string` | Yes | Valid UUID | Target role UUID. |

##### Example
```json
{
  "full_name": "Technician Joe",
  "email": "joe@gmao.com",
  "password": "securepassword123",
  "role_id": "b2f689e4-cc79-4d6d-85fa-7f8976a45612"
}
```

---

### 🟢 Update User (Privilege: `USER_UPDATE`)

Update an existing user's details, status, role, or team assignment.

- **HTTP Method**: `PUT`
- **Path**: `/api/user/users/:id`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Validation / Allowed Values | Description |
| :--- | :--- | :--- | :--- | :--- |
| `full_name` | `string` | No | Min 2, max 255 chars | Updated full name. |
| `email` | `string` | No | Valid email syntax | Updated login email. |
| `status` | `string` | No | `ACTIVE`, `INACTIVE`, `LOCKED` | Account status. |
| `role_id` | `string` | No | Valid UUID | Updated role. |
| `team_id` | `string` | No | Valid UUID | Updated team assignment. |

---

### 🟢 Admin Reset Password (Privilege: `USER_UPDATE`)

Administrators can forcefully reset a user's password if they lose access. The system generates a random 6-digit code, sets it as the user's password, and returns it so the admin can communicate it to the user. The user's account is flagged to immediately change the password on their next login.

- **HTTP Method**: `POST`
- **Path**: `/api/user/users/:id/reset-password`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "message": "Password reset successfully. Please communicate this code to the user.",
  "code": "847291"
}
```

---

### 🟢 Delete User (Privilege: `USER_DELETE`)

Permanently remove a user account. 

- **HTTP Method**: `DELETE`
- **Path**: `/api/user/users/:id`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "message": "User deleted successfully"
}
```

---

## 👥 Dynamic Role Management Endpoints

### 🟢 List Roles (Privilege: `ROLE_VIEW`)
- **HTTP Method**: `GET`
- **Path**: `/api/user/roles`

### 🟢 Get Specific Role (Privilege: `ROLE_VIEW`)
- **HTTP Method**: `GET`
- **Path**: `/api/user/roles/:id`

### 🟢 Create Role (Privilege: `ROLE_CREATE`)
- **HTTP Method**: `POST`
- **Path**: `/api/user/roles`
- **Request Body**:
```json
{
  "name": "Custom Technician",
  "description": "Custom role for third-party service provider",
  "privileges": ["ASSET_VIEW", "WORKORDER_VIEW", "WORKORDER_UPDATE"]
}
```

### 🟢 Update Role (Privilege: `ROLE_UPDATE`)
- **HTTP Method**: `PUT`
- **Path**: `/api/user/roles/:id`

### 🟢 Delete Role (Privilege: `ROLE_DELETE`)
- **HTTP Method**: `DELETE`
- **Path**: `/api/user/roles/:id`
- **Constraints**: A role cannot be deleted if there are any users actively assigned to it. Attempting to do so will return an error `ErrRoleHasUsers`.

### 🟢 Set Role Privileges (Privilege: `ROLE_UPDATE`)
- **HTTP Method**: `PUT`
- **Path**: `/api/user/roles/:id/privileges`
- **Request Body**:
```json
{
  "privileges": ["ASSET_VIEW", "ASSET_UPDATE"]
}
```

### 🟢 Get System Privileges List (Privilege: `SYSTEM_CONFIG` or `SYSTEM_ADMIN`)
Retrieve the complete hardcoded vocabulary of capability keys.
- **HTTP Method**: `GET`
- **Path**: `/api/user/roles/privileges`
- **Success Response**: `200 OK` with JSON array of strings e.g. `["USER_VIEW", "USER_CREATE", ...]`

---

## 🏗️ Team Management Endpoints

### 🟢 List Teams (Privilege: `TEAM_VIEW`)
- **HTTP Method**: `GET`
- **Path**: `/api/user/teams`

### 🟢 Get Specific Team (Privilege: `TEAM_VIEW`)
- **HTTP Method**: `GET`
- **Path**: `/api/user/teams/:id`
- **Success Response**: `200 OK`
```json
{
  "id": "f8a71b2d-3c4e-5f6a-7b8c-9d0e1f2a3b4c",
  "name": "Maintenance Alpha",
  "manager_id": "c3b99db1-d419-48e0-bb15-081079d38bb1",
  "description": "Primary maintenance squad",
  "created_at": "2026-05-19T10:00:00Z",
  "updated_at": "2026-05-19T10:00:00Z"
}
```

### 🟢 Create Team (Privilege: `TEAM_CREATE`)
- **HTTP Method**: `POST`
- **Path**: `/api/user/teams`
- **Request Body**:
```json
{
  "name": "Maintenance Alpha",
  "manager_id": "c3b99db1-d419-48e0-bb15-081079d38bb1",
  "description": "Primary maintenance squad"
}
```

### 🟢 Update Team (Privilege: `TEAM_UPDATE`)
- **HTTP Method**: `PUT`
- **Path**: `/api/user/teams/:id`
- **Request Body**:
```json
{
  "name": "Maintenance Alpha",
  "manager_id": "e4f5g6h7-8i9j-0k1l-2m3n-4o5p6q7r8s9t",
  "description": "Updated primary maintenance squad"
}
```

### 🟢 Delete Team (Privilege: `TEAM_DELETE`)
- **HTTP Method**: `DELETE`
- **Path**: `/api/user/teams/:id`
- **Constraints**: A team cannot be deleted if there are any users actively assigned to it. Attempting to do so will return an error `ErrTeamHasUsers`.

---

## 🔒 Internal Service-to-Service Endpoints

These endpoints are bypass-routed but blocked from public gateway access. They require `X-Internal-Service: true` header to run.

### 🟢 Get User By Email (Internal)
Retrieve complete user payload, including the hashed password, for validation checks.
- **HTTP Method**: `GET`
- **Path**: `/internal/by-email?email=...`
- **Headers**: `X-Internal-Service: true`

### 🟢 Get User By ID (Internal)
Retrieve complete user payload, including structural privilege list, for validation checks.
- **HTTP Method**: `GET`
- **Path**: `/internal/by-id?id=...`
- **Headers**: `X-Internal-Service: true`

### 🟢 Get User Name By ID (Internal)
Retrieve just the full name of a user. Used by services like audit-service.
- **HTTP Method**: `GET`
- **Path**: `/internal/user-name-by-id?id=...`
- **Headers**: `X-Internal-Service: true`
