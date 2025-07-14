package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
}

func startChromeAPIServer() {
	http.HandleFunc("/", lcSubmissionHandler)
	http.ListenAndServe(":9100", nil)
}
