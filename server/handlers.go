package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/wintergathering/daren2"
)

//API handlers --------

func handleAPIIndex(w http.ResponseWriter, r *http.Request) {
	msg := make(map[string]string)

	msg["hello"] = "Welcome to Daren's API"

	writeJSON(w, http.StatusOK, msg)
}

// api handler for dare creation
func (s *Server) handleAPICreateDare(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	var d *daren2.Dare

	err := json.NewDecoder(r.Body).Decode(&d)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	id := uuid.New().String()

	//populate other components of dare
	d.UUID = id
	d.Seen = false

	//validate struct
	validate := validator.New()
	err = validate.Struct(d)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	err = s.DareService.CreateDare(ctx, d)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, d)
}

func (s *Server) handleAPIGetAllDares(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	dares, err := s.DareService.GetAllDares(ctx)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, dares)
}

func (s *Server) handleAPIGetRandDare(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	dare, err := s.DareService.GetRandomDare(ctx)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, dare)
}

// HTML handlers ----------------------

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	s.Templates.ExecuteTemplate(w, "index.html", nil)
}

func (s *Server) handleCreateDare(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	err := r.ParseForm()

	if err != nil {
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	title := r.FormValue("title")
	text := r.FormValue("text")
	addedBy := r.FormValue("addedBy")

	d := daren2.NewDare(title, text, addedBy)

	//validate struct
	validate := validator.New()
	err = validate.Struct(d)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	err = s.DareService.CreateDare(ctx, d)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	s.Templates.ExecuteTemplate(w, "return.html", nil)
}

func (s *Server) handleGetRandDare(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	ctx := context.Background()

	d, err := s.DareService.GetRandomDare(ctx)

	if err != nil {
		w.WriteHeader(http.StatusOK)
		s.Templates.ExecuteTemplate(w, "no_dares.html", nil)
		return
	}

	s.Templates.ExecuteTemplate(w, "single_dare.html", d)
}

func (s *Server) handleGetAllDares(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	dares, err := s.DareService.GetAllDares(ctx)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	s.Templates.ExecuteTemplate(w, "all_dares.html", dares)
}

// utils ----------
// utility func to write out JSON
func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(v)
}
