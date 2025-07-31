package tests

import (
	"Leetcode-or-Explode-Bot/db"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"testing"
	"time"
)

func TestSameDaySubm(t *testing.T) {

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer mockDB.Close()

	today := time.Now().Format("2006-01-02")
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(today).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(true))

	exists := db.SameDaySubm(mockDB, "two-sum", "user123", today)

	// Check result
	if !exists {
		t.Error("❌ Expected submission to exist, got false")
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("❌ Unmet expectations: %v", err)
	}
	type testParams struct {
		problemName string
		userID      string
	}

}
