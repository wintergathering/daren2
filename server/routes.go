package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	daren "github.com/wintergathering/daren2"
)

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	//set header to html
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	headerStatus := http.StatusOK
	w.WriteHeader(headerStatus)

	s.logRequest(r, headerStatus, nil)

	s.Templates.ExecuteTemplate(w, "index.html", nil)
}

func (s *Server) handleSuccess(w http.ResponseWriter, r *http.Request) {
	//set header to html
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	headerStatus := http.StatusOK
	w.WriteHeader(headerStatus)

	s.logRequest(r, headerStatus, nil)

	s.Templates.ExecuteTemplate(w, "dare_success.html", nil)
}

// html routes ---------------
func (s *Server) handleCreateDare(w http.ResponseWriter, r *http.Request) {
	d := &daren.Dare{}

	//check request method
	if r.Method != http.MethodPost {
		s.logRequest(r, http.StatusMethodNotAllowed, nil)
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//access form values
	err := r.ParseForm()

	if err != nil {
		s.logRequest(r, http.StatusBadRequest, err)
		w.WriteHeader(http.StatusBadRequest)
		s.Templates.ExecuteTemplate(w, "error.html", nil)
		return
	}

	d.Title = r.FormValue("title")
	d.Text = r.FormValue("text")
	d.AddedBy = r.FormValue("addedBy")

	_, err = s.DareService.CreateDare(d)

	//keeping this for now, but will want to return a real page in the future
	if err != nil {
		s.logRequest(r, http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//redirect to /success
	s.logRequest(r, http.StatusSeeOther, nil)
	w.WriteHeader(http.StatusSeeOther)
	http.Redirect(w, r, "/success", http.StatusSeeOther)

}

func (s *Server) handleGetCreateDare(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	s.logRequest(r, http.StatusOK, nil)

	s.Templates.ExecuteTemplate(w, "create_dare.html", nil)
}

func (s *Server) HandleGetRandomDare(w http.ResponseWriter, r *http.Request) {
	dare, err := s.DareService.GetRandomDare()

	if errors.Is(err, daren.ErrNoDare) {
		s.logRequest(r, http.StatusNotFound, err)
		w.WriteHeader(http.StatusNotFound)
		s.Templates.ExecuteTemplate(w, "no_dares.html", nil)
		return
	}

	if err != nil {
		s.logRequest(r, http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		s.Templates.ExecuteTemplate(w, "error.html", nil)
		return
	}

	err = s.DareService.MarkDareSeen(dare.ID)

	if err != nil {
		s.logRequest(r, http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		s.Templates.ExecuteTemplate(w, "error.html", nil)
		return
	}

	s.logRequest(r, http.StatusOK, nil)
	w.WriteHeader(http.StatusOK)
	s.Templates.ExecuteTemplate(w, "single_dare.html", dare)

}

func (s *Server) handleGetAllDares(w http.ResponseWriter, r *http.Request) {
	dares, err := s.DareService.GetAllDares()

	if err != nil {
		s.logRequest(r, http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		s.Templates.ExecuteTemplate(w, "error.html", nil)
		return
	}

	if len(dares) == 0 {
		s.logRequest(r, http.StatusNotFound, err)
		w.WriteHeader(http.StatusNotFound)
		s.Templates.ExecuteTemplate(w, "no_dares.html", nil)
		return
	}

	s.logRequest(r, http.StatusOK, nil)
	w.WriteHeader(http.StatusOK)
	s.Templates.ExecuteTemplate(w, "all_dares.html", dares)
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
