package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Required for mysql driver
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

var DB *sql.DB // Global exported variable

func Init() {
	cwd, _ := os.Getwd()
	fmt.Println("📂 Current working dir:", cwd)

	if os.Getenv("ENV") != "production" {
		fmt.Println("Loading environment variables from .env file")
		if err := godotenv.Load(".env"); err != nil {
			log.Println("❌ Failed to load .env:", err)
		} else {
			fmt.Println("✅ .env file loaded")
		}
	}
	// load the env file from a tmeporatry test file durung tests

	dsn := os.Getenv("DNS_DB")
	fmt.Println("dns:", dsn)

	fmt.Printf("dsn:", dsn[22:])

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to open DB: %v", err)
	}

	// Ping with a short timeout; fail fast if unreachable
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("❌ Failed to ping DB: %v", err)
	}

	// Connection pool tuning to avoid stale idle sockets and broken pipes
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Promote to global
	DB = db
	log.Println("✅ Connected to DB")
}
