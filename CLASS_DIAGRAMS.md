# Backend Services Class Diagrams

This document contains the structural class diagrams (Domain Models) for each microservice in the GMAO backend. These diagrams illustrate the core entities, their properties, and their relationships.

## Analytics Service

```mermaid
classDiagram
    class TimeRange {
        +Time start
        +Time end
    }
    
    class CoreMetrics {
        +float64 total_uptime_seconds
        +float64 total_downtime_seconds
        +float64 availability_percentage
        +float64 mttr_hours
        +float64 mtbf_hours
    }
    
    class AssetHealthMetrics {
        +UUID asset_id
        +TimeRange period
        +CoreMetrics metrics
    }
    AssetHealthMetrics *-- TimeRange
    AssetHealthMetrics *-- CoreMetrics
    
    class CategoryHealthMetrics {
        +string category_name
        +int asset_count
        +TimeRange period
        +CoreMetrics metrics
    }
    CategoryHealthMetrics *-- TimeRange
    CategoryHealthMetrics *-- CoreMetrics
    
    class AnalyticsAssetDim {
        +UUID asset_id
        +UUID model_id
        +string category_name
        +Time created_at
    }
    
    class AnalyticsStateEvent {
        +UUID id
        +UUID asset_id
        +string old_state
        +string new_state
        +Time timestamp
    }
    
    class AnalyticsMaintenanceEvent {
        +UUID work_order_id
        +UUID asset_id
        +string type
        +string maintenance_type
        +Time started_at
        +Time completed_at
        +float64 uptime_seconds
        +float64 downtime_seconds
    }
    
    class Metric {
        +UUID id
        +string name
        +float64 value
        +string category
        +Time timestamp
    }
```

## Asset Service

```mermaid
classDiagram
    class EquipmentModel {
        +UUID id
        +string name
        +string category
        +string description
        +Time created_at
        +Time updated_at
    }
    EquipmentModel "1" *-- "many" ModelSupplier
    EquipmentModel "1" *-- "many" EquipmentModelPartRequirement
    
    class PartModel {
        +UUID id
        +string name
        +string category
        +int spare_quantity
        +bool is_serialized
    }
    PartModel "1" *-- "many" ModelSupplier
    
    class EquipmentModelPartRequirement {
        +UUID id
        +UUID equipment_model_id
        +UUID part_model_id
        +int quantity
    }
    EquipmentModelPartRequirement --> PartModel
    
    class EquipmentInstance {
        +UUID id
        +string code
        +UUID equipment_model_id
        +UUID supplier_id
        +string status
        +string location
        +float64 usage_hours
    }
    EquipmentInstance --> EquipmentModel
    EquipmentInstance --> Supplier
    EquipmentInstance "1" *-- "many" PartInstance
    
    class PartInstance {
        +UUID id
        +UUID equipment_instance_id
        +UUID part_model_id
        +UUID supplier_id
        +string serial_number
        +string status
        +string current_location
    }
    PartInstance --> PartModel
    PartInstance --> Supplier
    
    class PartConsumptionLog {
        +UUID id
        +UUID part_model_id
        +int quantity_used
        +UUID work_order_id
        +UUID consumed_by
        +string notes
    }

    class Consumable {
        +UUID id
        +string name
        +string category
        +string unit_of_measure
        +int total_stock
        +Time created_at
        +Time updated_at
    }
    Consumable "1" *-- "many" ModelSupplier
    Consumable "1" *-- "many" ConsumableLocationStock

    class ConsumableLocationStock {
        +UUID id
        +UUID consumable_id
        +string location
        +int quantity
    }

    class ConsumableConsumptionLog {
        +UUID id
        +UUID consumable_id
        +int quantity_used
        +UUID work_order_id
        +UUID consumed_by
        +string notes
    }
    
    class Supplier {
        +UUID id
        +string name
        +string contact_info
    }
    
    class ModelSupplier {
        +UUID id
        +UUID supplier_id
        +UUID equipment_model_id
        +UUID part_model_id
        +UUID consumable_id
        +string supplier_reference_code
        +string technical_doc_reference
    }
    ModelSupplier --> Supplier
```

## Audit Service

