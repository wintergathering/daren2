package sqlite // This service implementation lives in the sqlite package

import (
	"database/sql"
	"fmt" // For error wrapping

	// For CreateOriginalPurchase debtorIDs handling (potential)
	// Import the daren package to access the structs and interface
	// The module path is 'github.com/wintergathering/daren2'
	// So, to access types like daren.Participant, daren.PaybackService:
	daren "github.com/wintergathering/daren2"
)

// paybackService provides an implementation of the daren.PaybackService interface
// using a SQL database.
type paybackService struct {
	db *sql.DB // The database connection for payback.db
}

// NewPaybackService creates a new instance of paybackService.
func NewPaybackService(db *sql.DB) daren.PaybackService {
	return &paybackService{
		db: db,
	}
}

// --- Participant Methods ---

// CreateParticipant adds a new participant to the database.
// It ensures the participant name is unique.
func (s *paybackService) CreateParticipant(name string) (*daren.Participant, error) {
	// Check if participant with the same name already exists
	existingParticipant, err := s.GetParticipantByName(name)
	if err != nil && err != daren.ErrNoParticipantFound {
		return nil, fmt.Errorf("error checking for existing participant: %w", err)
	}
	if existingParticipant != nil {
		return nil, daren.ErrParticipantExists // Use the error defined in daren/payback.go
	}

	query := `INSERT INTO participants (name) VALUES (?) RETURNING participant_id, name, created_at;`
	// Note: SQLite's RETURNING clause is relatively new. If using an older version,
	// you'd do an INSERT then a SELECT. modernc.org/sqlite supports it.

	var p daren.Participant
	err = s.db.QueryRow(query, name).Scan(&p.ID, &p.Name, &p.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("could not create participant %s: %w", name, err)
	}
	return &p, nil
}

// GetParticipantByID retrieves a participant by their ID.
func (s *paybackService) GetParticipantByID(id int) (*daren.Participant, error) {
	query := `SELECT participant_id, name, created_at FROM participants WHERE participant_id = ?;`
	var p daren.Participant
	err := s.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, daren.ErrNoParticipantFound
		}
		return nil, fmt.Errorf("could not get participant by id %d: %w", id, err)
	}
	return &p, nil
}

// GetParticipantByName retrieves a participant by their name.
func (s *paybackService) GetParticipantByName(name string) (*daren.Participant, error) {
	query := `SELECT participant_id, name, created_at FROM participants WHERE name = ?;`
	var p daren.Participant
	// Using QueryRow because name is UNIQUE. If it wasn't, you'd use Query.
	err := s.db.QueryRow(query, name).Scan(&p.ID, &p.Name, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, daren.ErrNoParticipantFound
		}
		return nil, fmt.Errorf("could not get participant by name %s: %w", name, err)
	}
	return &p, nil
}

// GetAllParticipants retrieves all participants.
func (s *paybackService) GetAllParticipants() ([]*daren.Participant, error) {
	query := `SELECT participant_id, name, created_at FROM participants ORDER BY name;`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not get all participants: %w", err)
	}
	defer rows.Close()

	var participants []*daren.Participant
	for rows.Next() {
		var p daren.Participant
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning participant row: %w", err)
		}
		participants = append(participants, &p)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating participant rows: %w", err)
	}
	return participants, nil
}

// --- Trip Methods ---

// CreateTrip adds a new trip to the database.
func (s *paybackService) CreateTrip(name string) (*daren.Trip, error) {
	// Optional: Check if trip with the same name already exists
	// For simplicity, we'll skip this check for now, but you might want it.
	// If you add it, define and use daren.ErrTripExists.

	query := `INSERT INTO trips (name) VALUES (?) RETURNING trip_id, name, created_at;`
	var t daren.Trip
	err := s.db.QueryRow(query, name).Scan(&t.ID, &t.Name, &t.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("could not create trip %s: %w", name, err)
	}
	// Participants field will be empty initially.
	return &t, nil
}

// GetTripByID retrieves a trip by its ID.
// This version does NOT YET load participants for the trip. We can add that later.
func (s *paybackService) GetTripByID(id int) (*daren.Trip, error) {
	query := `SELECT trip_id, name, created_at FROM trips WHERE trip_id = ?;`
	var t daren.Trip
	err := s.db.QueryRow(query, id).Scan(&t.ID, &t.Name, &t.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, daren.ErrNoTripFound
		}
		return nil, fmt.Errorf("could not get trip by id %d: %w", id, err)
	}
	// To load participants:
	// t.Participants, err = s.GetParticipantsForTrip(t.ID)
	// if err != nil { return nil, fmt.Errorf("could not load participants for trip %d: %w", t.ID, err) }
	return &t, nil
}

