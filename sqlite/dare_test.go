package sqlite_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func NewTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper() //marks the function as a helper

	db, err := sql.Open("sqlite", ":memory:")

	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	migrationFile, err := os.ReadFile(filepath.Join("..", "sqlite", "migration.sql"))

	if err != nil {
		t.Fatalf("Failed to read migration file: %v", err)
	}

	_, err = db.Exec(string(migrationFile))

	if err != nil {
		t.Fatalf("Failed to execute migration file: %v", err)
	}

	cleanup := func() {
		err := db.Close()
		if err != nil {
			t.Fatalf("Failed to close database: %v", err)
		}
	}

	return db, cleanup
}
