.PHONY: build-app run clean create-tables

#defining db variables
MIGRATION_SCRIPT := ./sqlite/migration.sql
DATABASE := daren.db

build-app:
	go build -o bin/daren ./cmd/daren.go

run: build-app
	@./bin/daren

clean:
	@rm -rf bin

create-tables: $(MIGRATION_SCRIPT) $(DATABASE)
	sqlite3 $(DATABASE) < $(MIGRATION_SCRIPT)

test:
	go test ./...

test-verbose:
	go test ./... -v
