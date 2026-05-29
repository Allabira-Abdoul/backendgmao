# Maintenance Service API Specification

The **Maintenance Service** manages the lifecycle of preventive and corrective maintenance work orders and technician interventions. It registers with Consul under the name `maintenance-service` and operates on port `8084`.

All external client requests are routed through the API Gateway using the prefix `/api/maintenance/*`.

---

## 🔒 Security & Privileges

Endpoints require user authentication (Bearer JWT) with the appropriate dynamic privilege:

- **Create Work Order**: Requires `WORKORDER_CREATE`
- **List & Get Work Order**: Requires `WORKORDER_VIEW`
- **Update Work Order**: Requires `WORKORDER_UPDATE`
- **Delete Work Order**: Requires `WORKORDER_DELETE`

---

## 📡 Endpoints Specification

### 🟢 List Work Orders (Privilege: `WORKORDER_VIEW`)

Retrieve all registered work orders.

- **HTTP Method**: `GET`
- **Path**: `/api/maintenance/work-orders` (routed to `GET /work-orders` downstream)
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
[
  {
    "id": "2b09db24-f719-482a-a9bd-8311d9d93bb1",
    "title": "Replace Hydraulic Pump Seals",
    "description": "Slow fluid leak identified during routine inspection",
    "asset": {
      "id": "7b09db24-f719-482a-a9bd-8311d9d93bb1",
      "name": "Hydraulic Pump"
    },
    "type": "INTERVENTION",
    "scheduled_at": "2026-05-20T10:00:00Z",
    "priority": "HIGH",
    "status": "IN_PROGRESS",
    "maintenance_category": "CORRECTIVE",
    "maintenance_type": "CURATIVE",
    "is_metric_measurement": false,
    "assigned_to": {
      "id": "c3b99db1-d419-48e0-bb15-081079d38bb1",
      "name": "Jane Technician"
    },
    "interventions": [
      {
        "id": "9b09db24-f719-482a-a9bd-8311d9d93cc2",
        "maintenance_category": "CORRECTIVE",
        "maintenance_type": "CURATIVE",
        "is_metric_measurement": false,
        "started_at": "2026-05-19T10:15:00Z",
        "ended_at": "2026-05-19T11:00:00Z",
        "performed_by": {
          "id": "c3b99db1-d419-48e0-bb15-081079d38bb1",
          "name": "Jane Technician"
        },
        "created_at": "2026-05-19T11:00:00Z",
        "updated_at": "2026-05-19T11:00:00Z"
      }
    ],
    "inspections": [],
    "created_at": "2026-05-19T10:15:00Z",
    "updated_at": "2026-05-19T11:00:00Z"
  }
]
```

---

### 🟢 Get Specific Work Order (Privilege: `WORKORDER_VIEW`)

Retrieve a single work order by its unique UUID, including all recorded interventions.

- **HTTP Method**: `GET`
- **Path**: `/api/maintenance/work-orders/:id`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**: (Returns a single work order JSON object, matching the structure above)

---

### 🟢 Create Work Order (Privilege: `WORKORDER_CREATE`)

Register a new maintenance ticket/work order.

- **HTTP Method**: `POST`
- **Path**: `/api/maintenance/work-orders`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Validation | Description |
| :--- | :--- | :--- | :--- | :--- |
| `title` | `string` | Yes | Min 2, max 255 chars | Short summary of the work. |
| `description` | `string` | No | None | Detailed context of the work order. |
| `asset_id` | `string` | Yes | Valid UUID | Asset targeted by the work order. |
| `type` | `string` | No | `INTERVENTION`, `INSPECTION` | Defaults to INTERVENTION. |
| `scheduled_at` | `string` | No | ISO8601 Date | Scheduled start time for the order. |
| `priority` | `string` | Yes | `LOW`, `MEDIUM`, `HIGH`, `CRITICAL` | Task urgency level. |
| `maintenance_category` | `string` | No | `PREVENTIVE`, `CORRECTIVE` | Maintenance Category. |
| `maintenance_type` | `string` | No | `CURATIVE`, `PALLIATIVE`, `SYSTEMATIC`, etc. | Maintenance Type. |
| `is_metric_measurement`| `boolean`| No | boolean | True if measurements are expected. |
| `assigned_to` | `string` | No | Valid UUID | Assigned technician or team lead user ID. |

##### Example
```json
{
  "title": "Replace Hydraulic Pump Seals",
  "description": "Slow fluid leak identified during routine inspection",
  "asset_id": "7b09db24-f719-482a-a9bd-8311d9d93bb1",
  "type": "INTERVENTION",
  "scheduled_at": "2026-05-20T10:00:00Z",
  "priority": "HIGH",
  "assigned_to": "c3b99db1-d419-48e0-bb15-081079d38bb1"
}
```

#### Success Response
- **Status Code**: `201 Created`
- **Body**: (Returns the newly created work order JSON object)

---

### 🟢 Update Work Order (Privilege: `WORKORDER_UPDATE`)

Update details, progress status, priority, or assignment of an existing work order.

- **HTTP Method**: `PUT`
- **Path**: `/api/maintenance/work-orders/:id`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Validation / Allowed Values | Description |
| :--- | :--- | :--- | :--- | :--- |
| `title` | `string` | No | Min 2, max 255 chars | Updated title. |
| `description` | `string` | No | None | Updated description. |
| `status` | `string` | No | `PENDING`, `IN_PROGRESS`, `COMPLETED`, `CANCELLED` | Work order workflow status. |
| `priority` | `string` | No | `LOW`, `MEDIUM`, `HIGH`, `CRITICAL` | Updated priority. |
| `assigned_to` | `string` | No | Valid UUID | Updated technician assignment. |

##### Example
```json
{
  "status": "COMPLETED"
}
```

#### Success Response
- **Status Code**: `200 OK`
- **Body**: (Returns the updated work order JSON object)

---

### 🟢 Delete Work Order (Privilege: `WORKORDER_DELETE`)

Delete a work order ticket from the active repository.

- **HTTP Method**: `DELETE`
- **Path**: `/api/maintenance/work-orders/:id`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "message": "Work order deleted successfully"
}
```

