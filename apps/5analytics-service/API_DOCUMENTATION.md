# Analytics Service API Specification

The **Analytics Service** tracks, persists, and aggregates performance and operational metrics across the GMAO system. It serves as a central ledger for metrics related to machinery, system performance, response times, and failure rates. It uses CQRS and Event Sourcing via RabbitMQ to build high-speed materialized views. It registers with Consul under the name `analytics-service` and operates on port `8086`.

All external client requests are routed through the API Gateway using the prefix `/api/analytics/*`.

---

## 🔒 Security & Privileges

Endpoints require user authentication (Bearer JWT) with the appropriate dynamic privilege:

- **Record Metric**: Requires `ANALYTICS_WRITE` or `SYSTEM_ADMIN` privilege (internal services bypass this using internal service tokens).
- **List & Get Metrics/KPIs**: Requires `ANALYTICS_VIEW` or `SYSTEM_ADMIN` privilege.

---

## 📡 Endpoints Specification

### 🟢 List All Metrics (Privilege: `ANALYTICS_VIEW`)

Retrieve all registered system performance metrics.

- **HTTP Method**: `GET`
- **Path**: `/api/analytics/metrics`
- **Authentication**: Bearer JWT

---

### 🟢 Get Specific Metric (Privilege: `ANALYTICS_VIEW`)

Retrieve a single performance metric record by its UUID.

- **HTTP Method**: `GET`
- **Path**: `/api/analytics/metrics/:id`
- **Authentication**: Bearer JWT

---

### 🟢 List Metrics by Category (Privilege: `ANALYTICS_VIEW`)

Retrieve all metrics recorded within a specific category.

- **HTTP Method**: `GET`
- **Path**: `/api/analytics/metrics/category/:category`
- **Authentication**: Bearer JWT

---

### 🟢 Record a Metric (Privilege: `ANALYTICS_WRITE` or Service Token)

Submit and log a new operational or performance metric.

- **HTTP Method**: `POST`
- **Path**: `/api/analytics/metrics`
- **Authentication**: Bearer JWT or Internal Service Token

---

## 📊 Key Performance Indicators (KPIs) - CQRS

The analytics service uses event sourcing to listen to `asset.events` and `maintenance.events`. It stores raw facts in a dimension/fact model and exposes them via Materialized Views that refresh daily.

### 🟢 Get Category Health Metrics (Privilege: `ANALYTICS_VIEW`)

Retrieves aggregated KPIs (Availability, MTTR) grouped by Asset Category over the last 30 days.

- **HTTP Method**: `GET`
- **Path**: `/api/analytics/kpis/categories/health`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
{
  "status": "success",
  "data": [
    {
      "category_name": "HVAC",
      "asset_count": 12,
      "period": {
        "start": "0001-01-01T00:00:00Z",
        "end": "0001-01-01T00:00:00Z"
      },
      "metrics": {
        "total_uptime_seconds": 0,
        "total_downtime_seconds": 0,
        "availability_percentage": 0.99,
        "mttr_hours": 3.4,
        "mtbf_hours": 0
      }
    },
    {
      "category_name": "Baggage Handling",
      "asset_count": 5,
      "period": {
        "start": "0001-01-01T00:00:00Z",
        "end": "0001-01-01T00:00:00Z"
      },
      "metrics": {
        "total_uptime_seconds": 0,
        "total_downtime_seconds": 0,
        "availability_percentage": 0.87,
        "mttr_hours": 1.2,
        "mtbf_hours": 0
      }
    }
  ]
}
```
