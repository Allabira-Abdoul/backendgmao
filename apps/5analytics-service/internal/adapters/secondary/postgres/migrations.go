package postgres

import (
	"log"

	"backend-gmao/apps/analytics-service/internal/core/domain"
	"gorm.io/gorm"
)

// AutoMigrate runs GORM migrations and creates materialized views.
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running analytics database migrations...")

	// GORM AutoMigrate for tables
	if err := db.AutoMigrate(
		&domain.Metric{},
		&domain.AnalyticsAssetDim{},
		&domain.AnalyticsStateEvent{},
		&domain.AnalyticsMaintenanceEvent{},
	); err != nil {
		return err
	}

	// Create Materialized View
	viewSQL := `
CREATE MATERIALIZED VIEW IF NOT EXISTS category_health_metrics_daily AS
SELECT 
    d.category_name,
    COUNT(DISTINCT d.asset_id) AS asset_count,
    COALESCE(AVG(EXTRACT(EPOCH FROM (e.completed_at - e.started_at))/3600), 0) AS mttr_hours,
    COALESCE(SUM(e.uptime_seconds) / NULLIF(SUM(e.uptime_seconds + e.downtime_seconds), 0), 1) AS availability
FROM 
    analytics_asset_dim d
LEFT JOIN 
    analytics_maintenance_events e ON d.asset_id = e.asset_id
WHERE 
    COALESCE(e.completed_at, e.started_at) >= CURRENT_DATE - INTERVAL '30 days' OR e.started_at IS NULL
GROUP BY 
    d.category_name;
`
	if err := db.Exec(viewSQL).Error; err != nil {
		log.Printf("Failed to create materialized view: %v", err)
		return err
	}

	log.Println("Analytics database migrations completed")
	return nil
}
