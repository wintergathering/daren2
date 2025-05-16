package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

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

	//other stuff?
	//config? logger?
}

// server constructor
func NewServer(addr string, ds daren.DareService, ps daren.PaybackService, templatePaths string) *Server {
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
	s.Router.HandleFunc("GET /", s.handleIndex)
	s.Router.HandleFunc("GET /success", s.handleSuccess)
	s.Router.HandleFunc("GET /create", s.handleGetCreateDare)
	s.Router.HandleFunc("POST /create", s.handleCreateDare)
	s.Router.HandleFunc("GET /all", s.handleGetAllDares)
	s.Router.HandleFunc("GET /random", s.HandleGetRandomDare)
	// --- STEE ZONE ---
	s.Router.HandleFunc("GET /payback", s.handPaybackHome)
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
