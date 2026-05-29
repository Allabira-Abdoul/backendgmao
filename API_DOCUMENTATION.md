# GMAO Master API Documentation

Welcome to the GMAO (Gestion de Maintenance Assistée par Ordinateur) backend API documentation. This document serves as the master specification for our high-performance, resilient microservices ecosystem built in Go.

---

## 🏗️ System Architecture

Our backend employs an API Gateway routing pattern integrated with Consul for dynamic service registry and discovery.

```mermaid
graph TD
    Client[Client Browser / Mobile] -->|HTTPS Requests| Gateway[API Gateway :8080]
    Gateway -->|Service Discovery| Consul[Consul Registry :8500]
    Gateway -->|Dynamic Proxy /api/{service}| MS[Microservices]
    
    subgraph Microservices [Microservices Network]
        AuthService[auth-service]
        UserService[user-service]
        AssetService[asset-service]
        MaintService[maintenance-service]
        PredService[prediction-service]
        AnalyService[analytics-service]
        AuditService[audit-service]
    end
    
    MS -->|Data Persistence| DB[(PostgreSQL)]
```

---

## ⚡ Gateway Dynamic Routing & Path Prefixes

The API Gateway listens on port **8080** and routes requests dynamically to downstream services based on the pattern:
`http://localhost:8080/api/{service-suffix}/*path`

| Dynamic Prefix Path | Downstream Service | Consul Registered Name | Default Downstream Port |
| :--- | :--- | :--- | :--- |
| `/api/auth/*` | Authentication Service | `auth-service` | `8081` |
| `/api/user/*` | User & RBAC Service | `user-service` | `8082` |
| `/api/asset/*` | Asset Service | `asset-service` | `8083` |
| `/api/maintenance/*` | Maintenance Service | `maintenance-service` | `8084` |
| `/api/prediction/*` | Prediction Service | `prediction-service` | `8085` |
| `/api/analytics/*` | Analytics Service | `analytics-service` | `8086` |
| `/api/audit/*` | Audit Service | `audit-service` | `8087` |

> [!WARNING]
> **SSRF Protection:** The API Gateway drops internal headers (e.g. `X-Internal-Service`) and strictly sanitizes paths. Downstream paths beginning with `/internal` are strictly forbidden for external users and will return a `403 Forbidden` response.

---

## 🔑 Role-Based Access Control (RBAC) & Privileges

Our authorization system uses a strict privilege-based RBAC model. Users are assigned a **Role**, which consists of a set of unique capability keys called **Privileges**. Downstream service endpoints require specific privileges to be present in the user's JWT.

### Privilege List by Functional Domain

| Domain | Privilege Key | Description |
| :--- | :--- | :--- |
| **User Management** | `USER_VIEW` | View list of users and specific user details |
| | `USER_CREATE` | Register new users |
| | `USER_UPDATE` | Modify user records and statuses |
| | `USER_DELETE` | Permanently remove a user |
| | `USER_ASSIGN_ROLE` | Alter roles assigned to users |
| **Role Management** | `ROLE_VIEW` | View roles and their associated privilege lists |
| | `ROLE_CREATE` | Create new system roles |
| | `ROLE_UPDATE` | Alter permissions list of existing roles |
| | `ROLE_DELETE` | Remove a system role |
| **Asset Management** | `ASSET_VIEW` | View machinery inventory and locations |
| | `ASSET_CREATE` | Add new assets to the system |
| | `ASSET_UPDATE` | Update asset states (e.g. DOWN, OPERATIONAL) |
| | `ASSET_DELETE` | Delete/scrap physical assets |
| | `ASSET_TRANSFER` | Transfer assets across locations/teams |
| **Work Order Management**| `WORKORDER_VIEW` | View maintenance work orders |
| | `WORKORDER_CREATE` | Draft a new work order |
| | `WORKORDER_UPDATE` | Edit details, priority, or status |
| | `WORKORDER_DELETE` | Cancel or delete a work order |
| | `WORKORDER_ASSIGN` | Assign a work order to a technician or team |
| | `WORKORDER_CLOSE` | Close completed work orders |
| **Maintenance** | `MAINTENANCE_VIEW` | Inspect historical intervention logs |
| | `MAINTENANCE_PLAN_CREATE`| Build preventive maintenance plans |
| | `MAINTENANCE_SCHEDULE` | Schedule interventions |
| **Analytics** | `ANALYTICS_VIEW` | Read categorical performance metrics |
| **System Admin** | `SYSTEM_ADMIN` | Full admin capabilities |
| | `SYSTEM_AUDIT_VIEW` | Read low-level secure audit logs |

---

## 🌐 Response Format & Errors

All JSON responses follow a standardized payload structure. Success responses return the data payload directly, while error responses return an error message.

### Success Response Format (Single / Array)
```json
{
  "id": "...",
  "field": "value"
}
// or [ { ... }, { ... } ]
```

### Standard Error Response Format
Our backend returns a structured error JSON detailing the specific issue.
```json
{
  "error": "Detailed error message explaining what went wrong"
}
```

#### Common Errors
- `400 Bad Request`: Validation failure or malformed payload.
- `401 Unauthorized`: Provided Bearer JWT is missing, invalid, or expired.
- `403 Forbidden`: User does not possess the required privilege to perform the action.
- `404 Not Found`: The requested resource does not exist.
- `500 Internal Server Error`: Unexpected backend failure.

---

## 📂 Service-Specific Documentations

For detailed request payloads, validation limits, and response schemas, navigate directly to the specific service's markdown specification:

1. 🛡️ [API Gateway API Specification](file:///d:/1.%20STAGE%20DE%20FIN%20D%27ETUDES/2.PROJET/backend/apps/api-gateway/API_DOCUMENTATION.md)
2. 🔑 [Authentication Service API Specification](file:///d:/1.%20STAGE%20DE%20FIN%20D%27ETUDES/2.PROJET/backend/apps/auth-service/API_DOCUMENTATION.md)
3. 👤 [User & RBAC Service API Specification](file:///d:/1.%20STAGE%20DE%20FIN%20D%27ETUDES/2.PROJET/backend/apps/user-service/API_DOCUMENTATION.md)
4. 🏗️ [Asset Service API Specification](file:///d:/1.%20STAGE%20DE%20FIN%20D%27ETUDES/2.PROJET/backend/apps/asset-service/API_DOCUMENTATION.md)
5. ⚙️ [Maintenance Service API Specification](file:///d:/1.%20STAGE%20DE%20FIN%20D%27ETUDES/2.PROJET/backend/apps/maintenance-service/API_DOCUMENTATION.md)
6. 📈 [Prediction Service API Specification](file:///d:/1.%20STAGE%20DE%20FIN%20D%27ETUDES/2.PROJET/backend/apps/prediction-service/API_DOCUMENTATION.md)
7. 📊 [Analytics Service API Specification](file:///d:/1.%20STAGE%20DE%20FIN%20D%27ETUDES/2.PROJET/backend/apps/analytics-service/API_DOCUMENTATION.md)
8. 📜 [Audit Service API Specification](file:///d:/1.%20STAGE%20DE%20FIN%20D%27ETUDES/2.PROJET/backend/apps/audit-service/API_DOCUMENTATION.md)
