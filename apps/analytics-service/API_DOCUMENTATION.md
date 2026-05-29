# Analytics Service API Specification

The **Analytics Service** tracks, persists, and aggregates performance and operational metrics across the GMAO system. It serves as a central ledger for metrics related to machinery, system performance, response times, and failure rates. It registers with Consul under the name `analytics-service` and operates on port `8086`.

All external client requests are routed through the API Gateway using the prefix `/api/analytics/*`.

---

## 🔒 Security & Privileges

Endpoints require user authentication (Bearer JWT) with the appropriate dynamic privilege:

- **Record Metric**: Requires `ANALYTICS_WRITE` or `SYSTEM_ADMIN` privilege (internal services bypass this using internal service tokens).
- **List & Get Metrics**: Requires `ANALYTICS_VIEW` or `SYSTEM_ADMIN` privilege.

---

## 📡 Endpoints Specification

### 🟢 List All Metrics (Privilege: `ANALYTICS_VIEW`)

Retrieve all registered system performance metrics.

- **HTTP Method**: `GET`
- **Path**: `/api/analytics/metrics` (routed to `GET /metrics` downstream)
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
[
  {
    "id": "e44d5c41-862d-45df-bb78-ecb18360d8ef",
    "name": "cpu_utilization",
    "value": 45.2,
    "category": "system_performance",
    "timestamp": "2026-05-19T11:45:00Z"
  },
  {
    "id": "8fa1ad78-831e-4cb8-8c10-9bd74130be52",
    "name": "mean_time_to_repair",
    "value": 120.5,
    "category": "maintenance_kpi",
    "timestamp": "2026-05-19T11:50:00Z"
  }
]
```

---

### 🟢 Get Specific Metric (Privilege: `ANALYTICS_VIEW`)

Retrieve a single performance metric record by its UUID.

- **HTTP Method**: `GET`
- **Path**: `/api/analytics/metrics/:id`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "id": "e44d5c41-862d-45df-bb78-ecb18360d8ef",
  "name": "cpu_utilization",
  "value": 45.2,
  "category": "system_performance",
  "timestamp": "2026-05-19T11:45:00Z"
}
```

#### Error Responses
- **Status Code**: `400 Bad Request` (Invalid UUID format)
```json
{
  "error": "Invalid UUID format"
}
```
- **Status Code**: `404 Not Found` (Metric not found)
```json
{
  "error": "Metric not found"
}
```

---

### 🟢 List Metrics by Category (Privilege: `ANALYTICS_VIEW`)

Retrieve all metrics recorded within a specific category.

- **HTTP Method**: `GET`
- **Path**: `/api/analytics/metrics/category/:category`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
[
  {
    "id": "e44d5c41-862d-45df-bb78-ecb18360d8ef",
    "name": "cpu_utilization",
    "value": 45.2,
    "category": "system_performance",
    "timestamp": "2026-05-19T11:45:00Z"
  }
]
```

---

### 🟢 Record a Metric (Privilege: `ANALYTICS_WRITE` or Service Token)

Submit and log a new operational or performance metric.

- **HTTP Method**: `POST`
- **Path**: `/api/analytics/metrics`
- **Authentication**: Bearer JWT or Internal Service Token

#### Request Body
| Field | Type | Required | Validation | Description |
| :--- | :--- | :--- | :--- | :--- |
| `name` | `string` | Yes | Min length 2, Max length 255 | Descriptive name of the metric. |
| `value` | `float64` | Yes | Real number | Value recorded for this metric. |
| `category` | `string` | Yes | Non-empty string | Category of the metric (e.g. `system_performance`, `maintenance_kpi`). |

##### Example
```json
{
  "name": "mttr_hours",
  "value": 4.5,
  "category": "maintenance_kpi"
}
```

#### Success Response
- **Status Code**: `201 Created`
- **Body**:
```json
{
  "id": "7ca9da22-e421-4fba-bb89-11c9d9d3000b",
  "name": "mttr_hours",
  "value": 4.5,
  "category": "maintenance_kpi",
  "timestamp": "2026-05-19T12:00:00Z"
}
```

---

## 📊 Key Performance Indicators (KPIs)

The analytics service automatically calculates MTTR (Mean Time To Repair), MTBF (Mean Time Between Failures), and Global Availability based on maintenance interventions.

### 🟢 Get Global KPIs (Privilege: `ANALYTICS_VIEW`)

Retrieves aggregated KPIs across all equipments.

- **HTTP Method**: `GET`
- **Path**: `/api/analytics/kpi/global`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "status": "success",
  "data": {
    "level": "global",
    "identifier": "all",
    "mttr": 2.5,
    "mtbf": 1500.0,
    "availability": 99.8,
    "total_breakdowns": 14
  }
}
```

---

### 🟢 Get Category KPIs (Privilege: `ANALYTICS_VIEW`)

Retrieves aggregated KPIs for a specific equipment category.

- **HTTP Method**: `GET`
- **Path**: `/api/analytics/kpi/category/:category`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "status": "success",
  "data": {
    "level": "category",
    "identifier": "HVAC",
    "mttr": 3.1,
    "mtbf": 800.0,
    "availability": 99.6,
    "total_breakdowns": 5
  }
}
```

---

### 🟢 Get Asset KPIs (Privilege: `ANALYTICS_VIEW`)

Retrieves KPIs for a specific individual asset.

- **HTTP Method**: `GET`
- **Path**: `/api/analytics/kpi/asset/:asset_id`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "status": "success",
  "data": {
    "level": "asset",
    "identifier": "8fa1ad78-831e-4cb8-8c10-9bd74130be52",
    "mttr": 1.5,
    "mtbf": 450.0,
    "availability": 99.7,
    "total_breakdowns": 2
  }
}
```
