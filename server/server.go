package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	daren "github.com/wintergathering/daren2"
)

type Server struct {
	Router *http.ServeMux
	Srvr   *http.Server

	//services
	DareService daren.DareService

	//templates
	Templates *template.Template

	//other stuff?
	//config? logger?
}

// server constructor
func NewServer(addr string, ds daren.DareService, templatePaths string) *Server {
	tmpl, err := template.ParseGlob(templatePaths)

	if err != nil {
		log.Fatalf("couldn't parse templates: %v", err.Error())
	}

	return &Server{
		Router: http.NewServeMux(),
		Srvr: &http.Server{
			Addr: addr,
		},
		// DareService: ds,
		Templates: tmpl,
	}
}

func (s *Server) registerRoutes() {
	s.Router.HandleFunc("GET /", handleIndex)
	s.Router.HandleFunc("GET /weather", handleWeather)
}

func (s *Server) Run() {
	//register routes
	s.registerRoutes()

	s.Srvr.Handler = s.Router

	fmt.Println("Running Daren")

	s.Srvr.ListenAndServe()
}
