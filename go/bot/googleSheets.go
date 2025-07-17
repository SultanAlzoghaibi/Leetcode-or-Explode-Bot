// google_sheets.go
package bot

import (
	"context"
	"encoding/json"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Define spreadsheetID and writeRange globally or pass as needed.
var spreadsheetID = "1Gc3PhSnLSrlcVSEDtQFiQ-rS-pHjWrgQQ64GFR-dwFQ" // TODO: replace with your spreadsheet ID
var writeRange = "Sheet1!A1:I"                                     // TODO: replace with your target range

func addtoSheets(subm Submission) {
	ctx := context.Background()

	// Create Sheets service using service account credentials
	srv, err := sheets.NewService(ctx,
		option.WithCredentialsFile("go/credentials.json"),
		option.WithScopes(sheets.SpreadsheetsScope),
	)
	if err != nil {
		log.Fatalf("Unable to create Sheets service: %v", err)
	}

	// Marshal topics to JSON string
	topicsJSON, _ := json.Marshal(subm.Topics)

	// Create the row to append
	row := []interface{}{
		subm.SubmissionID,
		subm.UserID,
		subm.ProblemNumber,
		subm.Difficulty.String(),
		subm.ConfidenceScore,
		subm.SolveTime,
		subm.Notes,
		string(topicsJSON),
	}

	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{row},
	}

	// Append to spreadsheet
	_, err = srv.Spreadsheets.Values.Append(spreadsheetID, writeRange, valueRange).
		ValueInputOption("RAW").
		Context(ctx).
		Do()
	if err != nil {
		log.Fatalf("Unable to append data to sheet: %v", err)
	}
}

/* ---------- helpers for OAuth & Difficulty ---------- */

// tokenFromFile / getTokenFromWeb / saveToken are the standard helper
// functions shown in Google’s Go Sheets quick-start.
// (Paste them unchanged from https://developers.google.com/sheets/api/quickstart/go)

func (d Difficulty) String() string {
	switch d {
	case Easy:
		return "Easy"
	case Medium:
		return "Medium"
	case Hard:
		return "Hard"
	default:
		return "Unknown"
	}
}
