package test

import (
	"context"
	"database/sql"
	"testing"

	db "github.com/cyriljohn147/carpooling/db/sqlc"
	"github.com/cyriljohn147/carpooling/util"
	"github.com/stretchr/testify/require"
)

func createRandomTrip(t *testing.T) db.Trips {
	// First create a user to be the driver
	user := createRandomUser(t)

	arg := db.CreateTripParams{
		DriverID:  user.ID,
		Date:      util.NullTime(util.RandomDate()),
		Time:      util.NullTime(util.RandomTime()),
		TotalCost: util.NullString(util.RandomPrice()),
	}

	trip, err := testQueries.CreateTrip(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trip)

	require.Equal(t, arg.DriverID, trip.DriverID)
	require.True(t, util.TimeEqual(arg.Date, trip.Date), "Date should match")
	require.True(t, util.TimeEqual(arg.Time, trip.Time), "Time should match")
	require.Equal(t, arg.TotalCost, trip.TotalCost)

	require.NotZero(t, trip.ID)

	return trip
}

func TestCreateTrip(t *testing.T) {
	createRandomTrip(t)
}

func TestGetTrip(t *testing.T) {
	trip1 := createRandomTrip(t)
	trip2, err := testQueries.GetTrip(context.Background(), trip1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, trip2)

	require.Equal(t, trip1.ID, trip2.ID)
	require.Equal(t, trip1.DriverID, trip2.DriverID)
	require.True(t, util.TimeEqual(trip1.Date, trip2.Date), "Date should match")
	require.True(t, util.TimeEqual(trip1.Time, trip2.Time), "Time should match")
	require.Equal(t, trip1.TotalCost, trip2.TotalCost)
}

func TestUpdateTrip(t *testing.T) {
	trip1 := createRandomTrip(t)

	// Create another user to be the new driver
	newDriver := createRandomUser(t)

	arg := db.UpdateTripParams{
		ID:        trip1.ID,
		DriverID:  newDriver.ID,
		Date:      util.NullTime(util.RandomDate()),
		Time:      util.NullTime(util.RandomTime()),
		TotalCost: util.NullString(util.RandomPrice()),
	}

	trip2, err := testQueries.UpdateTrip(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trip2)

	require.Equal(t, trip1.ID, trip2.ID)
	require.Equal(t, arg.DriverID, trip2.DriverID)
	require.True(t, util.TimeEqual(arg.Date, trip2.Date), "Date should match")
	require.True(t, util.TimeEqual(arg.Time, trip2.Time), "Time should match")
	require.Equal(t, arg.TotalCost, trip2.TotalCost)
}

func TestDeleteTrip(t *testing.T) {
	trip1 := createRandomTrip(t)
	err := testQueries.DeleteTrip(context.Background(), trip1.ID)
	require.NoError(t, err)

	trip2, err := testQueries.GetTrip(context.Background(), trip1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, trip2)
}

func TestListTrips(t *testing.T) {
	var trips []db.Trips
	for i := 0; i < 10; i++ {
		trips = append(trips, createRandomTrip(t))
	}

	arg := db.ListTripsParams{
		Limit:  5,
		Offset: 5,
	}

	retrievedTrips, err := testQueries.ListTrips(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, retrievedTrips, 5)

	for _, trip := range retrievedTrips {
		require.NotEmpty(t, trip)
	}
}

func TestListTripsEmpty(t *testing.T) {
	arg := db.ListTripsParams{
		Limit:  5,
		Offset: 0,
	}

	trips, err := testQueries.ListTrips(context.Background(), arg)
	require.NoError(t, err)
	require.NotNil(t, trips)

	for _, trip := range trips {
		require.NotEmpty(t, trip)
	}
}

// Test with cancelled context to potentially trigger error paths
func TestListTripsWithCancelledContext(t *testing.T) {
	// Create some trips first
	for i := 0; i < 3; i++ {
		createRandomTrip(t)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	arg := db.ListTripsParams{
		Limit:  10,
		Offset: 0,
	}

	// This should fail due to cancelled context
	_, err := testQueries.ListTrips(ctx, arg)
	require.Error(t, err)
}
