package server

import (
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"encoding/json"

	daren "github.com/wintergathering/daren2"
)

type Server struct {
	Router *http.ServeMux
	Srvr   *http.Server

	//services
	DareService daren.DareService

	//templates
	Templates *template.Template

	//logging
	Logger *slog.Logger

	//other stuff?
	//config? logger?
}

// server constructor
func NewServer(addr string, ds daren.DareService, templatePaths string, logFilePath string) *Server {
	tmpl, err := template.ParseGlob(templatePaths)

	if err != nil {
		log.Fatalf("couldn't parse templates: %v", err.Error())
	}

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("couldn't open log file: %v", err.Error())
	}

	lh := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})

	logger := slog.New(lh)

	return &Server{
		Router: http.NewServeMux(),
		Srvr: &http.Server{
			Addr: addr,
		},
		DareService: ds,
		Templates:   tmpl,
		Logger:      logger,
	}
}

func (s *Server) registerRoutes() {

	//TODO -- use custom css here

	//html
	s.Router.HandleFunc("GET /", s.handleIndex)
	s.Router.HandleFunc("GET /success", s.handleSuccess)
	s.Router.HandleFunc("GET /create", s.handleGetCreateDare)
	s.Router.HandleFunc("POST /create", s.handleCreateDare)
	s.Router.HandleFunc("GET /all", s.handleGetAllDares)
	s.Router.HandleFunc("GET /random", s.HandleGetRandomDare)

	//api
	s.Router.HandleFunc("POST /api/v1/dare/create", s.handleApiCreateDare)
	s.Router.HandleFunc("GET /api/v1/dare/id/{id}", s.handleApiGetDareById)
	s.Router.HandleFunc("GET /api/v1/dare/random", s.handleApiGetRandomDare)
	s.Router.HandleFunc("GET /api/v1/dare/all", s.handleApiGetAllDares)
	s.Router.HandleFunc("DELETE /api/v1/dare/id/{id}", s.handleApiDeleteDare)
}

func (s *Server) Run() {
	//register routes
	s.registerRoutes()

	s.Srvr.Handler = s.Router

	s.Logger.Info("Starting server", slog.String("addr", s.Srvr.Addr))

	fmt.Printf("Starting server on %s\n", s.Srvr.Addr)

	s.Srvr.ListenAndServe()
}

// log utility
func (s *Server) logRequest(r *http.Request, statusCode int, err error) {
	var e string
	if err == nil {
		e = ""
	} else {
		e = err.Error()
	}
	g := slog.Group(
		"request_details",
		slog.Time("timestamp", time.Now()),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("remote_addr", r.RemoteAddr),
		slog.Int("status_code", statusCode),
		slog.Any("error", e),
	)

	s.Logger.Info("request received",
		g,
	)
}

// utility ------------
func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}
