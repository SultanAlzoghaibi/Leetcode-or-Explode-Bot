// google_sheets.go
package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

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

	setDifficultyValidationAndFormatting(srv, spreadsheetID)
	setConfidenceValidationAndFormatting(srv, spreadsheetID)

	// --- Set data validation for "Difficulty" column (column B) ---

	// Marshal topics to JSON string
	topicsJSON, _ := json.Marshal(subm.Topics)

	// Create the row to append
	row := []interface{}{
		fmt.Sprintf(`=HYPERLINK("https://leetcode.com/problems/%s/", "%s")`,
			strings.SplitN(subm.ProblemName, "-", 2)[1], // slug
			subm.ProblemName),

		subm.Difficulty.String(),
		subm.ConfidenceScore,

		fmt.Sprintf(`=TEXT(DATEVALUE("%s"), "yyyy-mm-dd")`, subm.SubmittedAt[:10]),
		fmt.Sprintf("%dmin%s", subm.SolveTime, func() string {
			if subm.SolveTime >= 60 {
				return fmt.Sprintf(" (%dh %dmin)", subm.SolveTime/60, subm.SolveTime%60)
			}
			return ""
		}()),
		string(topicsJSON),
		subm.Notes,
	}

	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{row},
	}

	// Append to spreadsheet
	_, err = srv.Spreadsheets.Values.Append(spreadsheetID, writeRange, valueRange).
		ValueInputOption("USER_ENTERED").
		Context(ctx).
		Do()
	if err != nil {
		log.Fatalf("Unable to append data to sheet: %v", err)
	}
}

func setDifficultyValidationAndFormatting(srv *sheets.Service, spreadsheetID string) {
	_, err := srv.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				SetDataValidation: &sheets.SetDataValidationRequest{
					Range: &sheets.GridRange{
						SheetId:          0, // Update if not the first sheet
						StartColumnIndex: 1,
						EndColumnIndex:   2,
					},
					Rule: &sheets.DataValidationRule{
						Condition: &sheets.BooleanCondition{
							Type: "ONE_OF_LIST",
							Values: []*sheets.ConditionValue{
								{UserEnteredValue: "Easy"},
								{UserEnteredValue: "Medium"},
								{UserEnteredValue: "Hard"},
							},
						},
						Strict:       true,
						ShowCustomUi: true,
					},
				},
			},
			{
				AddConditionalFormatRule: &sheets.AddConditionalFormatRuleRequest{
					Rule: &sheets.ConditionalFormatRule{
						Ranges: []*sheets.GridRange{{
							SheetId:          0,
							StartColumnIndex: 1,
							EndColumnIndex:   2,
						}},
						BooleanRule: &sheets.BooleanRule{
							Condition: &sheets.BooleanCondition{
								Type:   "TEXT_EQ",
								Values: []*sheets.ConditionValue{{UserEnteredValue: "Easy"}},
							},
							Format: &sheets.CellFormat{
								BackgroundColor: &sheets.Color{Red: 0.8, Green: 1.0, Blue: 0.8},
							},
						},
					},
					Index: 0,
				},
			},
			{
				AddConditionalFormatRule: &sheets.AddConditionalFormatRuleRequest{
					Rule: &sheets.ConditionalFormatRule{
						Ranges: []*sheets.GridRange{{
							SheetId:          0,
							StartColumnIndex: 1,
							EndColumnIndex:   2,
						}},
						BooleanRule: &sheets.BooleanRule{
							Condition: &sheets.BooleanCondition{
								Type:   "TEXT_EQ",
								Values: []*sheets.ConditionValue{{UserEnteredValue: "Medium"}},
							},
							Format: &sheets.CellFormat{
								BackgroundColor: &sheets.Color{Red: 1.0, Green: 0.8, Blue: 0.4},
							},
						},
					},
					Index: 0,
				},
			},
			{
				AddConditionalFormatRule: &sheets.AddConditionalFormatRuleRequest{
					Rule: &sheets.ConditionalFormatRule{
						Ranges: []*sheets.GridRange{{
							SheetId:          0,
							StartColumnIndex: 1,
							EndColumnIndex:   2,
						}},
						BooleanRule: &sheets.BooleanRule{
							Condition: &sheets.BooleanCondition{
								Type:   "TEXT_EQ",
								Values: []*sheets.ConditionValue{{UserEnteredValue: "Hard"}},
							},
							Format: &sheets.CellFormat{
								BackgroundColor: &sheets.Color{Red: 1.0, Green: 0.8, Blue: 0.8},
							},
						},
					},
					Index: 0,
				},
			},
		},
	}).Do()
	if err != nil {
		log.Fatalf("Unable to set data validation: %v", err)
	}
}

