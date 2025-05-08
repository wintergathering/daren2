package daren // Same package as dare.go

import (
	"errors"
	"time"
)

// --- Error definitions for Payback ---
var (
	ErrNoParticipantFound = errors.New("no participant found")
	ErrNoTripFound        = errors.New("no trip found")
	ErrNoPurchaseFound    = errors.New("no purchase entry found")
	ErrNoDebtFound        = errors.New("no debt entry found")
	ErrParticipantExists  = errors.New("participant with that name already exists")
	ErrTripExists         = errors.New("trip with that name already exists")
)

// --- Struct Definitions ---

type Participant struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	// We might add UpdatedAt later if needed
}

type Trip struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	CreatedAt    time.Time      `json:"created_at,omitempty"`
	Participants []*Participant `json:"participants,omitempty"` // For displaying trip details
}

// OriginalPurchase represents a single spending event.
type OriginalPurchase struct {
	ID                 int       `json:"id"`
	TripID             int       `json:"trip_id"`
	PayerParticipantID int       `json:"payer_participant_id"`
	TotalAmount        int       `json:"total_amount"` // In cents
	Description        string    `json:"description"`
	PurchaseDate       time.Time `json:"purchase_date"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
	// For display/API responses, you might want to include:
	PayerName string `json:"payer_name,omitempty"`
	TripName  string `json:"trip_name,omitempty"`
}

// IndividualDebt represents a single person's share of an OriginalPurchase.
type IndividualDebt struct {
	ID                  int       `json:"id"`
	OriginalPurchaseID  int       `json:"original_purchase_id"`
	DebtorParticipantID int       `json:"debtor_participant_id"`
	AmountOwed          int       `json:"amount_owed"` // In cents
	CreatedAt           time.Time `json:"created_at,omitempty"`
	// For display/API responses, you might want to include:
	DebtorName       string            `json:"debtor_name,omitempty"`
	OriginalPurchase *OriginalPurchase `json:"original_purchase,omitempty"` // To show context
}

// --- Service Interface Definition ---

type PaybackService interface {
	// Participant methods
	CreateParticipant(name string) (*Participant, error)
	GetParticipantByID(id int) (*Participant, error)
	GetParticipantByName(name string) (*Participant, error)
	GetAllParticipants() ([]*Participant, error)
	// UpdateParticipant(p *Participant) error // Maybe later
	// DeleteParticipant(id int) error         // Consider implications (cascade deletes?)

	// Trip methods
	CreateTrip(name string) (*Trip, error)
	GetTripByID(id int) (*Trip, error) // Should probably load participants
	GetAllTrips() ([]*Trip, error)     // Should probably load participants for each
	AddParticipantToTrip(tripID, participantID int) error
	RemoveParticipantFromTrip(tripID, participantID int) error
	GetParticipantsForTrip(tripID int) ([]*Participant, error)
	// UpdateTrip(t *Trip) error
	// DeleteTrip(id int) error // Consider implications

	// OriginalPurchase methods
	CreateOriginalPurchase(op *OriginalPurchase, debtorIDs []int) (*OriginalPurchase, []*IndividualDebt, error)
	GetOriginalPurchaseByID(id int) (*OriginalPurchase, error)
	GetOriginalPurchasesByTripID(tripID int) ([]*OriginalPurchase, error)
	// UpdateOriginalPurchase(op *OriginalPurchase, newDebtorIDs []int) (*OriginalPurchase, []*IndividualDebt, error) // Complex
	// DeleteOriginalPurchase(id int) error // Must also delete associated IndividualDebts

	// IndividualDebt methods
	GetIndividualDebtsByPurchaseID(purchaseID int) ([]*IndividualDebt, error)
	GetDebtsForParticipantInTrip(tripID, participantID int) ([]*IndividualDebt, error)   // Debts where they are the debtor
	GetCreditsForParticipantInTrip(tripID, participantID int) ([]*IndividualDebt, error) // Debts where they are the payer (via OriginalPurchase)

	// Balance/Summary methods (more complex, for later)
	// GetTripSummary(tripID int) (map[string]int, error) // participantName -> netBalance
}
