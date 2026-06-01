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
| `is_serialized` | `boolean` | No | Whether the part is tracked individually. |

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
| `part_model_id` | `uuid` | Yes | What kind of part it is. |
| `serial_number` | `string` | Yes | The unique serial number for this part instance. |
| `equipment_instance_id`| `uuid` | No | Where the part is installed. |
| `current_location`| `string` | Yes | Current location of the part. |

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

---

### 🟢 Move Part Instance
Move a part instance to a new location or equipment.

- **HTTP Method**: `POST`
- **Path**: `/api/asset/instances/parts/:id/move`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `equipment_instance_id`| `uuid` | No | The ID of the equipment to move the part to. |
| `current_location`| `string` | Yes | The new location of the part. |

#### Success Response: `200 OK`

---

### 🟢 Consume Part
Consume a non-serialized part from global inventory.

- **HTTP Method**: `POST`
- **Path**: `/api/asset/actions/consume-part`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `part_model_id` | `uuid` | Yes | The ID of the part model to consume. |
| `quantity` | `integer` | Yes | The quantity to consume (minimum 1). |
| `work_order_id` | `uuid` | No | The ID of the associated work order. |
| `notes` | `string` | No | Any notes regarding the consumption. |

#### Success Response: `200 OK`

---

### 🟢 Ingest Measurement
Ingest a new telemetry measurement for an equipment or part instance.

- **HTTP Method**: `POST`
- **Path**: `/api/asset/measurements`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `equipment_instance_id`| `uuid` | No | The ID of the equipment instance. |
| `part_instance_id` | `uuid` | No | The ID of the part instance. |
| `metric_name` | `string` | Yes | The name of the metric (e.g., Temperature, Pressure). |
| `value` | `float64`| Yes | The recorded value of the measurement. |
| `unit` | `string` | Yes | The unit of the measurement. |
| `recorded_at` | `string` | No | RFC3339 Date. Defaults to current time if omitted. |

*(Note: Provide either `equipment_instance_id` or `part_instance_id`)*

#### Success Response: `201 Created`

---

### 🟢 Get Measurements
Retrieve measurements for a specific target (equipment or part).

- **HTTP Method**: `GET`
- **Path**: `/api/asset/measurements/:targetType/:targetID`
- **Authentication**: Bearer JWT
- **Query Parameters**:
  - `since` (string, optional): Filter measurements recorded after this date (RFC3339).

*(Note: `:targetType` should be either `equipment` or `part`)*

#### Success Response: `200 OK`

---

### 🟢 Create Supplier
Register a new asset supplier.

- **HTTP Method**: `POST`
- **Path**: `/api/asset/suppliers`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `name` | `string` | Yes | Name of the supplier. |
| `contact_info` | `string` | No | Contact information. |

#### Success Response: `201 Created`

---

### 🟢 List Suppliers
Retrieve all suppliers.

- **HTTP Method**: `GET`
- **Path**: `/api/asset/suppliers`
- **Authentication**: Bearer JWT

#### Success Response: `200 OK`

---

### 🟢 Link Supplier to Equipment Model
Add a supplier reference to an equipment model.

- **HTTP Method**: `POST`
- **Path**: `/api/asset/models/equipment/:id/suppliers`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `supplier_id` | `uuid` | Yes | ID of the supplier. |
| `supplier_reference_code`| `string` | Yes | Supplier's reference code (e.g., ITS3). |
| `technical_doc_reference`| `string` | No | Page or link for technical doc. |

#### Success Response: `201 Created`

---

### 🟢 Link Supplier to Part Model
Add a supplier reference to a part model.

- **HTTP Method**: `POST`
- **Path**: `/api/asset/models/parts/:id/suppliers`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `supplier_id` | `uuid` | Yes | ID of the supplier. |
| `supplier_reference_code`| `string` | Yes | Supplier's reference code. |
| `technical_doc_reference`| `string` | No | Page or link for technical doc. |

#### Success Response: `201 Created`