```mermaid
classDiagram
    class AuditLog {
        +UUID id
        +Time performed_at
        +UUID actor_id
        +string service_name
        +string action
        +string resource_type
        +string resource_id
        +JSON changes
    }
```

## Auth Service

```mermaid
classDiagram
    class Session {
        +UUID id
        +UUID user_id
        +string access_token
        +string refresh_token
        +Time access_expired_at
        +Time refresh_expired_at
        +Time created_at
    }
    
    class User {
        +UUID id
        +string full_name
        +string email
        +string password
        +AccountStatus status
        +string role_name
        +string[] privileges
    }
```

## Maintenance Service

```mermaid
classDiagram
    class EquipmentModelMaintenanceRule {
        +UUID id
        +UUID equipment_model_id
        +string rule_name
        +float64 interval_hours
        +int interval_months
    }
    
    class EquipmentInstanceMaintenanceState {
        +UUID id
        +UUID equipment_instance_id
        +UUID maintenance_rule_id
        +Time last_maintenance_at
        +float64 last_maintenance_usage_hours
    }
    EquipmentInstanceMaintenanceState --> EquipmentModelMaintenanceRule

    class MaintenanceSchedule {
        +UUID id
        +UUID asset_id
        +string title
        +string description
        +string frequency
        +int interval_months
        +float64 interval_hours
        +Time start_date
        +Time end_date
        +Time next_scheduled_date
        +float64 next_scheduled_usage
        +string maintenance_category
        +string maintenance_type
        +bool is_active
        +bool require_counter_reading
    }
    MaintenanceSchedule "1" *-- "many" WorkOrder
    MaintenanceSchedule "1" *-- "many" Inspection
    
    class WorkOrder {
        +UUID id
        +string title
        +string description
        +UUID asset_id
        +UUID schedule_id
        +string priority
        +string status
        +string type
        +Time scheduled_start_date
        +Time scheduled_end_date
        +string maintenance_category
        +string maintenance_type
        +bool is_metric_measurement
        +UUID assigned_to
    }
    WorkOrder "1" *-- "many" Intervention
    
    class Intervention {
        +UUID id
        +UUID work_order_id
        +string description
        +string maintenance_category
        +string maintenance_type
        +bool is_metric_measurement
        +Time date
        +Time started_at
        +Time ended_at
        +UUID performed_by
    }
    Intervention "1" *-- "many" Measurement
    
    class Inspection {
        +UUID id
        +UUID asset_id
        +UUID schedule_id
        +string observations
        +float64 usage_hours_recorded
        +bool requires_attention
        +string attention_reason
        +Time date
        +Time started_at
        +Time ended_at
        +UUID performed_by
    }
    Inspection "1" *-- "many" Measurement
    
    class Measurement {
        +UUID id
        +UUID asset_id
        +UUID intervention_id
        +UUID inspection_id
        +UUID component_id
        +string metric_name
        +float64 value
        +string unit
        +bool is_threshold_breached
        +Time recorded_at
        +UUID recorded_by
    }
    
    class DefectAlert {
        +UUID id
        +UUID asset_id
        +UUID reported_by
        +string title
        +string description
        +string image_url
        +string status
    }
```

## Prediction Service

```mermaid
classDiagram
    class Prediction {
        +UUID id
        +UUID asset_id
        +float64 failure_probability
        +Time predicted_failure_date
        +Time created_at
    }
```

## User Service

```mermaid
classDiagram
    class User {
        +UUID id
        +string full_name
        +string email
        +string password
        +bool must_change_password
        +AccountStatus status
        +UUID role_id
        +UUID team_id
        +string location
    }
    User --> Role
    User --> Team
    
    class Role {
        +UUID id
        +string name
        +string description
        +string[] privileges
    }
    Role "1" *-- "many" RolePrivilege
    
    class RolePrivilege {
        +UUID role_id
        +string privilege
    }
    
    class Team {
        +UUID id
        +string name
        +UUID manager_id
        +string description
        +string location
    }
```

## API Gateway
_The API Gateway proxies requests based on Consul service discovery and strips out internal headers. It operates entirely as an HTTP router and middleware provider, thus having no persistent domain entities or database._
