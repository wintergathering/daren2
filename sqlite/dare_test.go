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

type testDare struct {
	name    string
	dare    *daren.Dare
	wantID  int
	wantErr bool
}

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

	err = populateDB(t, db)

	if err != nil {
		t.Fatalf("Failed to populate database: %v", err)
	}

	cleanup := func() {
		err := db.Close()
		if err != nil {
			t.Fatalf("Failed to close database: %v", err)
		}
	}

	return db, cleanup
}

func populateDB(t *testing.T, db *sql.DB) error {
	t.Helper() //marks the function as a helper

	_, err := db.Exec(`INSERT INTO dares (title, dare_text, added_by) VALUES (?, ?, ?), (?, ?, ?)`,
		"Test Dare", "This is a test dare", "Test User",
		"Test Dare 2", "This is another test dare", "EE")

	if err != nil {
		return err
	}

	return nil
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

	testCase := testDare{
		name: "create a dare",
		dare: &daren.Dare{
			Title:   "Test Dare 3",
			Text:    "This is a third test dare",
			AddedBy: "Stee",
		},
		wantID:  3,
		wantErr: false,
	}

	//run the test
	t.Run(testCase.name, func(t *testing.T) {
		gotID, err := ds.CreateDare(testCase.dare)

		if (err != nil) != testCase.wantErr {
			t.Errorf("CreateDare() error = %v, wantErr %v", err, testCase.wantErr)
			return
		}

		if gotID != testCase.wantID {
			t.Errorf("CreateDare() gotID = %v, want %v", gotID, testCase.wantID)
		}

	})
}

func TestDareService_GetDareByID(t *testing.T) {
	ds, cleanup := NewTestDareService(t)
	defer cleanup()

	testCase := testDare{
		name: "get a dare by id",
		dare: &daren.Dare{
			Title:   "Test Dare 2",
			Text:    "This is another test dare",
			AddedBy: "EE",
		},
		wantID:  2,
		wantErr: false,
	}

	t.Run(testCase.name, func(t *testing.T) {
		got, err := ds.GetDareByID(testCase.wantID)

		if (err != nil) != testCase.wantErr {
			t.Errorf("GetDareByID() error = %v, wantErr %v", err, testCase.wantErr)
			return
		}

		if got.Title != testCase.dare.Title {
			t.Errorf("GetDareByID() title got = %v, want %v", got.Title, testCase.dare.Title)
			return
		}

		if got.Text != testCase.dare.Text {
			t.Errorf("GetDareByID() text got = %v, want %v", got.Text, testCase.dare.Text)
			return
		}

		if got.AddedBy != testCase.dare.AddedBy {
			t.Errorf("GetDareByID() addedBy got = %v, want %v", got.Text, testCase.dare.Text)
			return
		}

	})

	//test case for a row that doesn't exist -- should error
	testCase = testDare{
		name:    "get a dare by id that doesn't exist",
		dare:    &daren.Dare{},
		wantID:  100,
		wantErr: true,
	}

	t.Run(testCase.name, func(t *testing.T) {
		_, err := ds.GetDareByID(testCase.wantID)
		if (err != nil) != testCase.wantErr {
			t.Errorf("GetDareByID() error = %v, wantErr %v", err, testCase.wantErr)
			return
		}

	})

}

// test GetRandomDare()
func TestDareService_GetRandomDare(t *testing.T) {
	ds, cleanup := NewTestDareService(t)
	defer cleanup()

	testCase := testDare{
		name:    "get a random dare",
		dare:    nil,
		wantErr: false,
	}

	t.Run(testCase.name, func(t *testing.T) {
		got, err := ds.GetRandomDare()

		if (err != nil) != testCase.wantErr {
			t.Errorf("GetRandomDare() err = %v, wantErr %v", err, testCase.wantErr)
			return
		}

		if got == nil {
			t.Errorf("GetRandomDare() returned nil, expected a dare")
			return
		}

		//check that the returned dare has the fields we need
		if got.ID == 0 {
			t.Errorf("GetRandomDare() returned dare with ID 0, expected a nonzero value")
			return
		}

		if got.Title == "" {
			t.Errorf("GetRandomDare() returned dare with empty title, expected a non-empty value")
			return
		}

		if got.Text == "" {
			t.Errorf("GetRandomDare() returned dare with empty text, expected a non-empty value")
			return
		}

		if got.AddedBy == "" {
			t.Errorf("GetRandomDare() returned dare with empty addedBy, expected a non-empty value")
			return
		}

		//i think this is getting coerced to a bool? it's not throwing an error
		if got.Seen {
			t.Errorf("GetRandomDare() returned dare with seen set to true, expected false")
			return
		}

	})
}

// test GetAllDares()
func TestDareService_GetAllDares(t *testing.T) {
	ds, cleanup := NewTestDareService(t)
	defer cleanup()

	testCase := testDare{
		name:    "get all dares",
		dare:    nil,
		wantErr: false,
	}

	t.Run(testCase.name, func(t *testing.T) {
		got, err := ds.GetAllDares()

		if (err != nil) != testCase.wantErr {
			t.Errorf("GetAllDares() err = %v, wantErr = %v", err, testCase.wantErr)
			return
		}

		//hardcoding 2 for now since I know there are 2 dares in the test db, but might want to address this later?
		if len(got) != 2 {
			t.Errorf("GetAllDares() returned %d dares, expected 2", len(got))
			return
		}

	})
}

func TestDareService_MarkDareSeen(t *testing.T) {
	ds, cleanup := NewTestDareService(t)
	defer cleanup()

	testCase := testDare{
		name: "mark dare as seen",
		dare: &daren.Dare{
			ID:   2,
			Seen: true,
		},
		wantID:  2,
		wantErr: false,
	}

	t.Run(testCase.name, func(t *testing.T) {
		err := ds.MarkDareSeen(testCase.wantID)

		if (err != nil) != testCase.wantErr {
			t.Errorf("MarkDareSeen() err = %v, wantErr = %v", err, testCase.wantErr)
			return
		}

		//ignoring the error here since i've tested this elsewhere
		got, _ := ds.GetDareByID(testCase.wantID)

		if got.Seen != testCase.dare.Seen {
			t.Errorf("MarkDareSeen() returned %v, expected %v", got.Seen, testCase.dare.Seen)
			return
		}
	})
}

func TestDareService_DeleteDare(t *testing.T) {
	ds, cleanup := NewTestDareService(t)
	defer cleanup()

	testCase := testDare{
		name:    "delete dare",
		dare:    nil,
		wantID:  1,
		wantErr: false,
	}

	t.Run(testCase.name, func(t *testing.T) {
		err := ds.DeleteDare(testCase.wantID)

		if (err != nil) != testCase.wantErr {
			t.Errorf("DeleteDare() err = %v, wantErr = %v", err, testCase.wantErr)
		}

		//try to retrieve the deleted dare
		_, err = ds.GetDareByID(testCase.wantID)

		if err == nil {
			t.Errorf("DeleteDare() did not delete dare with ID %d", testCase.wantID)
		}

	})
}
