package bot

import (
	"Leetcode-or-Explode-Bot/db"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Server struct {
}

type Message struct {
	Text string `json:"text"`
}

type Difficulty int8

const (
	Easy Difficulty = iota
	Medium
	Hard
)

type Submission struct {
	SubmissionID    string   `json:"submissionId"` // ex: "1696788684"
	ProblemName     string   `json:"problemName"`  // ex: 1   (Two-Sum)
	UserID          string   `json:"userID"`       // ex: "7syRMHE2MD"
	Difficulty      string   `json:"difficulty"`   // "Easy" | "Medium" | "Hard" / ENUM
	SubmittedAt     string   `json:"submittedAt"`  // ISO-8601 timestamp
	ConfidenceScore uint8    `json:"confidenceScore"`
	Notes           string   `json:"notes"`
	SolveTime       uint8    `json:"duration"`
	Topics          []string `json:"topics"`
}

func (d *Difficulty) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch strings.ToUpper(s) {
	case "EASY":
		*d = Easy
	case "MEDIUM":
		*d = Medium
	case "HARD":
		*d = Hard
	default:
		return fmt.Errorf("invalid difficulty: %s", s)
	}
	return nil
}

func lcSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("📩 we got a submission")

	validOrigins := map[string]bool{
		"chrome-extension://bphfdocncclgepoiabbjodikpeegopfd": true,
		"chrome-extension://lffamldlgnnlimjpggjcphdocgjbflna": true,
	}

	origin := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Vary", "Origin")
	if validOrigins[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	if r.Method == "OPTIONS" {
		if validOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
		return
	}

	if !validOrigins[origin] {
		log.Printf("❌ Rejected CORS request from origin: %s\n", origin)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("📦 Raw body:", string(body))

	var submission Submission

	err = json.Unmarshal(body, &submission)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("✅ Submission received:\n%+v\n", submission)

	database := db.DB

	subExists, err := tableExists(database, "submissions")
	if err != nil {
		log.Fatal(err)
	}
	userExists, err := tableExists(database, "users")
	if err != nil {
		log.Fatal(err)
	}

	if !subExists || !userExists {
		err := db.SetupDB(database)
		if err != nil {
			log.Fatal(err)
		}
	}
	//TODO: pu this into a function as its messy out here

	//Todo: have a check that if the User is not in the DB, pop up warning is called

	if !db.DoesExist(database, "users", "user_id", submission.UserID) {
		fmt.Println("!db.DoesExist(database, \"submissions\", \"user_id\", submission.UserID)")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("the user has never signed up via our discord bot, contact h82luzn on discord for more information"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received"))

	if db.SameDaySubm(database, submission.ProblemName, submission.UserID, submission.SubmittedAt) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("the submission already exists today"))
		fmt.Println("submission was already submited today")
		return
	}

	db.AddSubm(database,
		submission.SubmissionID,
		submission.ProblemName,
		submission.Difficulty,
		submission.ConfidenceScore,
		submission.SubmittedAt,
		submission.Topics,
		submission.SolveTime,
		submission.Notes,
		submission.UserID)

	//printDB(database)

	addtoSheets(submission)

}

func StartChromeAPIServer() {
	http.HandleFunc("/api/chrome", lcSubmissionHandler)

	if err := http.ListenAndServe(":9100", nil); err != nil {
		log.Fatalf("❌ Failed to start Chrome API server: %v", err)
	}
	fmt.Println("✅ Chrome API server listening on port 9100")
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
	query2 := `SELECT submission_id, problem_number, confidence_score, timestamp, notes, user_id FROM submissions`
	sRows, err := db.Query(query2)
	if err != nil {
		log.Fatal(err)
	}
	defer sRows.Close()

	fmt.Println("\nSubmissions:")
	for sRows.Next() {
		var id string
		var problemNumber uint16
		var confidenceScore uint8
		var submittedAt, userID, notes string

		err := sRows.Scan(&id, &problemNumber, &confidenceScore, &submittedAt, &userID, &notes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("- #%d | Q#%d | Conf: %d | %s | User: %s | Notes: %s\n",
			id, problemNumber, confidenceScore, submittedAt, userID, notes)
	}
}
