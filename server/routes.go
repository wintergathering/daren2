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

	s.Templates.ExecuteTemplate(w, "index.html", nil)
}

func (s *Server) handleSuccess(w http.ResponseWriter, r *http.Request) {
	//set header to html
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	s.Templates.ExecuteTemplate(w, "dare_success.html", nil)
}

// html routes ---------------
func (s *Server) handleCreateDare(w http.ResponseWriter, r *http.Request) {
	d := &daren.Dare{}

	//check request method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//access form values
	err := r.ParseForm()

	if err != nil {
		s.Templates.ExecuteTemplate(w, "error.html", nil)
		return
	}

	d.Title = r.FormValue("title")
	d.Text = r.FormValue("text")
	d.AddedBy = r.FormValue("addedBy")

	_, err = s.DareService.CreateDare(d)

	//keeping this for now, but will want to return a real page in the future
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//redirect to /success
	http.Redirect(w, r, "/success", http.StatusSeeOther)

}

func (s *Server) handleGetCreateDare(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	s.Templates.ExecuteTemplate(w, "create_dare.html", nil)
}

func (s *Server) HandleGetRandomDare(w http.ResponseWriter, r *http.Request) {
	//TODO
	dare, err := s.DareService.GetRandomDare()

	//TODO -- do something if there are no more dares

	if err != nil {
		s.Templates.ExecuteTemplate(w, "error.html", nil)
		return
	}

	err = s.DareService.MarkDareSeen(dare.ID)

	if err != nil {
		s.Templates.ExecuteTemplate(w, "error.html", nil)
		return
	}

	s.Templates.ExecuteTemplate(w, "single_dare.html", dare)

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
