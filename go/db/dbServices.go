package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sort"
	"strings"
)

func addSubm(db *sql.DB) {
	// TODO: implement submission insert logic here
}

//3306

func DoesExist(db *sql.DB, table, column, key string) bool {

	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s = ?)", table, column)
	var exists bool
	err := db.QueryRow(query, key).Scan(&exists)
	if err != nil {
		log.Printf("❌ DB error checking %s.%s = %s: %v", table, column, key, err)
		return false
	}
	return exists
}

func AddUser(db *sql.DB,
	userID string,
	isAdmin bool,
	monthlyLeetcode uint8,
	status string,
	discordUserID string,
	discordServerID string,
	username string) error {

	stmt, err := db.Prepare(`
        INSERT INTO users (user_id, discord_user_id, username, discord_server_id, is_admin, monthly_leetcode, status)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		return fmt.Errorf("Adduser perspare user failde: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, discordUserID, username, discordServerID, isAdmin, monthlyLeetcode, status)
	if err != nil {
		return fmt.Errorf("execution failed: %v", err)
	}

	return nil
}

func AddSubm(db *sql.DB,
	submissionID string,
	problemName string,
	difficulty string,
	confidenceScore uint8,
	timestamp string,
	topics []string,
	solveTime uint8,
	notes string,
	userID string,
) error {
	topicsJSON, err := json.Marshal(topics)
	if err != nil {
		return fmt.Errorf("json marshal failed: %v", err)
	}

	stmt, err := db.Prepare(`
        INSERT INTO submissions (submission_id, 
                                 problem_name, 
                                 difficulty, 
                                 confidence_score, 
                                 timestamp, 
                                 topics,
                                 solve_time,
                                 notes, 
                                 user_id)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `)

	if err != nil {
		fmt.Println("Addsubm perspare submission failde: %v", err)
		return fmt.Errorf("execution failed: %v", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		submissionID,
		problemName,
		difficulty,
		confidenceScore,
		timestamp,
		string(topicsJSON),
		solveTime,
		notes,
		userID)
	if err != nil {
		fmt.Println("Addsubm exucute submission failde: %v", err)
		return fmt.Errorf("execution failed: %v", err)
	}
	log.Printf("✅ Inserted submission %s for user %s", submissionID, userID)
	return nil
}

func RemoveSubm(db *sql.DB, submissionID string, userID string) error {
	return nil
}
func DeleteRow(db *sql.DB, table, column, key string) {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", table, column)

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("❌ Prepare delete failed: %v", err)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(key)
	if err != nil {
		log.Printf("❌ Execution delete failed: %v", err)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("❌ Could not fetch rows affected: %v", err)
		return
	}
	log.Printf("✅ Deleted %d row(s) from %s", rowsAffected, table)

}

func PrintDB(db *sql.DB) {
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
	query2 := `SELECT id, problemNumber, confidenceScore, timestamp, userID FROM submissions`
	sRows, err := db.Query(query2)
	if err != nil {
		log.Fatal(err)
	}
	defer sRows.Close()

	fmt.Println("\nSubmissions:")
	for sRows.Next() {
		var id, userID string
		var problemNumber int
		var confidenceScore uint8
		var timestamp string

		err := sRows.Scan(&id, &problemNumber, &confidenceScore, &timestamp, &userID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("- ID: %s | Problem: %d | Score: %d | Time: %s | UserID: %s\n",
			id, problemNumber, confidenceScore, timestamp, userID)
	}
}

type DailyStat struct {
	UserID     string `json:"userID"`
	Username   string `json:"username"`
	Easy       int    `json:"easy"`
	Medium     int    `json:"medium"`
	Hard       int    `json:"hard"`
	TotalToday int    `json:"totalToday"`
	MonthlyLC  uint8  `json:"monthlyLeetcode"`
}

func GetAllDailyLeets(db *sql.DB, date string) []DailyStat {
	query := `
		SELECT u.user_id, u.username, s.difficulty, COUNT(*) 
		FROM submissions s
		JOIN users u ON s.user_id = u.user_id
		WHERE DATE(s.timestamp) = ?
		GROUP BY u.user_id, u.username, s.difficulty
	`

	rows, err := db.Query(query, date)
	if err != nil {
		log.Fatalf("❌ Query failed: %v", err)
	}
	defer rows.Close()

	type innerStat struct {
		Username   string
		Easy       int
		Medium     int
		Hard       int
		TotalToday int
	}

	userStats := make(map[string]*innerStat)

	for rows.Next() {
		var userID, username, difficulty string
		var count int

		if err := rows.Scan(&userID, &username, &difficulty, &count); err != nil {
			log.Fatalf("❌ Row scan failed: %v", err)
		}

		if _, exists := userStats[userID]; !exists {
			userStats[userID] = &innerStat{Username: username}
		}

		stat := userStats[userID]
		stat.TotalToday += count

		switch strings.ToUpper(difficulty) {
		case "EASY":
			stat.Easy += count
		case "MEDIUM":
			stat.Medium += count
		case "HARD":
			stat.Hard += count
		}
	}

	var result []DailyStat
	for userID, s := range userStats {
		var monthlyLC uint8
		err := db.QueryRow("SELECT monthly_leetcode FROM users WHERE user_id = ?", userID).Scan(&monthlyLC)
		if err != nil {
			log.Printf("⚠️ Could not get monthly LC for user %s: %v", userID, err)
		}

		result = append(result, DailyStat{
			UserID:     userID,
			Username:   s.Username,
			Easy:       s.Easy,
			Medium:     s.Medium,
			Hard:       s.Hard,
			TotalToday: s.TotalToday,
			MonthlyLC:  monthlyLC,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalToday > result[j].TotalToday
	})

	return result
}

func increaseMonthlyLeetcode(db *sql.DB, userID string, increase uint8) {
	query := `SELECT monthly_leetcode FROM users WHERE user_id = ?`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var moLCAmount uint8
	err = stmt.QueryRow(userID).Scan(&moLCAmount)
	if err != nil {
		log.Fatal(err)
	}
	moLCAmount += increase

	updateQuery := `UPDATE users SET monthly_leetcode = ? WHERE user_id = ?`
	updateStmt, err := db.Prepare(updateQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer updateStmt.Close()

	_, err = updateStmt.Exec(moLCAmount, userID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✅ Increased monthly leetcode for user %s: %d\n", userID, moLCAmount)
}

func getMoLCColumn(db *sql.DB, userID string) {
	query := `SELECT monthly_leetcode from users WHERE user_id = ?`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var moLCAmount uint8
	err = stmt.QueryRow(userID).Scan(&moLCAmount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("moLC column: %d \n", moLCAmount)
}

func resetMoLC(db *sql.DB) {
	query := `INSERT INTO submissions (monthly_leetcode) VALUES 0`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("monthy LC rest to 0")
}

type LeaderEntry struct {
	Username   string
	MoLCAmount uint8
}

func GetLeaderboard(db *sql.DB) []LeaderEntry {
	query := `SELECT username, monthly_leetcode FROM users ORDER BY monthly_leetcode DESC`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var leaderboard []LeaderEntry

	for i := 0; i < 3 && rows.Next(); i++ {
		var username string
		var moLCAmount uint8
		err = rows.Scan(&username, &moLCAmount)
		if err != nil {
			log.Fatal(err)
		}
		leaderboard = append(leaderboard, LeaderEntry{
			Username:   username,
			MoLCAmount: moLCAmount,
		})
	}

	return leaderboard
}
