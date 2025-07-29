package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"sort"
	"strings"
	"time"
)

func addSubm(db *sql.DB) {
	// TODO: implement submission insert logic here
}

// 3306
// Use the submittedAt timestamp to extract the date
func SameDaySubm(db *sql.DB, problemName, userID string, submittedAt string) bool {

	t, err := time.Parse("2006-01-02T15:04:05", submittedAt)
	if err != nil {
		log.Println("❌ Failed to parse SubmittedAt:", err)
		return true
	}

	dateOnly := t.Format("2006-01-02") // Strip to date only

	query := `
		SELECT EXISTS (
			SELECT 1 FROM submissions
			WHERE problem_name = ?
			  AND user_id = ?
			  AND DATE(timestamp) = ?
		)
	`
	var exists bool
	err = db.QueryRow(query, problemName, userID, dateOnly).Scan(&exists)
	if err != nil {
		log.Println("DB error:", err)
		return false // fail open
	}

	fmt.Println("🚫 Already submitted today?", exists)
	return exists
}

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
	username string,
	streak uint,
) error {

	stmt, err := db.Prepare(`
        INSERT INTO users (user_id, discord_user_id, username, discord_server_id, is_admin, monthly_leetcode, status, streak)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		return fmt.Errorf("Adduser perspare user failde: %v", err)
	}
	defer stmt.Close()

	fmt.Println(streak)

	_, err = stmt.Exec(userID, discordUserID, username, discordServerID, isAdmin, monthlyLeetcode, status, streak)
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
	increaseMonthlyLeetcode(db, userID, 1)
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
	Streak     uint   `json:"streak"`
}

func GetUsernameByUserID(db *sql.DB, userID string) (string, error) {
	query := `SELECT username FROM users WHERE user_id = ?`
	var username string

	err := db.QueryRow(query, userID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no username found for user ID %s", userID)
		}
		return "", fmt.Errorf("query failed: %v", err)
	}

	return username, nil
}

func GetUserIDwithDiscordID(db *sql.DB, discordUserID string) (string, error) {
	query := `SELECT user_id FROM users WHERE discord_user_id = ?`

	var userID string
	err := db.QueryRow(query, discordUserID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no user found for Discord ID %s", discordUserID)
		}
		return "", fmt.Errorf("query error: %v", err)
	}
	return userID, nil
}

func GetAllDailyLeets(db *sql.DB, date string) []DailyStat {
	t, err := time.Parse("2006-01-02T15:04:05", date)
	if err != nil {
		t, err = time.Parse("2006-01-02", date)
		if err != nil {
			log.Println("❌ Failed to parse date:", err)
			return nil
		}
	}
	date = t.Format("2006-01-02")

	query := `
		SELECT u.user_id, u.username, s.difficulty, COUNT(*)
		FROM submissions s
		JOIN users u ON s.user_id = u.user_id
		WHERE DATE(s.timestamp) = ?
		GROUP BY u.user_id, u.username, s.difficulty
	`
	//TODO: add steaks to this list incomplete

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
		var streak uint
		err := db.QueryRow("SELECT monthly_leetcode, streak FROM users WHERE user_id = ?", userID).Scan(&monthlyLC, &streak)

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
			Streak:     streak,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalToday > result[j].TotalToday
	})
	// n log(n)

	return result
}

func increaseMonthlyLeetcode(db *sql.DB, userID string, increase uint8) {
	updateQuery := `UPDATE users SET monthly_leetcode = monthly_leetcode + ? WHERE user_id = ?`
	stmt, err := db.Prepare(updateQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(increase, userID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("✅ Increased monthly leetcode for user %s by %d", userID, increase)
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

func IncrementStreak(db *sql.DB, userID string) error {
	updateQuery := `UPDATE users SET streak = streak + 1 WHERE user_id = ?`

	_, err := db.Exec(updateQuery, userID)
	if err != nil {
		return fmt.Errorf("failed to increment streak for user %s: %w", userID, err)
	}
	return nil
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

func DeleteUserByDiscordID(db *sql.DB, discordUserID string) error {
	var userID string
	err := db.QueryRow("SELECT user_id FROM users WHERE discord_user_id = ?", discordUserID).Scan(&userID)
	if err != nil {
		return fmt.Errorf("❌ Could not find user with Discord ID %s: %v", discordUserID, err)
	}

	// Delete submissions first to satisfy foreign key constraint
	_, err = db.Exec("DELETE FROM submissions WHERE user_id = ?", userID)
	if err != nil {
		return fmt.Errorf("❌ Failed to delete submissions for user %s: %v", userID, err)
	}

	// Then delete the user
	_, err = db.Exec("DELETE FROM users WHERE user_id = ?", userID)
	if err != nil {
		return fmt.Errorf("❌ Failed to delete user %s: %v", userID, err)
	}

	log.Printf("✅ Successfully deleted user %s and their submissions", userID)
	return nil
}

func GetRandomSkewedLeetcode(db *sql.DB, userID string) string {
	query := `SELECT problem_name, max(confidence_score) AS max_score
FROM submissions 
WHERE user_id = ?
GROUP BY problem_name
HAVING max_score < 5`

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(userID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var problemNames []string
	var confidenceScores []int8

	for rows.Next() {
		var pname string
		var score int8
		if err := rows.Scan(&pname, &score); err != nil {
			log.Fatal(err)
		}
		problemNames = append(problemNames, pname)
		confidenceScores = append(confidenceScores, score)
	}

	if len(problemNames) == 0 {
		log.Println("⚠️ No problems found for user.")
		return ""
	}

	var randomLeet string
	var score int8

	for {
		randomNum := rand.Intn(len(problemNames))
		randomLeet = problemNames[randomNum]
		score = confidenceScores[randomNum]
		rerollRandom := rand.Intn(100)

		fmt.Println(rerollRandom, score)
		fmt.Println(ScoreToProbability[score] + rerollRandom)
		if rerollRandom+ScoreToProbability[score] >= 100 {
			break
		}
	}

	url := fmt.Sprintf("https://leetcode.com/problems/%s", randomLeet)
	return url
}

var ScoreToProbability = map[int8]int{
	0: 100,
	1: 100,
	2: 90,
	3: 50,
	4: 2,
	5: 0,
}

func ResetMoLCA(db *sql.DB) {

	query := `INSERT INTO users (monthly_leetcode) VALUES 0`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

}

func QueryAllSuerActivity(db *sql.DB) map[string]bool {
	inactiveMap := make(map[string]bool)

	// Get latest submission timestamp for each user
	query := `
        SELECT username, user_id, MAX(timestamp) as latest
        FROM submissions
        GROUP BY user_id
    `
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	now := time.Now()

	for rows.Next() {
		var userID string
		var latestTime time.Time
		var username string

		if err := rows.Scan(&username, &userID, &latestTime); err != nil {
			log.Printf("❌ Error scanning submission timestamp for user %s: %v", userID, err)
			continue
		}

		if now.Sub(latestTime).Hours() > 96 { // 4 days
			inactiveMap[username] = true
		}
	}

	return inactiveMap
}
