package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	dsn := "host=localhost port=5432 user=postgres password=123456 dbname=cal_salary sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM payroll_records")
	if err != nil {
		log.Fatalf("Error counting payroll_records: %v", err)
	}

	fmt.Printf("Total rows in payroll_records: %d\n", count)

	var detailsCount int
	err = db.Get(&detailsCount, "SELECT COUNT(*) FROM payroll_record_details")
	if err != nil {
		log.Fatalf("Error counting payroll_record_details: %v", err)
	}
	fmt.Printf("Total rows in payroll_record_details: %d\n", detailsCount)
}
