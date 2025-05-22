package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	daren "github.com/wintergathering/daren2"
)

// html routes ---------------
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	//set header to html
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	headerStatus := http.StatusOK
	w.WriteHeader(headerStatus)

	s.Templates.ExecuteTemplate(w, "index.html", nil)
}

func (s *Server) handleSuccess(w http.ResponseWriter, r *http.Request) {
	//set header to html
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	headerStatus := http.StatusOK
	w.WriteHeader(headerStatus)

	s.Templates.ExecuteTemplate(w, "dare_success.html", nil)
}

func (s *Server) handleCreateDare(w http.ResponseWriter, r *http.Request) {
	d := &daren.Dare{}

	//check request method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//access form values
	err := r.ParseForm()

	if err != nil {
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
		w.WriteHeader(http.StatusInternalServerError)
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

func (s *Server) handleGetRandomDare(w http.ResponseWriter, r *http.Request) {
	dare, err := s.DareService.GetRandomDare()

	if errors.Is(err, daren.ErrNoDare) {
		w.WriteHeader(http.StatusNotFound)
		s.Templates.ExecuteTemplate(w, "no_dares.html", nil)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.Templates.ExecuteTemplate(w, "error.html", nil)
		return
	}

	err = s.DareService.MarkDareSeen(dare.ID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.Templates.ExecuteTemplate(w, "error.html", nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	s.Templates.ExecuteTemplate(w, "single_dare.html", dare)

}

func (s *Server) handleGetAllDares(w http.ResponseWriter, r *http.Request) {
	dares, err := s.DareService.GetAllDares()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.Templates.ExecuteTemplate(w, "error.html", nil)
		return
	}

	if len(dares) == 0 {
		w.WriteHeader(http.StatusNotFound)
		s.Templates.ExecuteTemplate(w, "no_dares.html", nil)
		return
	}

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

// --- STEE ZONE --- ABANDON ALL HOPE --- QUESTIONABLE CODE HERE ---
func (s *Server) handPaybackHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Welcome to Payback!")
}

// --- Payback API Routes ---

// handlePaybackCreateParticipant handles POST requests to create a new participant.
// Expected JSON body: {"name": "Participant Name"}
func (s *Server) handlePaybackCreateParticipant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	var requestData struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload: " + err.Error()})
		return
	}
	defer r.Body.Close()

	if requestData.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Participant name cannot be empty"})
		return
	}

	participant, err := s.PaybackService.CreateParticipant(requestData.Name)
	if err != nil {
		if errors.Is(err, daren.ErrParticipantExists) {
			writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		} else {
			// Log the full error for server-side debugging
			log.Printf("Error creating participant '%s': %v", requestData.Name, err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create participant"})
		}
		return
	}

	writeJSON(w, http.StatusCreated, participant)
}

// handlePaybackGetAllParticipants handles GET requests to retrieve all participants.
func (s *Server) handlePaybackGetAllParticipants(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	participants, err := s.PaybackService.GetAllParticipants()
	if err != nil {
		// Log the full error for server-side debugging
		log.Printf("Error getting all participants: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve participants"})
		return
	}

	// If no participants, an empty slice is fine and will be marshalled to an empty JSON array []
	writeJSON(w, http.StatusOK, participants)
}

// handlePaybackCreateTrip handles POST requests to create a new trip.
// Expected JSON body: {"name": "Trip Name"}
func (s *Server) handlePaybackCreateTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	var requestData struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload: " + err.Error()})
		return
	}
	defer r.Body.Close()

	if requestData.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Trip name cannot be empty"})
		return
	}

	trip, err := s.PaybackService.CreateTrip(requestData.Name)
	if err != nil {
		// Optional: if you add ErrTripExists in payback.go and payback_service.go
		// if errors.Is(err, daren.ErrTripExists) {
		// 	writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		// } else {
		log.Printf("Error creating trip '%s': %v", requestData.Name, err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create trip"})
		// }
		return
	}

	writeJSON(w, http.StatusCreated, trip)
}

