package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	daren "github.com/wintergathering/daren2"
	"github.com/wintergathering/daren2/server"
	"github.com/wintergathering/daren2/sqlite"

	_ "modernc.org/sqlite"
)

const addr = ":8080"
const daresDsn = "./daren.db"     // DSN for the original dares database
const paybackDsn = "./payback.db" // DSN for your new payback database
const templatePaths = "templates/*.html"

// Function to execute a single migration file (this remains the same)
func executeMigration(db *sql.DB, migrationFilePath string) error {
	migrationSQL, err := os.ReadFile(migrationFilePath)
	if err != nil {
		if os.IsNotExist(err) { // Be more specific if file not found is okay
			log.Printf("Warning: Migration file %s not found. Skipping.", migrationFilePath)
			return nil
		}
		log.Printf("Warning: Could not read migration file %s: %v. Skipping.", migrationFilePath, err)
		return err // Return error if reading fails for other reasons
	}

	if len(migrationSQL) == 0 {
		log.Printf("Warning: Migration file %s is empty. Skipping.", migrationFilePath)
		return nil
	}

	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		return err
	}
	log.Printf("Successfully executed migration: %s on database connected to %v", migrationFilePath, db) // Added DB info for clarity
	return nil
}

func main() {
	// --- Setup Dares Database ---
	daresDb, err := sql.Open("sqlite", daresDsn)
	if err != nil {
		log.Fatalf("Couldn't open dares database (%s): %s", daresDsn, err)
	}
	defer daresDb.Close() // Good practice to close it

	// Run Dares Migration
	daresMigrationPath := filepath.Join("sqlite", "migration.sql") // Assumes dares migration is named migration.sql
	err = executeMigration(daresDb, daresMigrationPath)
	if err != nil {
		log.Fatalf("Failed to execute dares migration (%s) on %s: %v", daresMigrationPath, daresDsn, err)
	}
	log.Printf("Dares database setup complete on %s", daresDsn)

	// --- Setup Payback Database ---
	paybackDb, err := sql.Open("sqlite", paybackDsn)
	if err != nil {
		log.Fatalf("Couldn't open payback database (%s): %s", paybackDsn, err)
	}
	defer paybackDb.Close() // Good practice to close it

	// Run Payback Migration
	paybackMigrationPath := filepath.Join("sqlite", "payback_migration.sql")
	err = executeMigration(paybackDb, paybackMigrationPath)
	if err != nil {
		log.Fatalf("Failed to execute payback migration (%s) on %s: %v", paybackMigrationPath, paybackDsn, err)
	}
	log.Printf("Payback database setup complete on %s", paybackDsn)

	// --- Initialize Services ---
	// DareService uses the daresDb
	dareService := sqlite.NewDareService(daresDb)

	var paybackService daren.PaybackService
	paybackService = sqlite.NewPaybackService(paybackDb)

	// --- Setup and Run Server ---
	s := server.NewServer(addr, dareService, paybackService, templatePaths, daren.EmbeddedAssets)

	s.Run()
}
