package server

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/wintergathering/daren2"
)

// this will probably get more than this
type Server struct {
	Router      *mux.Router
	Srvr        *http.Server
	DareService daren2.DareService
	Templates   *template.Template
}

// constructor to build a new server
func NewServer(r *mux.Router, ds daren2.DareService, addr string) *Server {

	tmpl, err := template.ParseGlob("templates/*.html")

	if err != nil {
		log.Fatalf("couldn't parse templates: %v", err.Error())
	}

	return &Server{
		Router: r,
		Srvr: &http.Server{
			Addr: addr,
		},
		DareService: ds,
		Templates:   tmpl,
	}
}

// register routes on the server
// will want 'regular' routes that render HTML and 'api' routes (/api/..) that return JSON
func (s *Server) registerRoutes() {
	//api routes
	s.Router.HandleFunc("/api", handleAPIIndex).Methods("GET")
	s.Router.HandleFunc("/api/dare/create", s.handleAPICreateDare).Methods("POST")
	s.Router.HandleFunc("/api/dare/getRandom", s.handleAPIGetRandDare).Methods("GET")
	s.Router.HandleFunc("/api/dare/getAll", s.handleAPIGetAllDares).Methods("GET")

	//html routes
	s.Router.HandleFunc("/", s.handleIndex).Methods("GET")
	s.Router.HandleFunc("/", s.handleCreateDare).Methods("POST")
	s.Router.HandleFunc("/all_dares", s.handleGetAllDares).Methods("GET")
	s.Router.HandleFunc("/rand_dare", s.handleGetRandDare)
}

// run the server
func (s *Server) Run() {
	s.registerRoutes()

	fmt.Printf("Daren is running on port %s", s.Srvr.Addr)

	s.Srvr.Handler = s.Router

	s.Srvr.ListenAndServe()
}
