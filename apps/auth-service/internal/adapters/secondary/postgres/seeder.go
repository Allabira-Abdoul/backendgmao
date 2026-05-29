package postgres

import (
	"log"

	"gorm.io/gorm"
)

// Seed populates initial session data if necessary (no-op).
func Seed(db *gorm.DB) {
	log.Println("Seeding auth service default data (no-op)")
}
