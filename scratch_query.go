package main

import (
	"fmt"
	"log"

	"backend-gmao/pkg/db"
)

func main() {
	connStr := "host=127.0.0.1 user=gmao_user password=gmao_password dbname=gmao_db port=5432 sslmode=disable"
	database, err := db.InitPostgres(connStr)
	if err != nil {
		log.Fatal(err)
	}

	type Result struct {
		ID          string
		ServiceName string
		Action      string
		Details     string
	}

	var results []Result
	if err := database.Raw("SELECT id, service_name, action, details FROM audit_logs").Scan(&results).Error; err != nil {
		log.Fatal(err)
	}

	for _, r := range results {
		fmt.Printf("ID: %s, Service: %s, Action: %s, Details: %s\n", r.ID, r.ServiceName, r.Action, r.Details)
	}
	fmt.Printf("Total audit logs: %d\n", len(results))
}