// handlePaybackGetAllTrips handles GET requests to retrieve all trips.
func (s *Server) handlePaybackGetAllTrips(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	trips, err := s.PaybackService.GetAllTrips()
	if err != nil {
		log.Printf("Error getting all trips: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve trips"})
		return
	}
	// Note: The Trip struct in payback.go has an optional Participants field.
	// The current GetAllTrips in payback_service.go doesn't populate this.
	// If you want to include participants here, you'd modify the service method.
	// For now, it will return trips without participant details.
	writeJSON(w, http.StatusOK, trips)
}

// handlePaybackAddParticipantToTrip handles POST requests to add a participant to a trip.
// URL: /api/v1/payback/trips/{tripID}/participants
// Expected JSON body: {"participant_id": 123}
func (s *Server) handlePaybackAddParticipantToTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	tripIDStr := r.PathValue("tripID") // Using PathValue for path parameters
	tripID, err := strconv.Atoi(tripIDStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid trip ID in URL path"})
		return
	}

	var requestData struct {
		ParticipantID int `json:"participant_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload: " + err.Error()})
		return
	}
	defer r.Body.Close()

	if requestData.ParticipantID <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid participant_id"})
		return
	}

	err = s.PaybackService.AddParticipantToTrip(tripID, requestData.ParticipantID)
	if err != nil {
		// More specific error handling could be added here, e.g., if trip or participant not found,
		// or if the participant is already on the trip (SQLite might return a UNIQUE constraint error).
		log.Printf("Error adding participant %d to trip %d: %v", requestData.ParticipantID, tripID, err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to add participant to trip"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Participant added to trip successfully"})
}

// handlePaybackGetParticipantsForTrip handles GET requests to retrieve all participants for a specific trip.
// URL: /api/v1/payback/trips/{tripID}/participants
func (s *Server) handlePaybackGetParticipantsForTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	tripIDStr := r.PathValue("tripID")
	tripID, err := strconv.Atoi(tripIDStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid trip ID in URL path"})
		return
	}

	participants, err := s.PaybackService.GetParticipantsForTrip(tripID)
	if err != nil {
		// Could check for daren.ErrNoTripFound if GetParticipantsForTrip was modified to return it
		// when the trip itself doesn't exist.
		log.Printf("Error getting participants for trip %d: %v", tripID, err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve participants for trip"})
		return
	}

	writeJSON(w, http.StatusOK, participants)
}

// --- Purchase and Debt API Routes ---

// createOriginalPurchaseRequest defines the expected JSON body for creating a purchase.
type createOriginalPurchaseRequest struct {
	TripID             int    `json:"trip_id"`
	PayerParticipantID int    `json:"payer_participant_id"`
	TotalAmount        int    `json:"total_amount"` // In cents
	Description        string `json:"description"`
	PurchaseDate       string `json:"purchase_date"` // Expected format: "YYYY-MM-DDTHH:MM:SSZ" (RFC3339)
	DebtorIDs          []int  `json:"debtor_ids"`
}

// handlePaybackCreateOriginalPurchase handles POST requests to create a new original purchase and its debts.
func (s *Server) handlePaybackCreateOriginalPurchase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	var reqData createOriginalPurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload: " + err.Error()})
		return
	}
	defer r.Body.Close()

	// Basic validation
	if reqData.TripID <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid trip_id"})
		return
	}
	if reqData.PayerParticipantID <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid payer_participant_id"})
		return
	}
	if reqData.TotalAmount <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Total amount must be positive"})
		return
	}
	if reqData.Description == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Description cannot be empty"})
		return
	}
	if len(reqData.DebtorIDs) == 0 {
		// The service handles this by defaulting to payer if empty, but API could enforce it.
		// For now, let service handle it, or uncomment below for stricter API.
		// writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Debtor IDs cannot be empty"})
		// return
	}

	purchaseDateTime, err := time.Parse(time.RFC3339, reqData.PurchaseDate)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid purchase_date format. Expected RFC3339 (e.g., YYYY-MM-DDTHH:MM:SSZ): " + err.Error()})
		return
	}

	op := &daren.OriginalPurchase{
		TripID:             reqData.TripID,
		PayerParticipantID: reqData.PayerParticipantID,
		TotalAmount:        reqData.TotalAmount,
		Description:        reqData.Description,
		PurchaseDate:       purchaseDateTime,
	}

	createdPurchase, createdDebts, err := s.PaybackService.CreateOriginalPurchase(op, reqData.DebtorIDs)
	if err != nil {
		// More specific error handling can be added here based on errors from the service
		log.Printf("Error creating original purchase: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create original purchase"})
		return
	}

	response := struct {
		Purchase *daren.OriginalPurchase `json:"purchase"`
		Debts    []*daren.IndividualDebt `json:"debts"`
	}{
		Purchase: createdPurchase,
		Debts:    createdDebts,
	}

	writeJSON(w, http.StatusCreated, response)
}

// handlePaybackGetPurchasesForTrip handles GET requests to retrieve all purchases for a specific trip.
// URL: /api/v1/payback/trips/{tripID}/purchases
func (s *Server) handlePaybackGetPurchasesForTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	tripIDStr := r.PathValue("tripID")
	tripID, err := strconv.Atoi(tripIDStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid trip ID in URL path"})
		return
	}

	purchases, err := s.PaybackService.GetOriginalPurchasesByTripID(tripID)
	if err != nil {
		log.Printf("Error getting purchases for trip %d: %v", tripID, err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve purchases for trip"})
		return
	}

	writeJSON(w, http.StatusOK, purchases)
}

// handlePaybackGetDebtsForPurchase handles GET requests to retrieve all individual debts for a specific purchase.
// URL: /api/v1/payback/purchases/{purchaseID}/debts
func (s *Server) handlePaybackGetDebtsForPurchase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	purchaseIDStr := r.PathValue("purchaseID")
	purchaseID, err := strconv.Atoi(purchaseIDStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid purchase ID in URL path"})
		return
	}

	debts, err := s.PaybackService.GetIndividualDebtsByPurchaseID(purchaseID)
	if err != nil {
		log.Printf("Error getting debts for purchase %d: %v", purchaseID, err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve debts for purchase"})
		return
	}

	writeJSON(w, http.StatusOK, debts)
}

// handlePaybackRemoveParticipantFromTrip handles DELETE requests to remove a participant from a trip.
// URL: /api/v1/payback/trips/{tripID}/participants/{participantID}
func (s *Server) handlePaybackRemoveParticipantFromTrip(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	tripIDStr := r.PathValue("tripID")
	tripID, err := strconv.Atoi(tripIDStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid trip ID in URL path"})
		return
	}

	participantIDStr := r.PathValue("participantID")
	participantID, err := strconv.Atoi(participantIDStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid participant ID in URL path"})
		return
	}

	err = s.PaybackService.RemoveParticipantFromTrip(tripID, participantID)
	if err != nil {
		// More specific error handling could be added here, e.g., if trip or participant not found,
		// or if the participant wasn't on the trip.
		log.Printf("Error removing participant %d from trip %d: %v", participantID, tripID, err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to remove participant from trip"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Participant removed from trip successfully"})
}

// handlePaybackGetTripBalances handles GET requests to retrieve net balances for all participants in a trip.
// URL: /api/v1/payback/trips/{tripID}/balances
func (s *Server) handlePaybackGetTripBalances(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	tripIDStr := r.PathValue("tripID")
	tripID, err := strconv.Atoi(tripIDStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid trip ID in URL path"})
		return
	}

	balances, err := s.PaybackService.GetTripBalances(tripID)
	if err != nil {
		log.Printf("Error getting balances for trip %d: %v", tripID, err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve balances for trip"})
		return
	}

	writeJSON(w, http.StatusOK, balances)
}

// --- END STEE ZONE ---
