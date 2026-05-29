# Asset Service API Specification

The **Asset Service** manages the lifecycle of physical machinery, technical equipment, and operational inventory assets, including spare parts. It registers with Consul under the name `asset-service` and operates on port `8083`.

All external client requests are routed through the API Gateway using the prefix `/api/asset/*`.

The service uses a **Type/Instance** (Model-Asset) architecture. This separates the definition of catalog items (Models) from physical assets installed in the real world (Instances).

---

## 🔒 Security & Privileges

Endpoints require user authentication (Bearer JWT) with the appropriate dynamic privilege:

- **Create Asset (Model or Instance)**: Requires `ASSET_CREATE`
- **List & Get Asset (Model or Instance)**: Requires `ASSET_VIEW`
- **Update Asset**: Requires `ASSET_UPDATE`
- **Delete Asset**: Requires `ASSET_DELETE`

---

## 📡 Endpoints Specification

### 🟢 Legacy Endpoint (Backward Compatibility)
**Path**: `/api/asset/assets` (GET)
**Privilege**: `ASSET_VIEW`
Retrieves a list of mapped instances looking like the old `Asset` schema to avoid breaking existing clients.

---

### 🟢 Create Equipment Model
Define a new standard type of equipment (the blueprint).

- **HTTP Method**: `POST`
- **Path**: `/api/asset/models/equipment`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `name` | `string` | Yes | Unique name of the equipment model. |
| `category` | `string` | Yes | Operational category (e.g. HVAC, HEATING). |
| `description` | `string` | No | Details about the equipment. |

#### Success Response: `201 Created`

---

### 🟢 Create Part Model
Define a standard part type, used to manage global spare quantities.

- **HTTP Method**: `POST`
- **Path**: `/api/asset/models/parts`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `name` | `string` | Yes | Unique name of the part. |
| `category` | `string` | Yes | Part category. |
| `spare_quantity` | `integer` | No | Current global inventory count. |

#### Success Response: `201 Created`

---

### 🟢 List Equipment Models
Retrieve all equipment models.

- **HTTP Method**: `GET`
- **Path**: `/api/asset/models/equipment`
- **Authentication**: Bearer JWT

#### Success Response: `200 OK` (Array of Equipment Models)

---

### 🟢 List Part Models
Retrieve all part models.

- **HTTP Method**: `GET`
- **Path**: `/api/asset/models/parts`
- **Authentication**: Bearer JWT

#### Success Response: `200 OK` (Array of Part Models)

---

### 🟢 Create Equipment Instance (Physical Asset)
Register a physical machine linking to an `EquipmentModel`.

- **HTTP Method**: `POST`
- **Path**: `/api/asset/instances/equipment`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `code` | `string` | Yes | Unique business/inventory code (e.g., PT1). |
| `equipment_model_id` | `uuid` | Yes | The ID of the model blueprint. |
| `location` | `string` | Yes | Physical location. |
| `purchase_date` | `string` | No | RFC3339 Date. |
| `purchase_value` | `float64`| No | Monetary value. |

#### Success Response: `201 Created`

---

### 🟢 Create Part Instance (Installed Part)
Register a physical part installed on a physical equipment instance.

- **HTTP Method**: `POST`
- **Path**: `/api/asset/instances/parts`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `equipment_instance_id` | `uuid` | Yes | Where the part is installed. |
| `part_model_id` | `uuid` | Yes | What kind of part it is. |

#### Success Response: `201 Created`

---

### 🟢 List Equipment Instances
Retrieve all physical equipment with nested models, parts, and thresholds.

- **HTTP Method**: `GET`
- **Path**: `/api/asset/instances/equipment`
- **Authentication**: Bearer JWT

#### Success Response: `200 OK`

---

### 🟢 Get Equipment Instance (By ID or Code)
Retrieve a specific physical equipment.

- **HTTP Method**: `GET`
- **Path**: `/api/asset/instances/equipment/:id`
- **Path**: `/api/asset/instances/equipment/code/:code`
- **Authentication**: Bearer JWT

#### Success Response: `200 OK`
