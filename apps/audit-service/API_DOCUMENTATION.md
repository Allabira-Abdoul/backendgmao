# Audit Service API Specification

The **Audit Service** provides an immutable, centralized ledger for tracking security events, access logs, and critical operational actions across all microservices. It registers with Consul under the name `audit-service` and operates on port `8086`.

All external client requests are routed through the API Gateway using the prefix `/api/audit/*`.

---

## 🔒 Security & Privileges

This service enforces strict privilege-based RBAC for public reading, while writing is restricted to internal service-to-service communication only.

- **Read Logs**: Requires `AUDITOR` privilege.
- **Write Logs**: Restricted to internal microservices via the `/internal` routes.

---

## 📡 Endpoints Specification

### 🟢 List Audit Logs (Privilege: `AUDITOR`)

Retrieve a paginated list of system-wide audit logs, ordered by the most recent events first.

- **HTTP Method**: `GET`
- **Path**: `/api/audit/audit-logs`
- **Authentication**: Bearer JWT

#### Request Query Parameters
| Parameter | Type | Required | Default | Description |
| :--- | :--- | :--- | :--- | :--- |
| `page` | `integer` | No | `1` | The page number to fetch. |
| `per_page` | `integer` | No | `100` | The page size (limit). Max 100 recommended. |

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "meta": {
    "page": 1,
    "per_page": 100,
    "total": 5302
  },
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "service_name": "maintenance-service",
      "action": "WORKORDER_UPDATE",
      "details": "Status changed from PENDING to IN_PROGRESS",
      "user_name": "Technician Joe",
      "performed_at": "2026-05-25T10:15:30Z"
    },
    {
      "id": "987e6543-e21b-34d1-b546-537614175111",
      "service_name": "user-service",
      "action": "USER_LOGIN",
      "details": "Successful login",
      "user_name": "Admin Jane",
      "performed_at": "2026-05-25T09:12:10Z"
    }
  ]
}
```

---

## 🔒 Internal Service-to-Service Endpoints

These endpoints are bypass-routed but blocked from public gateway access. They require the internal service authentication mechanism (e.g., specific internal JWT or `X-Internal-Service` headers) to run.

### 🟢 Create Audit Log (Internal)

Record a new system action. The `Audit Service` will dynamically resolve the user's full name by calling the `user-service` if a `user_id` is provided.

- **HTTP Method**: `POST`
- **Path**: `/internal/audit-logs`
- **Headers**: `X-Internal-Service: <calling-service-name>`, `Authorization: Bearer <internal_token>`

#### Request Body
| Field | Type | Required | Validation | Description |
| :--- | :--- | :--- | :--- | :--- |
| `service_name` | `string` | Yes | - | Name of the microservice generating the log. |
| `action` | `string` | Yes | - | Event or action type (e.g., `ASSET_CREATED`). |
| `details` | `string` | No | - | Optional descriptive JSON or text details. |
| `user_id` | `string` | No | Valid UUID | The ID of the user performing the action. |

##### Example
```json
{
  "service_name": "maintenance-service",
  "action": "ASSET_CREATED",
  "details": "Created new HVAC asset ID 992",
  "user_id": "c3b99db1-d419-48e0-bb15-081079d38bb1"
}
```

#### Success Response
- **Status Code**: `201 Created`
- **Body**: Returns the recorded audit log containing the resolved `user_name` and auto-generated `performed_at` timestamp.
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "service_name": "maintenance-service",
  "action": "ASSET_CREATED",
  "details": "Created new HVAC asset ID 992",
  "user_name": "Admin Jane",
  "performed_at": "2026-05-25T11:42:15Z"
}
```
