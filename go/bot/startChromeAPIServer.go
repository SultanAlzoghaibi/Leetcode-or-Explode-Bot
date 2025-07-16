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
	UserID          string `json:"userID"`        // ex: "7syRMHE2MD"
	SubmissionID    string `json:"submissionId"`  // ex: "1696788684"
	ProblemNumber   uint16 `json:"problemNumber"` // ex: 1   (Two-Sum)
	Difficulty      string `json:"difficulty"`    // "Easy" | "Medium" | "Hard"
	SubmittedAt     string `json:"submittedAt"`   // ISO-8601 timestamp
	ConfidenceScore uint8  `json:"confidenceScore"`
	Notes           string `json:"notes"`
	// TODO: Add a confidence in REDOING SCOREs
}

func lcSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodPost {
		log.Println("❌ Rejected non-POST request")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(string(body))

	var submission Submission
	err = json.Unmarshal(body, &submission)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("✅ Submission received:\n%+v\n", submission)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received"))

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

	printDB(database)
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
	query2 := `SELECT id, problemNumber, confidenceScore, timestamp, userID, notes FROM submissions`
	sRows, err := db.Query(query2)
	if err != nil {
		log.Fatal(err)
	}
	defer sRows.Close()

	fmt.Println("\nSubmissions:")
	for sRows.Next() {
		var id int
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
