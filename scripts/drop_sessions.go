package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=127.0.0.1 user=gmao_user password=gmao_password dbname=gmao_db port=5432 sslmode=disable TimeZone=UTC"
	if envDsn := os.Getenv("DATABASE_URL"); envDsn != "" {
		dsn = envDsn
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.Migrator().DropTable("sessions")
	if err != nil {
		log.Fatalf("failed to drop table: %v", err)
	}
	log.Println("Sessions table dropped successfully")
}
