package test

import (
	"context"
	"database/sql"
	"testing"

	db "github.com/cyriljohn147/carpooling/db/sqlc"
	"github.com/cyriljohn147/carpooling/util"
	"github.com/stretchr/testify/require"
)

func createRandomParticipant(t *testing.T) db.Participants {
	// First create a trip and a user
	trip := createRandomTrip(t)
	user := createRandomUser(t)

	arg := db.AddParticipantParams{
		TripID:      trip.ID,
		UserID:      user.ID,
		ShareAmount: util.NullString(util.RandomPrice()),
	}

	participant, err := testQueries.AddParticipant(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, participant)

	require.Equal(t, arg.TripID, participant.TripID)
	require.Equal(t, arg.UserID, participant.UserID)
	require.Equal(t, arg.ShareAmount, participant.ShareAmount)

	require.NotZero(t, participant.ID)

	return participant
}

func TestAddParticipant(t *testing.T) {
	createRandomParticipant(t)
}

func TestGetParticipant(t *testing.T) {
	participant1 := createRandomParticipant(t)
	participant2, err := testQueries.GetParticipant(context.Background(), participant1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, participant2)

	require.Equal(t, participant1.ID, participant2.ID)
	require.Equal(t, participant1.TripID, participant2.TripID)
	require.Equal(t, participant1.UserID, participant2.UserID)
	require.Equal(t, participant1.ShareAmount, participant2.ShareAmount)
}

func TestUpdateParticipantShare(t *testing.T) {
	participant1 := createRandomParticipant(t)

	arg := db.UpdateParticipantShareParams{
		ID:          participant1.ID,
		ShareAmount: util.NullString(util.RandomPrice()),
	}

	participant2, err := testQueries.UpdateParticipantShare(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, participant2)

	require.Equal(t, participant1.ID, participant2.ID)
	require.Equal(t, participant1.TripID, participant2.TripID)
	require.Equal(t, participant1.UserID, participant2.UserID)
	require.Equal(t, arg.ShareAmount, participant2.ShareAmount)
}

func TestDeleteParticipant(t *testing.T) {
	participant1 := createRandomParticipant(t)
	err := testQueries.DeleteParticipant(context.Background(), participant1.ID)
	require.NoError(t, err)

	participant2, err := testQueries.GetParticipant(context.Background(), participant1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, participant2)
}

func TestListParticipantsByTrip(t *testing.T) {
	// Create a trip first
	trip := createRandomTrip(t)

	var participants []db.Participants
	// Create 10 participants for the same trip
	for i := 0; i < 10; i++ {
		user := createRandomUser(t)
		arg := db.AddParticipantParams{
			TripID:      trip.ID,
			UserID:      user.ID,
			ShareAmount: util.NullString(util.RandomPrice()),
		}

		participant, err := testQueries.AddParticipant(context.Background(), arg)
		require.NoError(t, err)
		participants = append(participants, participant)
	}

	arg := db.ListParticipantsByTripParams{
		TripID: trip.ID,
		Limit:  5,
		Offset: 5,
	}

	retrievedParticipants, err := testQueries.ListParticipantsByTrip(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, retrievedParticipants, 5)

	for _, participant := range retrievedParticipants {
		require.NotEmpty(t, participant)
		require.Equal(t, trip.ID, participant.TripID)
	}
}

func TestListParticipantsByTripEmpty(t *testing.T) {
	trip := createRandomTrip(t)

	arg := db.ListParticipantsByTripParams{
		TripID: trip.ID,
		Limit:  5,
		Offset: 0,
	}

	participants, err := testQueries.ListParticipantsByTrip(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, participants, 0) // Should be empty since no participants added
}

// Test with cancelled context to potentially trigger error paths
func TestListParticipantsByTripWithCancelledContext(t *testing.T) {
	trip := createRandomTrip(t)

	// Create some participants first
	for i := 0; i < 3; i++ {
		user := createRandomUser(t)
		arg := db.AddParticipantParams{
			TripID:      trip.ID,
			UserID:      user.ID,
			ShareAmount: util.NullString(util.RandomPrice()),
		}
		_, err := testQueries.AddParticipant(context.Background(), arg)
		require.NoError(t, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	arg := db.ListParticipantsByTripParams{
		TripID: trip.ID,
		Limit:  10,
		Offset: 0,
	}

	// This should fail due to cancelled context
	_, err := testQueries.ListParticipantsByTrip(ctx, arg)
	require.Error(t, err)
}
