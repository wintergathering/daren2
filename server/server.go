package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wintergathering/daren2"
)

// this will probably get more than this
type Server struct {
	Router      *mux.Router
	Srvr        *http.Server
	DareService daren2.DareService
}

// constructor to build a new server
func NewServer(r *mux.Router, ds daren2.DareService, addr string) *Server {
	return &Server{
		Router: r,
		Srvr: &http.Server{
			Addr: addr,
		},
		DareService: ds,
	}
}

// helper to handle a generic index right now
func handleAPIIndex(w http.ResponseWriter, r *http.Request) {
	msg := make(map[string]string)

	msg["hello"] = "Welcome to Daren's API"

	writeJSON(w, http.StatusOK, msg)
}

// register routes on the server
// will want 'regular' routes that render HTML and 'api' routes (/api/..) that return JSON
func (s *Server) registerRoutes() {
	//api routes
	s.Router.HandleFunc("/api", handleAPIIndex).Methods("GET")
	s.Router.HandleFunc("/api/dare/create", s.handleAPICreateDare).Methods("POST")
	s.Router.HandleFunc("/api/dare/getRandom", s.handleAPIGetRandDare).Methods("GET")
	s.Router.HandleFunc("/api/dare/getAll", s.handleAPIGetAllDares).Methods("GET")
}

// run the server
func (s *Server) Run() {
	s.registerRoutes()

	fmt.Printf("Daren is running on port %s", s.Srvr.Addr)

	s.Srvr.Handler = s.Router

	s.Srvr.ListenAndServe()
}