func setConfidenceValidationAndFormatting(srv *sheets.Service, spreadsheetID string) {
	_, err := srv.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			// Set data validation
			{
				SetDataValidation: &sheets.SetDataValidationRequest{
					Range: &sheets.GridRange{
						SheetId:          0, // Update if not the first sheet
						StartColumnIndex: 2,
						EndColumnIndex:   3,
					},
					Rule: &sheets.DataValidationRule{
						Condition: &sheets.BooleanCondition{
							Type: "ONE_OF_LIST",
							Values: []*sheets.ConditionValue{
								{UserEnteredValue: "0 – No clue"},
								{UserEnteredValue: "1 – Struggle to repeat"},
								{UserEnteredValue: "2 – Might redo poorly"},
								{UserEnteredValue: "3 – Could redo maybe"},
								{UserEnteredValue: "4 – Confident redo"},
								{UserEnteredValue: "5 – Perfectly repeatable"},
							},
						},
						Strict:       true,
						ShowCustomUi: true,
					},
				},
			},
			// Add conditional formatting for "0 – No clue" (Dark Red)
			{
				AddConditionalFormatRule: &sheets.AddConditionalFormatRuleRequest{
					Rule: &sheets.ConditionalFormatRule{
						Ranges: []*sheets.GridRange{
							{
								SheetId:          0,
								StartColumnIndex: 2,
								EndColumnIndex:   3,
							},
						},
						BooleanRule: &sheets.BooleanRule{
							Condition: &sheets.BooleanCondition{
								Type: "TEXT_EQ",
								Values: []*sheets.ConditionValue{
									{UserEnteredValue: "0 – No clue"},
								},
							},
							Format: &sheets.CellFormat{
								BackgroundColor: &sheets.Color{
									Red:   0.8,
									Green: 0.0,
									Blue:  0.0,
									Alpha: 1.0,
								},
							},
						},
					},
					Index: 0,
				},
			},
			// Add conditional formatting for "1 – Struggle to repeat" (Red-Orange)
			{
				AddConditionalFormatRule: &sheets.AddConditionalFormatRuleRequest{
					Rule: &sheets.ConditionalFormatRule{
						Ranges: []*sheets.GridRange{
							{
								SheetId:          0,
								StartColumnIndex: 2,
								EndColumnIndex:   3,
							},
						},
						BooleanRule: &sheets.BooleanRule{
							Condition: &sheets.BooleanCondition{
								Type: "TEXT_EQ",
								Values: []*sheets.ConditionValue{
									{UserEnteredValue: "1 – Struggle to repeat"},
								},
							},
							Format: &sheets.CellFormat{
								BackgroundColor: &sheets.Color{
									Red:   0.9,
									Green: 0.3,
									Blue:  0.0,
									Alpha: 1.0,
								},
							},
						},
					},
					Index: 1,
				},
			},
			// Add conditional formatting for "2 – Might redo poorly" (Orange)
			{
				AddConditionalFormatRule: &sheets.AddConditionalFormatRuleRequest{
					Rule: &sheets.ConditionalFormatRule{
						Ranges: []*sheets.GridRange{
							{
								SheetId:          0,
								StartColumnIndex: 2,
								EndColumnIndex:   3,
							},
						},
						BooleanRule: &sheets.BooleanRule{
							Condition: &sheets.BooleanCondition{
								Type: "TEXT_EQ",
								Values: []*sheets.ConditionValue{
									{UserEnteredValue: "2 – Might redo poorly"},
								},
							},
							Format: &sheets.CellFormat{
								BackgroundColor: &sheets.Color{
									Red:   1.0,
									Green: 0.6,
									Blue:  0.0,
									Alpha: 1.0,
								},
							},
						},
					},
					Index: 2,
				},
			},
			// Add conditional formatting for "3 – Could redo maybe" (Yellow)
			{
				AddConditionalFormatRule: &sheets.AddConditionalFormatRuleRequest{
					Rule: &sheets.ConditionalFormatRule{
						Ranges: []*sheets.GridRange{
							{
								SheetId:          0,
								StartColumnIndex: 2,
								EndColumnIndex:   3,
							},
						},
						BooleanRule: &sheets.BooleanRule{
							Condition: &sheets.BooleanCondition{
								Type: "TEXT_EQ",
								Values: []*sheets.ConditionValue{
									{UserEnteredValue: "3 – Could redo maybe"},
								},
							},
							Format: &sheets.CellFormat{
								BackgroundColor: &sheets.Color{
									Red:   1.0,
									Green: 0.8,
									Blue:  0.0,
									Alpha: 1.0,
								},
							},
						},
					},
					Index: 3,
				},
			},
			// Add conditional formatting for "4 – Confident redo" (Light Green)
			{
				AddConditionalFormatRule: &sheets.AddConditionalFormatRuleRequest{
					Rule: &sheets.ConditionalFormatRule{
						Ranges: []*sheets.GridRange{
							{
								SheetId:          0,
								StartColumnIndex: 2,
								EndColumnIndex:   3,
							},
						},
						BooleanRule: &sheets.BooleanRule{
							Condition: &sheets.BooleanCondition{
								Type: "TEXT_EQ",
								Values: []*sheets.ConditionValue{
									{UserEnteredValue: "4 – Confident redo"},
								},
							},
							Format: &sheets.CellFormat{
								BackgroundColor: &sheets.Color{
									Red:   0.6,
									Green: 1.0,
									Blue:  0.0,
									Alpha: 1.0,
								},
							},
						},
					},
					Index: 4,
				},
			},
			// Add conditional formatting for "5 – Perfectly repeatable" (Green)
			{
				AddConditionalFormatRule: &sheets.AddConditionalFormatRuleRequest{
					Rule: &sheets.ConditionalFormatRule{
						Ranges: []*sheets.GridRange{
							{
								SheetId:          0,
								StartColumnIndex: 2,
								EndColumnIndex:   3,
							},
						},
						BooleanRule: &sheets.BooleanRule{
							Condition: &sheets.BooleanCondition{
								Type: "TEXT_EQ",
								Values: []*sheets.ConditionValue{
									{UserEnteredValue: "5 – Perfectly repeatable"},
								},
							},
							Format: &sheets.CellFormat{
								BackgroundColor: &sheets.Color{
									Red:   0.0,
									Green: 0.8,
									Blue:  0.0,
									Alpha: 1.0,
								},
							},
						},
					},
					Index: 5,
				},
			},
		},
	}).Do()
	if err != nil {
		log.Fatalf("Unable to set confidence validation and formatting: %v", err)
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
