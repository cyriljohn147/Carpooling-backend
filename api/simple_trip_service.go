package api

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	db "github.com/cyriljohn147/carpooling/db/sqlc"
)

// TripService handles simple trip operations without complex transactions
type TripService struct {
	store *db.Queries
	db    *sql.DB
}

// NewTripService creates a new trip service
func NewTripService(store *db.Queries, database *sql.DB) *TripService {
	return &TripService{
		store: store,
		db:    database,
	}
}

// CreateTrip creates a simple trip with automatic date/time
func (s *TripService) CreateTrip(ctx context.Context, req TripRequest) (*TripResponse, error) {
	// Start a database transaction to ensure atomicity
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Create a new store with the transaction
	qtx := s.store.WithTx(tx)

	// Step 1: Validate and get driver user ID
	driverUser, err := qtx.GetUserByName(ctx, sql.NullString{String: req.DriverName, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("driver '%s' not found", req.DriverName)
		}
		return nil, fmt.Errorf("failed to get driver: %w", err)
	}

	// Step 2: Validate and get participant user IDs
	participantUsers := make(map[string]db.Users)
	var participantIDs []int32

	for _, participantName := range req.Participants {
		user, err := qtx.GetUserByName(ctx, sql.NullString{String: participantName, Valid: true})
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("participant '%s' not found", participantName)
			}
			return nil, fmt.Errorf("failed to get participant '%s': %w", participantName, err)
		}
		participantUsers[participantName] = user
		participantIDs = append(participantIDs, user.ID)
	}

	// Step 3: Use current date and time automatically
	now := time.Now()
	tripDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tripTime := time.Date(1, 1, 1, now.Hour(), now.Minute(), now.Second(), 0, now.Location())

	// Step 4: Create the trip
	trip, err := qtx.CreateTrip(ctx, db.CreateTripParams{
		DriverID:  driverUser.ID,
		Date:      sql.NullTime{Time: tripDate, Valid: true},
		Time:      sql.NullTime{Time: tripTime, Valid: true},
		TotalCost: sql.NullString{String: fmt.Sprintf("%.2f", req.TotalCost), Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create trip: %w", err)
	}

	// Step 5: Calculate individual shares
	shareAmount := req.TotalCost / float64(len(req.Participants))

	// Step 6: Add participants with their share amounts
	var participants []ParticipantInfo
	for _, participantName := range req.Participants {
		user := participantUsers[participantName]

		participant, err := qtx.AddParticipant(ctx, db.AddParticipantParams{
			TripID:      trip.ID,
			UserID:      user.ID,
			ShareAmount: sql.NullString{String: fmt.Sprintf("%.2f", shareAmount), Valid: true},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to add participant '%s': %w", participantName, err)
		}

		participants = append(participants, ParticipantInfo{
			UserID:      participant.UserID,
			Name:        participantName,
			ShareAmount: shareAmount,
		})
	}

	// Step 7: Update balances for each participant (except the driver)
	for _, participantName := range req.Participants {
		user := participantUsers[participantName]

		// Skip driver as they don't owe themselves
		if user.ID == driverUser.ID {
			continue
		}

		// Update user's balance by adding their share
		err = qtx.UpdateUserBalance(ctx, db.UpdateUserBalanceParams{
			ID:      user.ID,
			Balance: sql.NullString{String: fmt.Sprintf("%.2f", shareAmount), Valid: true},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update balance for user '%s': %w", participantName, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Step 8: Build response
	response := &TripResponse{
		TripID:       trip.ID,
		DriverID:     trip.DriverID,
		DriverName:   req.DriverName,
		Date:         nullTimeToString(trip.Date),
		Time:         nullTimeToTimeString(trip.Time),
		TotalCost:    req.TotalCost,
		Participants: participants,
		Message:      fmt.Sprintf("Trip created successfully with %d participants. Total cost: %.2f, Share per person: %.2f", len(req.Participants), req.TotalCost, shareAmount),
		CreatedAt:    time.Now().Format(time.RFC3339),
	}

	return response, nil
}

// GetTrip retrieves a trip by ID with participant information
func (s *TripService) GetTrip(ctx context.Context, tripID int32) (*TripResponse, error) {
	// Get trip details
	trip, err := s.store.GetTrip(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}

	// Get driver details
	driver, err := s.store.GetUser(ctx, trip.DriverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get driver: %w", err)
	}

	// Get participants
	participantRows, err := s.store.ListParticipantsByTrip(ctx, db.ListParticipantsByTripParams{
		TripID: tripID,
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}

	var participants []ParticipantInfo

	for _, participant := range participantRows {
		// Get user details for this participant
		user, err := s.store.GetUser(ctx, participant.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get participant user: %w", err)
		}

		shareAmount, _ := strconv.ParseFloat(nullStringToString(participant.ShareAmount), 64)

		participants = append(participants, ParticipantInfo{
			UserID:      participant.UserID,
			Name:        nullStringToString(user.Name),
			ShareAmount: shareAmount,
		})
	}

	totalCost, _ := strconv.ParseFloat(nullStringToString(trip.TotalCost), 64)

	response := &TripResponse{
		TripID:       trip.ID,
		DriverID:     trip.DriverID,
		DriverName:   nullStringToString(driver.Name),
		Date:         nullTimeToString(trip.Date),
		Time:         nullTimeToTimeString(trip.Time),
		TotalCost:    totalCost,
		Participants: participants,
		Message:      "Trip retrieved successfully",
		CreatedAt:    "", // Not available from database in current schema
	}

	return response, nil
}
