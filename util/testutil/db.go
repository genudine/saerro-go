package testutil

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

// GetTestDB standardizes what type of DB tests use.
func GetTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", t.TempDir()+"/test.db")
	if err != nil {
		t.Fatalf("test shim: sqlite open failed, %v", err)
	}

	return db
}
