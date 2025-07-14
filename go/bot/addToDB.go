package bot

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func submisionAddtoDB(db sql.DB, subm Submission) {

	stmt, err := db.Prepare("INSERT INTO difficulty (submission) VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

}

func SetDB(db *sql.DB) {

	if err := db.Ping(); err != nil {
		if setupErr := setupDB(db); setupErr != nil {
			fmt.Println("Error setting up database", setupErr)
			// optionally log or handle setup error
		}
	}
	fmt.Println("WE Are running the DB")
}

func setupDB(db *sql.DB) error {

	submisionTable := `
	CREATE TABLE IF NOT EXISTS submissions (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255),
		problem VARCHAR(255),
		score INT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	userTable := `
CREATE TABLE IF NOT EXISTS users (
    user_id           VARCHAR PRIMARY KEY   -- from LeetCode
	discord_user_id   VARCHAR NOT NULL      -- for @ mentions in Discord
	discord_server_id VARCHAR NOT NULL      -- for ping context
	is_admin          BOOLEAN
	monthly_leetcode  TINYINT UNSIGNED
	status            ENUM('NONE', 'DEFAULT', 'NO-PING', 'EXTREME')
)`

	if _, err := db.Exec(submisionTable); err != nil {
		return err
	}

	if _, err := db.Exec(userTable); err != nil {
		return err
	}
	return nil

}