// GetAllTrips retrieves all trips.
// This version does NOT YET load participants for each trip.
func (s *paybackService) GetAllTrips() ([]*daren.Trip, error) {
	query := `SELECT trip_id, name, created_at FROM trips ORDER BY name;`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not get all trips: %w", err)
	}
	defer rows.Close()

	var trips []*daren.Trip
	for rows.Next() {
		var t daren.Trip
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning trip row: %w", err)
		}
		// To load participants for each trip:
		// t.Participants, err = s.GetParticipantsForTrip(t.ID)
		// if err != nil { /* handle error, maybe log and continue */ }
		trips = append(trips, &t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating trip rows: %w", err)
	}
	return trips, nil
}

// AddParticipantToTrip links a participant to a trip in the trip_participants table.
func (s *paybackService) AddParticipantToTrip(tripID, participantID int) error {
	// Optional: Check if trip and participant exist first.
	// Optional: Check if the link already exists to avoid duplicate errors or handle them.
	query := `INSERT INTO trip_participants (trip_id, participant_id) VALUES (?, ?);`
	_, err := s.db.Exec(query, tripID, participantID)
	if err != nil {
		// You might want to check for specific SQLite errors, e.g., UNIQUE constraint violation
		// if the pair already exists, or FOREIGN KEY constraint violation if IDs are invalid.
		return fmt.Errorf("could not add participant %d to trip %d: %w", participantID, tripID, err)
	}
	return nil
}

// RemoveParticipantFromTrip unlinks a participant from a trip.
func (s *paybackService) RemoveParticipantFromTrip(tripID, participantID int) error {
	query := `DELETE FROM trip_participants WHERE trip_id = ? AND participant_id = ?;`
	res, err := s.db.Exec(query, tripID, participantID)
	if err != nil {
		return fmt.Errorf("could not remove participant %d from trip %d: %w", participantID, tripID, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected for removing participant from trip: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no link found for participant %d in trip %d to remove", participantID, tripID) // Or a specific error
	}
	return nil
}

// GetParticipantsForTrip retrieves all participants associated with a given tripID.
func (s *paybackService) GetParticipantsForTrip(tripID int) ([]*daren.Participant, error) {
	query := `
		SELECT p.participant_id, p.name, p.created_at
		FROM participants p
		JOIN trip_participants tp ON p.participant_id = tp.participant_id
		WHERE tp.trip_id = ?
		ORDER BY p.name;
	`
	rows, err := s.db.Query(query, tripID)
	if err != nil {
		return nil, fmt.Errorf("could not get participants for trip %d: %w", tripID, err)
	}
	defer rows.Close()

	var participants []*daren.Participant
	for rows.Next() {
		var p daren.Participant
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning participant row for trip: %w", err)
		}
		participants = append(participants, &p)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating participant rows for trip: %w", err)
	}
	return participants, nil
}

// --- OriginalPurchase and IndividualDebt Methods (Basic stubs for now) ---

// CreateOriginalPurchase creates an original purchase and its associated individual debts.
// This is a more complex method involving a transaction.
func (s *paybackService) CreateOriginalPurchase(op *daren.OriginalPurchase, debtorIDs []int) (*daren.OriginalPurchase, []*daren.IndividualDebt, error) {
	// This method needs to be transactional to ensure atomicity.
	tx, err := s.db.Begin()
	if err != nil {
		return nil, nil, fmt.Errorf("could not begin transaction: %w", err)
	}
	// Defer a rollback in case of error, commit if successful.
	// This is a common pattern for transactions.
	defer func() {
		if p := recover(); p != nil { // Handle panics during transaction
			_ = tx.Rollback()
			panic(p) // Re-panic after rollback
		} else if err != nil { // Handle explicit errors
			_ = tx.Rollback() // Rollback on error
		} else {
			err = tx.Commit() // Commit on success
			if err != nil {
				// Log or handle commit error
				err = fmt.Errorf("error committing transaction: %w", err)
			}
		}
	}()

	// 1. Insert the OriginalPurchase
	purchaseQuery := `
		INSERT INTO original_purchases (trip_id, payer_participant_id, total_amount, description, purchase_date)
		VALUES (?, ?, ?, ?, ?)
		RETURNING purchase_id, created_at;
	`
	// We assume op.PurchaseDate is already set correctly (e.g., to time.Now() or user input)
	err = tx.QueryRow(purchaseQuery, op.TripID, op.PayerParticipantID, op.TotalAmount, op.Description, op.PurchaseDate).Scan(&op.ID, &op.CreatedAt)
	if err != nil {
		return nil, nil, fmt.Errorf("could not insert original purchase: %w", err)
	}

	// 2. Calculate and Insert IndividualDebts
	if len(debtorIDs) == 0 {
		// Or handle as an error: return nil, nil, errors.New("no debtors specified for the purchase")
		// For now, if no debtors, assume only the payer is involved (their share is the full amount)
		debtorIDs = []int{op.PayerParticipantID}
	}

	amountPerDebtor := op.TotalAmount / len(debtorIDs) // Simple equal split
	// Handle remainder if TotalAmount doesn't divide evenly (e.g., add to first debtor or payer)
	remainder := op.TotalAmount % len(debtorIDs)

	var createdDebts []*daren.IndividualDebt
	debtQuery := `
		INSERT INTO individual_debts (original_purchase_id, debtor_participant_id, amount_owed)
		VALUES (?, ?, ?)
		RETURNING debt_id, created_at;
	`
	stmt, err := tx.Prepare(debtQuery)
	if err != nil {
		return nil, nil, fmt.Errorf("could not prepare individual debt statement: %w", err)
	}
	defer stmt.Close()

	for i, debtorID := range debtorIDs {
		currentAmountOwed := amountPerDebtor
		if i == 0 && remainder > 0 { // Distribute remainder to the first debtor (can be payer)
			currentAmountOwed += remainder
		}

		var debt daren.IndividualDebt
		debt.OriginalPurchaseID = op.ID
		debt.DebtorParticipantID = debtorID
		debt.AmountOwed = currentAmountOwed

		err = stmt.QueryRow(op.ID, debtorID, currentAmountOwed).Scan(&debt.ID, &debt.CreatedAt)
		if err != nil {
			// Error is handled by the defer func which will rollback
			return nil, nil, fmt.Errorf("could not insert individual debt for participant %d: %w", debtorID, err)
		}
		createdDebts = append(createdDebts, &debt)
	}

	// If we reach here without error, the defer func will commit.
	return op, createdDebts, err // err will be nil on success, or commit error
}

