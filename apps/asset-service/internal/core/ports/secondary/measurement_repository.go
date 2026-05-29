package secondary

import (
	"context"
	"time"

	"backend-gmao/apps/asset-service/internal/core/domain"
	"github.com/google/uuid"
)

type MeasurementRepository interface {
	CreateMeasurement(ctx context.Context, measurement *domain.Measurement) error
	GetMeasurementsByEquipment(ctx context.Context, equipmentID uuid.UUID, since time.Time) ([]domain.Measurement, error)
	GetMeasurementsByPart(ctx context.Context, partID uuid.UUID, since time.Time) ([]domain.Measurement, error)
}
