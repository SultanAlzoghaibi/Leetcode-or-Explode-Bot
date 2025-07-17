package db

import (
	"database/sql"
	"fmt"
	"log"
)

func SetupDB(db *sql.DB) error {

	submisionTable := `
	CREATE TABLE IF NOT EXISTS submissions (
	    submission_id VARCHAR(32) PRIMARY KEY,
	    problem_name VARCHAR(80) NOT NULL,
	    difficulty ENUM('EASY', 'MEDIUM', 'HARD'),
	    confidence_score TINYINT,
	    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
	    topics TEXT, -- array of strings
	    solve_time TINYINT UNSIGNED,
	    notes TEXT,
	    user_id VARCHAR(32),
	    FOREIGN KEY (user_id) REFERENCES users(user_id)
	);`
	userTable := `
CREATE TABLE IF NOT EXISTS users (
    user_id           VARCHAR(32) PRIMARY KEY,
    discord_user_id   VARCHAR(32) NOT NULL UNIQUE,
    discord_server_id VARCHAR(32) NOT NULL,
    is_admin          BOOLEAN,
    monthly_leetcode  TINYINT UNSIGNED,
    status            ENUM('NONE', 'DEFAULT', 'NO-PING', 'EXTREME')
)`

	// Execute user table creation
	_, err := db.Exec(userTable)
	if err != nil {
		return err
	}

	// Execute submissions table creation
	_, err = db.Exec(submisionTable)
	if err != nil {
		return err
	}

	return nil
}

func printDB(db *sql.DB) {
	// --- Print Users ---
	query1 := `SELECT user_id, discord_user_id, discord_server_id, is_admin, monthly_leetcode, status FROM users`
	uRows, err := db.Query(query1)
	if err != nil {
		log.Fatal(err)
	}
	defer uRows.Close()

	fmt.Println("\nUsers:")
	for uRows.Next() {
		var userID, discordUserID, discordServerID, status string
		var isAdmin bool
		var monthlyLeetcode uint8

		err := uRows.Scan(&userID, &discordUserID, &discordServerID, &isAdmin, &monthlyLeetcode, &status)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("- %s | %s | %s | Admin: %v | Monthly: %d | Status: %s\n",
			userID, discordUserID, discordServerID, isAdmin, monthlyLeetcode, status)
	}

	// --- Print Submissions ---
	query2 := `SELECT submission_id, problem_name, difficulty, confidence_score, timestamp, topics, solve_time, notes, user_id FROM submissions`
	sRows, err := db.Query(query2)
	if err != nil {
		log.Fatal(err)
	}
	defer sRows.Close()

	fmt.Println("\nSubmissions:")
	for sRows.Next() {
		var id, problemName, difficulty, submittedAt, topics, notes, userID string
		var confidenceScore, solveTime uint8

		err := sRows.Scan(&id, &problemName, &difficulty, &confidenceScore, &submittedAt, &topics, &solveTime, &notes, &userID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("- %s | %s | %s | Conf: %d | Time: %dmin | %s | Topics: %s | Notes: %s\n",
			id, userID, problemName, confidenceScore, solveTime, submittedAt, topics, notes)
	}
}