// GetOriginalPurchaseByID retrieves an original purchase by its ID.
func (s *paybackService) GetOriginalPurchaseByID(id int) (*daren.OriginalPurchase, error) {
	query := `
		SELECT purchase_id, trip_id, payer_participant_id, total_amount, description, purchase_date, created_at
		FROM original_purchases
		WHERE purchase_id = ?;
	`
	var op daren.OriginalPurchase
	err := s.db.QueryRow(query, id).Scan(
		&op.ID, &op.TripID, &op.PayerParticipantID, &op.TotalAmount,
		&op.Description, &op.PurchaseDate, &op.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, daren.ErrNoPurchaseFound
		}
		return nil, fmt.Errorf("could not get original purchase by id %d: %w", id, err)
	}
	return &op, nil
}

// GetOriginalPurchasesByTripID retrieves all original purchases for a given tripID.
func (s *paybackService) GetOriginalPurchasesByTripID(tripID int) ([]*daren.OriginalPurchase, error) {
	query := `
		SELECT purchase_id, trip_id, payer_participant_id, total_amount, description, purchase_date, created_at
		FROM original_purchases
		WHERE trip_id = ?
		ORDER BY purchase_date DESC, created_at DESC;
	`
	rows, err := s.db.Query(query, tripID)
	if err != nil {
		return nil, fmt.Errorf("could not get original purchases for trip %d: %w", tripID, err)
	}
	defer rows.Close()

	var purchases []*daren.OriginalPurchase
	for rows.Next() {
		var op daren.OriginalPurchase
		if err := rows.Scan(
			&op.ID, &op.TripID, &op.PayerParticipantID, &op.TotalAmount,
			&op.Description, &op.PurchaseDate, &op.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning original purchase row for trip: %w", err)
		}
		purchases = append(purchases, &op)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating original purchase rows for trip: %w", err)
	}
	return purchases, nil
}

// GetIndividualDebtsByPurchaseID retrieves all individual debts for a given original_purchase_id.
func (s *paybackService) GetIndividualDebtsByPurchaseID(purchaseID int) ([]*daren.IndividualDebt, error) {
	query := `
		SELECT debt_id, original_purchase_id, debtor_participant_id, amount_owed, created_at
		FROM individual_debts
		WHERE original_purchase_id = ?;
	`
	rows, err := s.db.Query(query, purchaseID)
	if err != nil {
		return nil, fmt.Errorf("could not get individual debts for purchase %d: %w", purchaseID, err)
	}
	defer rows.Close()

	var debts []*daren.IndividualDebt
	for rows.Next() {
		var d daren.IndividualDebt
		if err := rows.Scan(
			&d.ID, &d.OriginalPurchaseID, &d.DebtorParticipantID, &d.AmountOwed, &d.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning individual debt row: %w", err)
		}
		debts = append(debts, &d)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating individual debt rows: %w", err)
	}
	return debts, nil
}

// --- Stubs for remaining PaybackService methods (to satisfy the interface) ---
// We will implement these properly later.

func (s *paybackService) GetDebtsForParticipantInTrip(tripID, participantID int) ([]*daren.IndividualDebt, error) {
	// TODO: Implement this method
	// Query individual_debts where trip_id (via original_purchases) and debtor_participant_id match
	return nil, fmt.Errorf("GetDebtsForParticipantInTrip not yet implemented")
}

func (s *paybackService) GetCreditsForParticipantInTrip(tripID, participantID int) ([]*daren.IndividualDebt, error) {
	// TODO: Implement this method
	// This is a bit more complex: find original_purchases where participantID is the payer,
	// then find all individual_debts for those purchases where the debtor is NOT the payer.
	return nil, fmt.Errorf("GetCreditsForParticipantInTrip not yet implemented")
}

// Ensure paybackService implements daren.PaybackService
var _ daren.PaybackService = (*paybackService)(nil)
