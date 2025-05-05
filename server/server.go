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
	s.Router.Handle("GET /", s.LoggingMiddleware(http.HandlerFunc(s.handleIndex)))
	s.Router.Handle("GET /success", s.LoggingMiddleware(http.HandlerFunc(s.handleSuccess)))
	s.Router.Handle("GET /create", s.LoggingMiddleware(http.HandlerFunc(s.handleGetCreateDare)))
	s.Router.Handle("POST /create", s.LoggingMiddleware(http.HandlerFunc(s.handleCreateDare)))
	s.Router.Handle("GET /all", s.LoggingMiddleware(http.HandlerFunc(s.handleGetAllDares)))
	s.Router.Handle("GET /random", s.LoggingMiddleware(http.HandlerFunc(s.handleGetRandomDare)))

	//api with logging
	s.Router.Handle("POST /api/v1/dare/create", s.LoggingMiddleware(http.HandlerFunc(s.handleApiCreateDare)))
	s.Router.Handle("GET /api/v1/dare/id/{id}", s.LoggingMiddleware(http.HandlerFunc(s.handleApiGetDareById)))
	s.Router.Handle("GET /api/v1/dare/random", s.LoggingMiddleware(http.HandlerFunc(s.handleApiGetRandomDare)))
	s.Router.Handle("GET /api/v1/dare/all", s.LoggingMiddleware(http.HandlerFunc(s.handleApiGetAllDares)))
	s.Router.Handle("DELETE /api/v1/dare/id/{id}", s.LoggingMiddleware(http.HandlerFunc(s.handleApiDeleteDare)))
}

func (s *Server) Run() {
	//register routes
	s.registerRoutes()

	s.Srvr.Handler = s.Router

	s.Logger.Info("Starting server", slog.String("addr", s.Srvr.Addr))

	fmt.Printf("Starting server on %s\n", s.Srvr.Addr)

	s.Srvr.ListenAndServe()
}

// logging middleware
func (s *Server) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		s.Logger.Info("request received",
			slog.Time("timestamp", time.Now()),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.Duration("duration", duration),
		)
	})
}

// utility ------------
func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}
