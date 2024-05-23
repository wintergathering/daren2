FROM golang:latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

ENV PORT 8080

RUN go build -o bin/daren ./cmd/daren.go

EXPOSE $PORT

CMD ["./bin/daren"]