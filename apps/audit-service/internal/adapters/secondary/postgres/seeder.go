package postgres

import (
	"log"

	"gorm.io/gorm"
)

// Seed populates initial audit log data if necessary (no-op as logs represent active runtime logs).
func Seed(db *gorm.DB) {
	log.Println("Seeding audit service default data (no-op)")
}
