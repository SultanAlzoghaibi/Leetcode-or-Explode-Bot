package shared

import (
	"encoding/json"
	"fmt"
	"strings"
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

type Difficulty int8

const (
	Easy Difficulty = iota
	Medium
	Hard
)

func (d *Difficulty) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch strings.ToUpper(strings.TrimSpace(s)) {
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
