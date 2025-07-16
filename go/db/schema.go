package db

import "database/sql"

func SetupDB(db *sql.DB) error {

	submisionTable := `
	CREATE TABLE IF NOT EXISTS submissions (
	    submission_id VARCHAR(32) PRIMARY KEY,
	    problem_number INT,
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
