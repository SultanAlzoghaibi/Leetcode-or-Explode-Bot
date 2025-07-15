package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // Required for mysql driver
	"github.com/joho/godotenv"
	"log"
	"os"
)

var DB *sql.DB // Global exported variable

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	dsn := os.Getenv("DNS_DB")

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("❌ Ping failed: %v", err)
	}

	log.Println("✅ Connected to DB")
}
