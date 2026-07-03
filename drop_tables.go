package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=127.0.0.1 user=gmao_user password=gmao_password dbname=gmao_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	
	err = db.Exec("DROP TABLE IF EXISTS inspections CASCADE").Error
	if err != nil {
		log.Fatalf("Failed to drop inspections: %v", err)
	}

	err = db.Exec("DROP TABLE IF EXISTS metric_measurements CASCADE").Error
	if err != nil {
		log.Fatalf("Failed to drop metric_measurements: %v", err)
	}

	err = db.Exec("DROP TABLE IF EXISTS interventions CASCADE").Error
	if err != nil {
		log.Fatalf("Failed to drop interventions: %v", err)
	}

	err = db.Exec("DROP TABLE IF EXISTS work_orders CASCADE").Error
	if err != nil {
		log.Fatalf("Failed to drop work_orders: %v", err)
	}

	log.Println("Successfully dropped tables")
}
