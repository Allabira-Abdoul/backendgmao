package postgres

import (
	"log"

	"gorm.io/gorm"
)

// Seed populates initial maintenance work orders if necessary (no-op).
func Seed(db *gorm.DB) {
	log.Println("Seeding maintenance service default data (no-op)")
}
