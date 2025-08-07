package chrome

import (
	"Leetcode-or-Explode-Bot/internal/chrome/cybersec"
	db2 "Leetcode-or-Explode-Bot/internal/db"
	"Leetcode-or-Explode-Bot/internal/shared"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Server struct {
}

type Message struct {
	Text string `json:"text"`
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

	var submission shared.Submission

	err = json.Unmarshal(body, &submission)
	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = forwarded
	}
	if !cybersec.Blackbox(
		submission,    // from JSON body
		r.Header,      // all request headers
		ip,            // extracted client IP
		http.Client{}, // reusable HTTP client for outbound validation calls
		r.UserAgent(), // User-Agent string
		origin,        // CORS origin
	) {
		http.Error(w, "Sorry lil bro, no ctf kids cracking this lol", http.StatusForbidden)
		return
	}
	fmt.Println("CYBERPASS")

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("✅ Submission received:\n%+v\n", submission)

	database := db2.DB

	subExists, err := tableExists(database, "submissions")
	if err != nil {
		log.Fatal(err)
	}
	userExists, err := tableExists(database, "users")
	if err != nil {
		log.Fatal(err)
	}

	if !subExists || !userExists {
		err := db2.SetupDB(database)
		if err != nil {
			log.Fatal(err)
		}
	}
	//TODO: put this into a function as its messy out here

	if !validSubm(submission) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !db2.DoesExist(database, "users", "user_id", submission.UserID) {
		fmt.Println("!db.DoesExist(database, \"submissions\", \"user_id\", submission.UserID)")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("the user has never signed up via our discord bot, contact h82luzn on discord for more information"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received"))

	if db2.SameDaySubm(database, submission.ProblemName, submission.UserID, submission.SubmittedAt) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("the submission already exists today"))
		fmt.Println("submission was already submited today")
		return
	}

	db2.AddSubm(database,
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
	shared.AddtoSheets(submission)

}

func validSubm(submission shared.Submission) bool {
	// ✅ ProblemName must contain at least one dash (e.g. "1234-two-sum")
	if len(submission.ProblemName) < 5 || !strings.Contains(submission.ProblemName, "-") {
		fmt.Println("❌ Invalid ProblemName: must be at least 5 characters and contain a dash (e.g. '1234-two-sum')")
		return false
	}

	// ✅ ProblemName must follow format like "1234-two-sum"
	match, _ := regexp.MatchString(`^\d+-[a-zA-Z0-9\-]+$`, submission.ProblemName)
	if !match {
		fmt.Println("❌ Invalid ProblemName format: must start with digits followed by a dash and a slug (e.g. '1234-two-sum')")
		return false
	}

	// ✅ ConfidenceScore must be between 1 and 10
	if submission.ConfidenceScore < 0 || submission.ConfidenceScore > 10 {
		fmt.Println("❌ Invalid ConfidenceScore: must be between 0 and 10")
		return false
	}

	// ✅ SubmissionID must be non-empty
	if strings.TrimSpace(submission.SubmissionID) == "" {
		fmt.Println("❌ SubmissionID is empty")
		return false
	}

	// ✅ SolveTime (duration) must be a positive number (minutes)
	if submission.SolveTime < 0 {
		fmt.Println("❌ SolveTime is negative: must be a positive number representing minutes spent")
		return false
	}

	return true
}

func StartChromeAPIServer() {
	fmt.Println("StartChromeAPIServer")
	http.HandleFunc("/api/chrome", lcSubmissionHandler)

	if err := http.ListenAndServe(":9100", nil); err != nil {
		log.Fatalf("❌ Failed to start Chrome API server: %v", err)
	}

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
