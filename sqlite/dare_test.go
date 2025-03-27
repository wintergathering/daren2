package sqlite_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	daren "github.com/wintergathering/daren2"
	"github.com/wintergathering/daren2/sqlite"
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

// create a new test dare service
func NewTestDareService(t *testing.T) (daren.DareService, func()) {
	t.Helper() //marks the function as a helper

	db, cleanup := NewTestDB(t)

	ds := sqlite.NewDareService(db)

	return ds, cleanup
}

func TestDareService_CreateDare(t *testing.T) {
	ds, cleanup := NewTestDareService(t)
	defer cleanup()

	//define a test case -- clean this up a bit later
	testCase := struct {
		name    string
		dare    *daren.Dare
		wantID  int
		wantErr bool
	}{
		name: "create a dare",
		dare: &daren.Dare{
			Title:   "Test Dare",
			Text:    "This is a test dare",
			AddedBy: "Test User",
			Seen:    false,
		},
		wantID:  1,
		wantErr: false,
	}

	//RESUME HERE!
}
