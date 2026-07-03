# Audit Service API Specification

The **Audit Service** provides an immutable, centralized ledger for tracking security events, access logs, and critical operational actions across all microservices. It registers with Consul under the name `audit-service` and operates on port `8086`.

All external client requests are routed through the API Gateway using the prefix `/api/audit/*`.

---

## 🔒 Security & Privileges

This service enforces strict privilege-based RBAC for public reading, while writing is restricted to internal asynchronous EventBus ingestion.

- **Read Logs**: Requires `AUDIT_LOG_VIEW` privilege.
- **Write Logs**: Log ingestion is handled via RabbitMQ EventBus, there is no public or internal HTTP endpoint for writes.

---

## 📡 Endpoints Specification

### 🟢 List Audit Logs (Privilege: `AUDIT_LOG_VIEW`)

Retrieve a paginated and filterable list of system-wide audit logs, ordered by the most recent events first.

- **HTTP Method**: `GET`
- **Path**: `/api/audit/audit-logs`
- **Authentication**: Bearer JWT

#### Request Query Parameters
| Parameter | Type | Required | Default | Description |
| :--- | :--- | :--- | :--- | :--- |
| `page` | `integer` | No | `1` | The page number to fetch. |
| `per_page` | `integer` | No | `100` | The page size (limit). Max 100 recommended. |
| `service_name` | `string` | No | - | Filter by the microservice generating the log. |
| `action` | `string` | No | - | Filter by event or action type (e.g., `WORKORDER_UPDATE`). |
| `resource_type` | `string` | No | - | Filter by the type of resource affected (e.g., `WorkOrder`). |
| `resource_id` | `string` | No | - | Filter by the ID of the resource affected. |

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
      "performed_at": "2026-05-25T10:15:30Z",
      "actor_id": "c3b99db1-d419-48e0-bb15-081079d38bb1",
      "service_name": "maintenance-service",
      "action": "WORKORDER_UPDATE",
      "resource_type": "WorkOrder",
      "resource_id": "wo-9912",
      "changes": {
        "status": {
          "old": "PENDING",
          "new": "IN_PROGRESS"
        }
      }
    },
    {
      "id": "987e6543-e21b-34d1-b546-537614175111",
      "performed_at": "2026-05-25T09:12:10Z",
      "actor_id": "a1b22db1-d419-48e0-bb15-081079d38bb2",
      "service_name": "user-service",
      "action": "USER_LOGIN",
      "resource_type": "Session",
      "resource_id": "sess-44",
      "changes": null
    }
  ]
}
```

---

## 🔒 Internal Service-to-Service Ingestion (EventBus)

The Audit Service no longer exposes synchronous HTTP routes for writing logs. Instead, it subscribes to the RabbitMQ EventBus on the `audit.logs` exchange with the routing key `audit.log.*`.

Other services should publish events to this routing key with the following payload structure:

#### Event Payload structure
| Field | Type | Required | Validation | Description |
| :--- | :--- | :--- | :--- | :--- |
| `service_name` | `string` | Yes | - | Name of the microservice generating the log. |
| `action` | `string` | Yes | - | Event or action type (e.g., `ASSET_CREATED`). |
| `actor_id` | `string` | No | Valid UUID | The ID of the user performing the action. |
| `resource_type` | `string` | No | - | The type of resource being acted upon. |
| `resource_id` | `string` | No | - | The unique identifier of the resource. |
| `changes` | `JSON Object` | No | - | The specific changes applied to the resource, or generic JSON details. |

##### Example RabbitMQ Message Publish
```json
{
  "Type": "ASSET_CREATED",
  "Payload": {
    "service_name": "maintenance-service",
    "action": "ASSET_CREATED",
    "actor_id": "c3b99db1-d419-48e0-bb15-081079d38bb1",
    "resource_type": "Asset",
    "resource_id": "992",
    "changes": {
      "name": "HVAC Unit A",
      "location": "Terminal 1"
    }
  }
}
```
