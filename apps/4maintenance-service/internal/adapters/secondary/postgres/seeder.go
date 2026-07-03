package postgres

import (
	"log"
	"math/rand"
	"time"

	"backend-gmao/apps/maintenance-service/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	log.Println("Seeding maintenance service default data...")

	var count int64
	db.Model(&domain.OrdreTravail{}).Count(&count)
	if count > 0 {
		log.Println("Maintenance data already exists. Skipping seed.")
		return
	}

	// Try to get a user and an asset to assign the work orders to
	var userIdStr, assetIdStr string
	db.Table("users").Select("id").Limit(1).Pluck("id", &userIdStr)
	db.Table("equipment_instances").Select("id").Limit(1).Pluck("id", &assetIdStr)

	var userId, assetId uuid.UUID
	var assignedTo *uuid.UUID

	if userIdStr != "" {
		userId = uuid.MustParse(userIdStr)
		assignedTo = &userId
	} else {
		log.Println("Warning: No users found. 'AssignedTo' will be null.")
	}

	if assetIdStr != "" {
		assetId = uuid.MustParse(assetIdStr)
	} else {
		// Create a dummy asset directly via GORM map to avoid failing
		log.Println("Warning: No equipment instances found. Creating dummy asset.")
		assetId = uuid.New()
		db.Table("equipment_instances").Create(map[string]interface{}{
			"id":                 assetId,
			"code":               "DUMMY-001",
			"equipment_model_id": uuid.New(), // Note: will fail if equipment_models has FK, but we assume cascade or no constraints if it's missing
			"status":             "OPERATIONAL",
			"location":           "DUMMY LOCATION",
			"created_at":         time.Now(),
			"updated_at":         time.Now(),
		})
	}

	now := time.Now()
	// Generate 20 diverse work orders spanning the current month (-15 days to +15 days)
	priorities := []string{"LOW", "MEDIUM", "HIGH", "CRITICAL"}
	statuses := []string{"PENDING", "IN_PROGRESS", "COMPLETED"}
	types := []string{"INTERVENTION", "INSPECTION"}
	categories := []string{"CORRECTIVE", "PREVENTIVE"}
	mTypes := []string{"PALLIATIVE", "CURATIVE", "SYSTEMATIC", "CONDITIONAL", "PREDICTIVE"}

	for i := 1; i <= 20; i++ {
		daysOffset := rand.Intn(30) - 15 // -15 to +15 days
		scheduled := now.AddDate(0, 0, daysOffset)

		wo := domain.OrdreTravail{
			ID:                  uuid.New(),
			Title:               "Automated Scheduled Task #" + string(rune('A'+i)),
			Description:         "This is an auto-generated work order for testing all possibilities.",
			AssetID:             assetId,
			Priority:            priorities[rand.Intn(len(priorities))],
			Status:              statuses[rand.Intn(len(statuses))],
			Type:                types[rand.Intn(len(types))],
			ScheduledAt:         &scheduled,
			MaintenanceCategory: categories[rand.Intn(len(categories))],
			MaintenanceType:     mTypes[rand.Intn(len(mTypes))],
			IsMetricMeasurement: rand.Intn(2) == 0,
			AssignedTo:          assignedTo,
			CreatedAt:           now,
			UpdatedAt:           now,
		}
		
		if wo.Status == "COMPLETED" && wo.Type == "INTERVENTION" {
			started := scheduled.Add(time.Hour)
			ended := started.Add(time.Hour * 2)
			wo.Interventions = []domain.Intervention{
				{
					ID:                  uuid.New(),
					WorkOrderID:         wo.ID,
					Description:         "Completed intervention steps.",
					MaintenanceCategory: wo.MaintenanceCategory,
					MaintenanceType:     wo.MaintenanceType,
					IsMetricMeasurement: wo.IsMetricMeasurement,
					StartedAt:           &started,
					EndedAt:             &ended,
					PerformedBy:         userId, // Requires user to exist, else could fail if FK strict
					CreatedAt:           now,
					UpdatedAt:           now,
				},
			}
		}

		db.Create(&wo)
	}

	log.Println("Successfully seeded 20 diverse work orders.")
}
