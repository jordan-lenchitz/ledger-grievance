package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "ledger")
	dbPass := getEnv("DB_PASSWORD", "ledgerpass")
	dbName := getEnv("DB_NAME", "grievance_ledger")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}
	defer db.Close()

	reporters := []string{"jordan", "admin", "system", "user_123", "manager_bot"}
	subjects := []string{"Unfair feedback", "Broken elevator", "Slow response", "Coffee was cold", "Meeting ran late"}
	categories := []string{"communication", "facility", "infrastructure", "etiquette"}

	fmt.Println("Seeding database with 100 random incidents...")

	for i := 0; i < 100; i++ {
		reporter := reporters[rand.Intn(len(reporters))]
		subject := subjects[rand.Intn(len(subjects))]
		category := categories[rand.Intn(len(categories))]
		severity := rand.Intn(5) + 1
		description := fmt.Sprintf("This is a generated grievance number %d. It reflects a systemic issue that needs immediate attention.", i)
		accommodation := rand.Float32() < 0.2

		_, err := db.Exec(`
			INSERT INTO incidents 
			(reporter_id, subject, category, severity, description, requires_accommodation) 
			VALUES (?, ?, ?, ?, ?, ?)`,
			reporter, subject, category, severity, description, accommodation,
		)
		if err != nil {
			log.Printf("Error inserting incident %d: %v", i, err)
		}
	}

	fmt.Println("Seeding complete.")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
