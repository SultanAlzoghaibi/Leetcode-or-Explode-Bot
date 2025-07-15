package db

import "database/sql"

func SetupDB(db *sql.DB) error {

	submisionTable := `
	CREATE TABLE IF NOT EXISTS submissions (
		id VARCHAR(32) PRIMARY KEY,
		problemNumber INT,
		confidenceScore TINYINT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		userID VARCHAR(32),
		FOREIGN KEY (userID) REFERENCES users(user_id)
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
