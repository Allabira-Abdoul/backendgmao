package postgres

import (
	"log"

	"gorm.io/gorm"
)

// Seed populates initial prediction data if necessary (no-op).
func Seed(db *gorm.DB) {
	log.Println("Seeding prediction service default data (no-op)")
}
