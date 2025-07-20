// google_sheets.go
package bot

import (
	"Leetcode-or-Explode-Bot/db"
	"context"
	"fmt"
	"log"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Define spreadsheetID and writeRange globally or pass as needed.
var spreadsheetID = "1Gc3PhSnLSrlcVSEDtQFiQ-rS-pHjWrgQQ64GFR-dwFQ" // TODO: replace with your spreadsheet ID

var writeRange = "Sheet1!A2:I" // TODO: be dynamic based o the userID/dicord name

var scoreMap = map[int8]string{
	0: "0 – No clue",
	1: "1 – Struggle to repeat",
	2: "2 – Might redo poorly",
	3: "3 – Could redo maybe",
	4: "4 – Confident redo",
	5: "5 – Perfectly repeatable",
}

func getGoogleSheets() (*sheets.Service, error) {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile("go/credentials.json"))
	if err != nil {
		return nil, err
	}
	return srv, nil
}

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

	// Create the row to append
	row := []interface{}{
		fmt.Sprintf(`=HYPERLINK("https://leetcode.com/problems/%s/", "%s")`,
			strings.SplitN(subm.ProblemName, "-", 2)[1], // slug
			subm.ProblemName),

		subm.Difficulty,
		scoreMap[int8(subm.ConfidenceScore)],

		subm.SubmittedAt[:10],
		subm.SolveTime,
		strings.Join(subm.Topics, ", "),
		subm.Notes,
	}

	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{row},
	}

	// Append to spreadsheet
	var sheetsName string

	sheetsName, err = db.GetUsernameByUserID(db.DB, subm.UserID)
	if err != nil {
		log.Fatalf("Unable to get username for submission: %v", err)
	}
	writeRange := fmt.Sprintf(sheetsName + "!A2:I")

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
	labels := []string{
		"0 – No clue",
		"1 – Struggle to repeat",
		"2 – Might redo poorly",
		"3 – Could redo maybe",
		"4 – Confident redo",
		"5 – Perfectly repeatable",
	}

	colors := []*sheets.Color{
		{Red: 0.8, Green: 0.0, Blue: 0.0, Alpha: 1.0}, // red
		{Red: 0.9, Green: 0.3, Blue: 0.0, Alpha: 1.0}, // red-orange
		{Red: 1.0, Green: 0.6, Blue: 0.0, Alpha: 1.0}, // orange
		{Red: 1.0, Green: 0.8, Blue: 0.0, Alpha: 1.0}, // yellow
		{Red: 0.6, Green: 1.0, Blue: 0.0, Alpha: 1.0}, // light green
		{Red: 0.0, Green: 0.8, Blue: 0.0, Alpha: 1.0}, // green
	}

	var requests []*sheets.Request

	// Add Data Validation
	validationValues := make([]*sheets.ConditionValue, len(labels))
	for i, label := range labels {
		validationValues[i] = &sheets.ConditionValue{UserEnteredValue: label}
	}
	requests = append(requests, &sheets.Request{
		SetDataValidation: &sheets.SetDataValidationRequest{
			Range: &sheets.GridRange{
				SheetId:          0,
				StartColumnIndex: 2,
				EndColumnIndex:   3,
			},
			Rule: &sheets.DataValidationRule{
				Condition: &sheets.BooleanCondition{
					Type:   "ONE_OF_LIST",
					Values: validationValues,
				},
				Strict:       true,
				ShowCustomUi: true,
			},
		},
	})

	// Add conditional formatting for each label
	for i, label := range labels {
		requests = append(requests, &sheets.Request{
			AddConditionalFormatRule: &sheets.AddConditionalFormatRuleRequest{
				Index: int64(i),
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
								{UserEnteredValue: label},
							},
						},
						Format: &sheets.CellFormat{
							BackgroundColor: colors[i],
						},
					},
				},
			},
		})
	}

	resp, err := srv.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		log.Fatalf("Unable to set data validation: %v", err)
	}
	print(resp)
	fmt.Println("✅ Confidence validation and formatting set successfully.")
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

// createNewSheetWithTitle creates a new sheet with the given title in the specified spreadsheet.
func createNewSheetWithTitle(srv *sheets.Service, spreadsheetID, title string) error {
	// Step 1: Create the new sheet
	_, err := srv.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				AddSheet: &sheets.AddSheetRequest{
					Properties: &sheets.SheetProperties{
						Title: title,
					},
				},
			},
		},
	}).Do()
	if err != nil {
		return fmt.Errorf("failed to create sheet '%s': %v", title, err)
	}
	fmt.Printf("✅ Sheet '%s' created successfully.\n", title)

	// Step 2: Write headers
	headers := &sheets.ValueRange{
		Range: fmt.Sprintf("%s!A1:I1", title),
		Values: [][]interface{}{{
			"Problem", "Difficulty", "Confidence", "Date", "Time (min)", "Topics", "Notes", "Counter", "UserID",
		}},
	}
	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, headers.Range, headers).
		ValueInputOption("RAW").
		Context(context.Background()).
		Do()
	if err != nil {
		return fmt.Errorf("failed to write headers: %v", err)
	}

	// Step 3: Freeze header row and enable filters
	sheetMeta, err := srv.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		return fmt.Errorf("failed to get spreadsheet metadata: %v", err)
	}

	var sheetID int64 = -1
	for _, s := range sheetMeta.Sheets {
		if s.Properties.Title == title {
			sheetID = s.Properties.SheetId
			break
		}
	}
	if sheetID == -1 {
		return fmt.Errorf("failed to find sheet ID for '%s'", title)
	}

	_, err = srv.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				UpdateSheetProperties: &sheets.UpdateSheetPropertiesRequest{
					Properties: &sheets.SheetProperties{
						SheetId:        sheetID,
						GridProperties: &sheets.GridProperties{FrozenRowCount: 1},
					},
					Fields: "gridProperties.frozenRowCount",
				},
			},
		},
	}).Do()
	if err != nil {
		return fmt.Errorf("failed to apply filters or freeze header: %v", err)
	}

	return nil
}

// deleteSheetByTitle deletes a sheet from the spreadsheet by its title.
func deleteSheetByTitle(srv *sheets.Service, spreadsheetID, title string) error {
	// First, get spreadsheet metadata to find the sheet ID
	resp, err := srv.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve spreadsheet metadata: %v", err)
	}

	var sheetID int64 = -1
	for _, sheet := range resp.Sheets {
		if sheet.Properties.Title == title {
			sheetID = sheet.Properties.SheetId
			break
		}
	}

	if sheetID == -1 {
		return fmt.Errorf("sheet with title '%s' not found", title)
	}

	// Build the delete sheet request
	requests := []*sheets.Request{
		{
			DeleteSheet: &sheets.DeleteSheetRequest{
				SheetId: sheetID,
			},
		},
	}

	// Execute the request
	_, err = srv.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		return fmt.Errorf("failed to delete sheet '%s': %v", title, err)
	}
	fmt.Printf("🗑️ Sheet '%s' deleted successfully.\n", title)
	return nil
}
