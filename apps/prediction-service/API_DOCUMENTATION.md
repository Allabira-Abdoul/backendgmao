# Predictive Maintenance Service API Specification

The **Predictive Maintenance Service** leverages sensor telemetry metrics, historical breakdown logs, and machine learning models to forecast asset failures and estimated wear lifespans. It registers with Consul under the name `prediction-service` and operates on port `8085`.

All external client requests are routed through the API Gateway using the prefix `/api/prediction/*`.

---

## 🔒 Security & Privileges

Endpoints require user authentication (Bearer JWT) with the appropriate dynamic privilege:

- **Create Prediction**: Requires `PREDICTION_CREATE`
- **List & Get Predictions**: Requires `PREDICTION_VIEW`

---

## 📡 Endpoints Specification

### 🟢 List All Predictions (Privilege: `PREDICTION_VIEW`)

Retrieve all logged historical failure predictions.

- **HTTP Method**: `GET`
- **Path**: `/api/prediction/predictions` (routed to `GET /predictions` downstream)
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**:
```json
[
  {
    "id": "3b09db24-f719-482a-a9bd-8311d9d93bb1",
    "asset_id": "7b09db24-f719-482a-a9bd-8311d9d93bb1",
    "failure_probability": 84.5,
    "predicted_failure_date": "2026-06-15T00:00:00Z",
    "created_at": "2026-05-19T10:15:00Z"
  }
]
```

---

### 🟢 Get Specific Prediction (Privilege: `PREDICTION_VIEW`)

Retrieve a single failure prediction record by its UUID.

- **HTTP Method**: `GET`
- **Path**: `/api/prediction/predictions/:id`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**: (Returns a single prediction JSON object, matching the structure above)

---

### 🟢 List Predictions for Asset (Privilege: `PREDICTION_VIEW`)

Retrieve all historical predictions recorded for a specific physical asset.

- **HTTP Method**: `GET`
- **Path**: `/api/prediction/predictions/asset/:assetId`
- **Authentication**: Bearer JWT

#### Success Response
- **Status Code**: `200 OK`
- **Body**: (Returns a JSON array of predictions logged for the asset)

---

### 🟢 Create Prediction (Privilege: `PREDICTION_CREATE`)

Log a new AI/ML generated health prediction.

- **HTTP Method**: `POST`
- **Path**: `/api/prediction/predictions`
- **Authentication**: Bearer JWT

#### Request Body
| Field | Type | Required | Validation | Description |
| :--- | :--- | :--- | :--- | :--- |
| `asset_id` | `string` | Yes | Valid UUID | Asset targeted by the prediction. |
| `failure_probability` | `float64` | Yes | Real number `[0.0, 100.0]` | Forecasted percentage chance of breakdown. |
| `predicted_failure_date` | `string (RFC3339)` | Yes | Valid date format | Predicted breakdown timeline date. |

##### Example
```json
{
  "asset_id": "7b09db24-f719-482a-a9bd-8311d9d93bb1",
  "failure_probability": 84.50,
  "predicted_failure_date": "2026-06-15T00:00:00Z"
}
```

#### Success Response
- **Status Code**: `201 Created`
- **Body**: (Returns the recorded prediction JSON response object)
