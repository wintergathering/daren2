.PHONY: build-app run clean create-tables populate-tables depopulate-tables test test-verbose

#defining db variables
MIGRATION_SCRIPT := ./sqlite/migration.sql
POPULATE_DB_SCRIPT := ./sqlite/populate.sql
DEPOPULATE_DB_SCRIPT := ./sqlite/depopulate.sql
DATABASE := daren.db

build-app:
	go build -o bin/daren ./cmd/daren.go

run: build-app
	@./bin/daren

clean:
	@rm -rf bin

create-tables: $(MIGRATION_SCRIPT) $(DATABASE)
	sqlite3 $(DATABASE) < $(MIGRATION_SCRIPT)

populate-tables: $(POPULATE_DB_SCRIPT) $(DATABASE)
	sqlite3 $(DATABASE) < $(POPULATE_DB_SCRIPT)

depopulate-tables: $(DEPOPULATE_DB_SCRIPT) $(DATABASE)
	sqlite3 $(DATABASE) < $(DEPOPULATE_DB_SCRIPT)

test:
	go test ./...

test-verbose:
	go test ./... -v
