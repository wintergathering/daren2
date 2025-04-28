package main

import (
	"database/sql"
	"log"

	"github.com/wintergathering/daren2/server"
	"github.com/wintergathering/daren2/sqlite"

	_ "modernc.org/sqlite"
)

const addr = ":8080"
const dsn = "./daren.db"
const templatePaths = "templates/*.html"
const logFilePath = "daren.log"

func main() {
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		log.Fatalf("Couldn't open database: %s", err)
	}

	ds := sqlite.NewDareService(db)

	s := server.NewServer(addr, ds, templatePaths, logFilePath)

	s.Run()
}
