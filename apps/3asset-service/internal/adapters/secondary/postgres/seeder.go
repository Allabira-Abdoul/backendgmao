package postgres

import (
	"log"

	"gorm.io/gorm"
)

// SeedData seeds the database with initial data for the 5-level hierarchy.
func SeedData(db *gorm.DB) {
	log.Println("Seeding Asset database started...")
	// TODO: Implement seeding for Site -> System -> Asset -> Subsystem -> Component
	log.Println("Seeding Asset database completed")
}
