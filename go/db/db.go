package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Required for mysql driver
	"github.com/joho/godotenv"
	"log"
	"os"
)

var DB *sql.DB // Global exported variable

func Init() {
	if os.Getenv("ENV") != "production" {
		fmt.Println("Loading environment variables from .env file")
		_ = godotenv.Load() // Only load .env locally
	}
	// load the env file from a tmeporatry test file durung tests

	dsn := os.Getenv("DNS_DB")

	fmt.Printf("dsn:", dsn[22:])

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to open DB: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping DB: %v", err)
	}

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("❌ Ping failed: %v", err)
	}

	DB.SetMaxOpenConns(20)
	DB.SetMaxIdleConns(20)

	log.Println("✅ Connected to DB")
}
