package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	daren "github.com/wintergathering/daren2"
)

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	//set header to html
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	s.Templates.ExecuteTemplate(w, "placeholder.html", nil)
}

// html routes ---------------
func (s *Server) handleCreateDare(w http.ResponseWriter, r *http.Request) {
	//todo
}

// api routes ----------------
func (s *Server) handleApiCreateDare(w http.ResponseWriter, r *http.Request) {
	var d *daren.Dare

	err := json.NewDecoder(r.Body).Decode(&d)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	dareID, err := s.DareService.CreateDare(d)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	msg := make(map[string]int)

	msg["Dare created"] = dareID

	writeJSON(w, http.StatusOK, msg)
}

// get dare by id
func (s *Server) handleApiGetDareById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	dare, err := s.DareService.GetDareByID(id)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, dare)
}

// get random dare
func (s *Server) handleApiGetRandomDare(w http.ResponseWriter, r *http.Request) {
	dare, err := s.DareService.GetRandomDare()

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	id := dare.ID

	err = s.DareService.MarkDareSeen(id)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, dare)
}

// get all dares
func (s *Server) handleApiGetAllDares(w http.ResponseWriter, r *http.Request) {
	dares, err := s.DareService.GetAllDares()

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, dares)

}

// delete dare
func (s *Server) handleApiDeleteDare(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	err = s.DareService.DeleteDare(id)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	msg := make(map[string]int)

	msg["Dare deleted"] = id

	writeJSON(w, http.StatusOK, msg)
}
