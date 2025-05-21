package server

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"

	"encoding/json"

	daren "github.com/wintergathering/daren2"
)

type Server struct {
	Router *http.ServeMux
	Srvr   *http.Server

	//services
	DareService    daren.DareService
	PaybackService daren.PaybackService
	//templates
	Templates *template.Template
	StaticFS  embed.FS

	//other stuff?
	//config? logger?
}
type spaFileSystem struct {
	staticFS  fs.FS
	indexPath string // Typically "index.html"
}

func (sfs spaFileSystem) Open(name string) (fs.File, error) {
	// Try to open the requested file directly.
	file, err := sfs.staticFS.Open(name)
	if err == nil {
		// If the file exists and there's no error, return it.
		return file, nil
	}

	// If the file does not exist (os.IsNotExist(err)),
	// or if it's a directory (which http.FileServer would normally handle by looking for index.html),
	// we want to serve the SPA's index.html for client-side routing.
	// A simple check for os.IsNotExist is often sufficient for basic SPAs.
	if errors.Is(err, os.ErrNotExist) {
		// If the file wasn't found, try to open the index.html file instead.
		return sfs.staticFS.Open(sfs.indexPath)
	}

	// For any other error (e.g., permission denied), return the error.
	return nil, err
}

// server constructor
func NewServer(addr string, ds daren.DareService, ps daren.PaybackService, templatePaths string, staticFS embed.FS) *Server {
	tmpl, err := template.ParseGlob(templatePaths)

	if err != nil {
		log.Fatalf("couldn't parse templates: %v", err.Error())
	}

	return &Server{
		Router: http.NewServeMux(),
		Srvr: &http.Server{
			Addr: addr,
		},
		DareService:    ds,
		PaybackService: ps,
		Templates:      tmpl,
		StaticFS:       staticFS,
	}
}

// CORS wrapper function for front-end dev, allows request from localhost:5173
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers
		// Allow requests from your Svelte dev server.
		// For production, you'd want to be more specific or configure this.
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // Common SvelteKit dev port
		// Or for more permissive development: w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // If you plan to use cookies/sessions

		// Handle preflight requests (OPTIONS method)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Next
		next.ServeHTTP(w, r)
	})
}

func (s *Server) registerRoutes() {

	//TODO -- use custom css here

	//html
	//s.Router.HandleFunc("GET /", s.handleIndex) <- this is banned for some reason
	s.Router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			if r.Method == http.MethodGet {
				s.handleIndex(w, r) // Call your original index handler
			} else {
				http.Error(w, "Method not allowed for /", http.StatusMethodNotAllowed)
			}
			return // Important: return after handling /
		}
	}))
	s.Router.HandleFunc("GET /success", s.handleSuccess)
	s.Router.HandleFunc("GET /create", s.handleGetCreateDare)
	s.Router.HandleFunc("POST /create", s.handleCreateDare)
	s.Router.HandleFunc("GET /all", s.handleGetAllDares)
	s.Router.HandleFunc("GET /random", s.HandleGetRandomDare)
	// --- STEE ZONE ---
	// --- END STEE ZONE

	//api
	s.Router.HandleFunc("POST /api/v1/dare/create", s.handleApiCreateDare)
	s.Router.HandleFunc("GET /api/v1/dare/id/{id}", s.handleApiGetDareById)
	s.Router.HandleFunc("GET /api/v1/dare/random", s.handleApiGetRandomDare)
	s.Router.HandleFunc("GET /api/v1/dare/all", s.handleApiGetAllDares)
	s.Router.HandleFunc("DELETE /api/v1/dare/id/{id}", s.handleApiDeleteDare)

	// --- STEE ZONE --- Payback API Routes ---
	s.Router.HandleFunc("POST /api/v1/payback/participants", s.handlePaybackCreateParticipant)
	s.Router.HandleFunc("GET /api/v1/payback/participants", s.handlePaybackGetAllParticipants)
	// --- Trip API Routes ---
	s.Router.HandleFunc("POST /api/v1/payback/trips", s.handlePaybackCreateTrip)
	s.Router.HandleFunc("GET /api/v1/payback/trips", s.handlePaybackGetAllTrips)
	s.Router.HandleFunc("POST /api/v1/payback/trips/{tripID}/participants", s.handlePaybackAddParticipantToTrip)
	s.Router.HandleFunc("GET /api/v1/payback/trips/{tripID}/participants", s.handlePaybackGetParticipantsForTrip)
	s.Router.HandleFunc("DELETE /api/v1/payback/trips/{tripID}/participants/{participantID}", s.handlePaybackRemoveParticipantFromTrip)
	// --- Purchase and Debt API Routes ---
	s.Router.HandleFunc("POST /api/v1/payback/purchases", s.handlePaybackCreateOriginalPurchase)
	s.Router.HandleFunc("GET /api/v1/payback/trips/{tripID}/purchases", s.handlePaybackGetPurchasesForTrip)
	s.Router.HandleFunc("GET /api/v1/payback/purchases/{purchaseID}/debts", s.handlePaybackGetDebtsForPurchase)
	s.Router.HandleFunc("GET /api/v1/payback/trips/{tripID}/balances", s.handlePaybackGetTripBalances) // New route
	// --- Frontend Static File Serving for /payback ---
	staticContentFS, err := fs.Sub(s.StaticFS, "web/static")
	if err != nil {
		// This is a fatal error during startup if the web/static directory isn't found or embedded correctly.
		log.Fatalf("Failed to create sub FS for static files: %v", err)
	}

	paybackHandler := http.StripPrefix("/payback/", http.FileServer(http.FS(spaFileSystem{staticFS: staticContentFS, indexPath: "index.html"})))
	s.Router.Handle("/payback/", paybackHandler)
}

func (s *Server) Run() {
	//register routes
	s.registerRoutes()

	s.Srvr.Handler = enableCORS(s.Router)

	fmt.Println("Running Daren on ", s.Srvr.Addr)

	s.Srvr.ListenAndServe()
}

// utility ------------
func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}