---

## 🛠️ Interventions (Sub-resource under Work Orders)

Technicians can log separate operational actions (interventions) under an active work order.

### 🟢 Record Intervention (Privilege: `WORKORDER_UPDATE`)

Record a new session of work or repair completed on this work order.

- **HTTP Method**: `POST`
- **Path**: `/api/maintenance/work-orders/:id/interventions`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Validation | Description |
| :--- | :--- | :--- | :--- | :--- |
| `description` | `string` | Yes | Min 2 chars | Summary of work done. |
| `performed_by` | `string` | Yes | Valid UUID | User ID of the performing technician. |
| `measurements` | `array` | No | - | List of metric measurements taken. |

##### Example
```json
{
  "description": "Inspected leak severity and prepared gasket replacements",
  "performed_by": "c3b99db1-d419-48e0-bb15-081079d38bb1"
}
```

#### Success Response
- **Status Code**: `201 Created`
- **Body**: (Returns the recorded intervention object)

---

### 🟢 Start Intervention (Privilege: `WORKORDER_UPDATE`)

Marks an intervention as started (updates `started_at` to now and turns the asset DOWN).

- **HTTP Method**: `POST`
- **Path**: `/api/maintenance/work-orders/:id/interventions/:inv_id/start`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`

---

### 🟢 End Intervention (Privilege: `WORKORDER_UPDATE`)

Marks an intervention as ended (updates `ended_at` to now, restores asset to OPERATIONAL, triggers analytics).

- **HTTP Method**: `POST`
- **Path**: `/api/maintenance/work-orders/:id/interventions/:inv_id/end`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`

---

### 🟢 List Interventions (Privilege: `WORKORDER_VIEW`)

List all recorded interventions for a given work order.

- **HTTP Method**: `GET`
- **Path**: `/api/maintenance/work-orders/:id/interventions`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**: (Returns JSON array of interventions)

---

## 🛠️ Inspections (Sub-resource under Work Orders)

Technicians can log separate observation actions (inspections) under an active work order.

### 🟢 Record Inspection (Privilege: `WORKORDER_UPDATE`)

Record a new inspection observation.

- **HTTP Method**: `POST`
- **Path**: `/api/maintenance/work-orders/:id/inspections`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Validation | Description |
| :--- | :--- | :--- | :--- | :--- |
| `observations` | `string` | Yes | Min 2 chars | Summary of observations. |
| `performed_by` | `string` | Yes | Valid UUID | User ID of the performing technician. |
| `measurements` | `array` | No | - | List of metric measurements taken. |

##### Example
```json
{
  "observations": "Vibration levels are high",
  "performed_by": "c3b99db1-d419-48e0-bb15-081079d38bb1"
}
```

#### Success Response
- **Status Code**: `201 Created`

---

### 🟢 Start Inspection (Privilege: `WORKORDER_UPDATE`)

Marks an inspection as started.

- **HTTP Method**: `POST`
- **Path**: `/api/maintenance/work-orders/:id/inspections/:ins_id/start`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`

---

### 🟢 End Inspection (Privilege: `WORKORDER_UPDATE`)

Marks an inspection as ended.

- **HTTP Method**: `POST`
- **Path**: `/api/maintenance/work-orders/:id/inspections/:ins_id/end`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
