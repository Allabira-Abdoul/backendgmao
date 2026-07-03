package postgres

import (
	"context"
	"fmt"
	"time"

	"backend-gmao/apps/audit-service/internal/core/domain"
	"backend-gmao/apps/audit-service/internal/core/ports/secondary"
	"gorm.io/gorm"
)

type auditRepository struct {
	db *gorm.DB
}

// NewAuditRepository creates a GORM audit log repository.
func NewAuditRepository(db *gorm.DB) secondary.AuditRepository {
	return &auditRepository{db: db}
}

// InitPartitionedTable drops the existing non-partitioned table (if any) and creates a partitioned one.
func InitPartitionedTable(db *gorm.DB) error {
	// Check if the table already exists
	var exists bool
	err := db.Raw(`SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE  table_schema = 'public'
		AND    table_name   = 'audit_logs'
	);`).Scan(&exists).Error

	if err == nil && exists {
		// Table already exists, just ensure current partitions are created
		now := time.Now()
		createPartition(db, now)
		createPartition(db, now.AddDate(0, 1, 0))
		return nil
	}

	createTableSQL := `
	CREATE TABLE audit_logs (
		id UUID DEFAULT gen_random_uuid(),
		performed_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
		actor_id UUID,
		service_name TEXT NOT NULL,
		action TEXT NOT NULL,
		resource_type TEXT,
		resource_id TEXT,
		changes JSONB,
		PRIMARY KEY (id, performed_at)
	) PARTITION BY RANGE (performed_at);`

	if err := db.Exec(createTableSQL).Error; err != nil {
		return fmt.Errorf("failed to create partitioned table: %w", err)
	}

	// Create partition for current month and next month
	now := time.Now()
	createPartition(db, now)
	createPartition(db, now.AddDate(0, 1, 0))

	return nil
}

func createPartition(db *gorm.DB, t time.Time) {
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	end := start.AddDate(0, 1, 0)
	
	tableName := fmt.Sprintf("audit_logs_y%dm%02d", t.Year(), t.Month())
	
	sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s PARTITION OF audit_logs FOR VALUES FROM ('%s') TO ('%s');`, 
		tableName, start.Format("2006-01-02"), end.Format("2006-01-02"))
	
	db.Exec(sql)
}

func (r *auditRepository) Save(ctx context.Context, log domain.AuditLog) error {
	return r.db.WithContext(ctx).Create(&log).Error
}

func (r *auditRepository) Find(ctx context.Context, filter domain.AuditFilter) ([]domain.AuditLog, int64, error) {
	var logs []domain.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.AuditLog{})

	if filter.ServiceName != "" {
		query = query.Where("service_name = ?", filter.ServiceName)
	}
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if filter.ResourceType != "" {
		query = query.Where("resource_type = ?", filter.ResourceType)
	}
	if filter.ResourceID != "" {
		query = query.Where("resource_id = ?", filter.ResourceID)
	}
	if filter.ActorID != nil {
		query = query.Where("actor_id = ?", filter.ActorID)
	}
	if filter.StartDate != nil {
		query = query.Where("performed_at >= ?", filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("performed_at <= ?", filter.EndDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("performed_at DESC").Limit(filter.Limit).Offset(filter.Offset).Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}
