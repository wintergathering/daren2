package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/wintergathering/daren2"
)

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

// utility func to write out JSON
func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(v)
}
