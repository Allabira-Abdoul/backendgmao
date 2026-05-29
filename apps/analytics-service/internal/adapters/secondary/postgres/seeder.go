package postgres

import (
	"log"

	"backend-gmao/apps/analytics-service/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Seed populates initial analytics metrics.
func Seed(db *gorm.DB) {
	var count int64
	db.Model(&domain.Metric{}).Count(&count)
	if count > 0 {
		return
	}

	log.Println("Seeding default metrics data...")
	initialMetrics := []domain.Metric{
		{
			ID:       uuid.New(),
			Name:     "mttr",
			Value:    120.5,
			Category: "MAINTENANCE",
		},
		{
			ID:       uuid.New(),
			Name:     "mtbf",
			Value:    4500.0,
			Category: "MAINTENANCE",
		},
	}

	for _, m := range initialMetrics {
		db.Create(&m)
	}
	log.Println("Seeding default metrics data completed")
}
