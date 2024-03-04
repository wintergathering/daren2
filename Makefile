build-app:
		go build -o bin/daren ./cmd/daren.go

run: build-app
		@./bin/daren

clean:
		@rm -rf bin