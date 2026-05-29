package postgres

import (
	"context"
	"fmt"
	"time"

	"backend-gmao/apps/asset-service/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type measurementRepository struct {
	db *gorm.DB
}

func NewMeasurementRepository(db *gorm.DB) *measurementRepository {
	return &measurementRepository{db: db}
}

// InitPartitionedTable creates the partitioned table and the current/next month partitions
// if they don't exist. This replaces GORM AutoMigrate for this specific table.
func InitPartitionedTable(db *gorm.DB) error {
	// Create master partitioned table
	// We use IF NOT EXISTS so it's idempotent.
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS measurements (
		id UUID,
		equipment_instance_id UUID,
		part_instance_id UUID,
		metric_name VARCHAR NOT NULL,
		value DOUBLE PRECISION NOT NULL,
		unit VARCHAR NOT NULL,
		recorded_at TIMESTAMP WITH TIME ZONE NOT NULL,
		recorded_by UUID,
		PRIMARY KEY (id, recorded_at)
	) PARTITION BY RANGE (recorded_at);
	`
	if err := db.Exec(createTableSQL).Error; err != nil {
		return fmt.Errorf("failed to create measurements partitioned table: %w", err)
	}

	// Create partitions for current and next month
	now := time.Now()
	
	monthsToCreate := []time.Time{
		now,
		now.AddDate(0, 1, 0),
	}

	for _, m := range monthsToCreate {
		startOfMonth := time.Date(m.Year(), m.Month(), 1, 0, 0, 0, 0, m.Location())
		startOfNextMonth := startOfMonth.AddDate(0, 1, 0)

		partitionName := fmt.Sprintf("measurements_y%04dm%02d", m.Year(), m.Month())
		
		createPartitionSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s PARTITION OF measurements 
		FOR VALUES FROM ('%s') TO ('%s');
		`, partitionName, startOfMonth.Format(time.RFC3339), startOfNextMonth.Format(time.RFC3339))

		if err := db.Exec(createPartitionSQL).Error; err != nil {
			return fmt.Errorf("failed to create partition %s: %w", partitionName, err)
		}
	}

	return nil
}

func (r *measurementRepository) CreateMeasurement(ctx context.Context, measurement *domain.Measurement) error {
	return r.db.WithContext(ctx).Create(measurement).Error
}

func (r *measurementRepository) GetMeasurementsByEquipment(ctx context.Context, equipmentID uuid.UUID, since time.Time) ([]domain.Measurement, error) {
	var measurements []domain.Measurement
	err := r.db.WithContext(ctx).
		Where("equipment_instance_id = ? AND recorded_at >= ?", equipmentID, since).
		Order("recorded_at desc").
		Find(&measurements).Error
	return measurements, err
}

func (r *measurementRepository) GetMeasurementsByPart(ctx context.Context, partID uuid.UUID, since time.Time) ([]domain.Measurement, error) {
	var measurements []domain.Measurement
	err := r.db.WithContext(ctx).
		Where("part_instance_id = ? AND recorded_at >= ?", partID, since).
		Order("recorded_at desc").
		Find(&measurements).Error
	return measurements, err
}
