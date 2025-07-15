package bot

import (
	"Leetcode-or-Explode-Bot/db"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Server struct {
}

type Message struct {
	Text string `json:"text"`
}
type Submission struct {
	UserID        string `json:"userID"`        // ex: "7syRMHE2MD"
	SubmissionID  string `json:"submissionId"`  // ex: "1696788684"
	ProblemNumber int    `json:"problemNumber"` // ex: 1   (Two-Sum)
	Difficulty    string `json:"difficulty"`    // "Easy" | "Medium" | "Hard"
	SubmittedAt   string `json:"submittedAt"`   // ISO-8601 timestamp
	// TODO: Add a confidence in REDOING SCORE
}

func lcSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // or set specific origin
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != "POST" {
		print("POST REQ invalid")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(string(body))
	defer r.Body.Close()

	var submission Submission
	err = json.Unmarshal(body, &submission)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Printf("✅ Submission received:\n%+v\n", submission)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received"))

	// DB part
	database := db.DB

	subExists, err := tableExists(database, "submissions")
	userExists, err := tableExists(database, "users")

	if !subExists || !userExists {
		err := db.SetupDB(database)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("before")

}

func StartChromeAPIServer() {
	http.HandleFunc("/", lcSubmissionHandler)
	http.ListenAndServe(":9100", nil)
}

func tableExists(db *sql.DB, tableName string) (bool, error) {

	query := `	SELECT COUNT(*) 
				FROM INFORMATION_SCHEMA.TABLES 
				WHERE table_schema = DATABASE() 
				  AND table_name = ?`
	var count int
	err := db.QueryRow(query, tableName).Scan(&count)
	if err != nil {
		fmt.Println("error in counting DBS")
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
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
	query2 := `SELECT id, username, problem, score, timestamp FROM submissions`
	sRows, err := db.Query(query2)
	if err != nil {
		log.Fatal(err)
	}
	defer sRows.Close()

	fmt.Println("\nSubmissions:")
	for sRows.Next() {
		var id int
		var username, problem string
		var score int
		var timestamp string // you can also use time.Time if your column is DATETIME

		err := sRows.Scan(&id, &username, &problem, &score, &timestamp)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("- #%d | %s | %s | Score: %d | %s\n", id, username, problem, score, timestamp)
	}
}
