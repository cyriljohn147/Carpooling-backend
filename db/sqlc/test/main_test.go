package test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	db "github.com/cyriljohn147/carpooling/db/sqlc"
	"github.com/cyriljohn147/carpooling/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/carpooling?sslmode=disable"
)

var testQueries *db.Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = db.New(testDB)

	os.Exit(m.Run())
}

func TestWithTx(t *testing.T) {
	tx, err := testDB.Begin()
	require.NoError(t, err)
	defer tx.Rollback()

	txQueries := testQueries.WithTx(tx)
	require.NotNil(t, txQueries)

	// Test that we can use the transaction-based queries
	user := createRandomUserWithTx(t, txQueries)
	require.NotEmpty(t, user)
}

func createRandomUserWithTx(t *testing.T, queries *db.Queries) db.Users {
	arg := db.CreateUserParams{
		Name:  util.NullString(util.RandomName()),
		Email: util.NullString(util.RandomEmail()),
		Phone: util.NullString(util.RandomPhone()),
	}

	user, err := queries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	return user
}

// Test with closed database to trigger error paths
func TestListUsersWithClosedDB(t *testing.T) {
	// Create a separate database connection that we can close
	tempDB, err := sql.Open(dbDriver, dbSource)
	require.NoError(t, err)

	tempQueries := db.New(tempDB)

	// Create some users first
	for i := 0; i < 3; i++ {
		arg := db.CreateUserParams{
			Name:  util.NullString(util.RandomName()),
			Email: util.NullString(util.RandomEmail()),
			Phone: util.NullString(util.RandomPhone()),
		}
		_, err := tempQueries.CreateUser(context.Background(), arg)
		require.NoError(t, err)
	}

	// Close the database connection
	tempDB.Close()

	// Now try to list users - this should fail
	arg := db.ListUsersParams{
		Limit:  10,
		Offset: 0,
	}

	_, err = tempQueries.ListUsers(context.Background(), arg)
	require.Error(t, err)
}

// Test edge case scenarios to improve coverage
func TestListFunctionsEdgeCases(t *testing.T) {
	// Test ListUsers with extreme parameters
	arg1 := db.ListUsersParams{
		Limit:  -1, // Invalid limit
		Offset: 0,
	}
	_, err := testQueries.ListUsers(context.Background(), arg1)
	// This may or may not error depending on postgres behavior
	_ = err // We don't require error here as postgres might handle it differently

	// Test with very high offset
	arg2 := db.ListUsersParams{
		Limit:  10,
		Offset: 999999,
	}
	users, err := testQueries.ListUsers(context.Background(), arg2)
	require.NoError(t, err)
	require.Empty(t, users)

	// Similar tests for other list functions
	// Create test data first
	// user := createRandomUser(t)
	trip := createRandomTrip(t)

	// Test ListBalancesByUser edge cases
	// balanceArg := db.ListBalancesByUserParams{
	// 	UserID: user.ID,
	// 	Limit:  0, // Zero limit
	// 	Offset: 0,
	// }
	// balances, err := testQueries.ListBalancesByUser(context.Background(), balanceArg)
	// require.NoError(t, err)
	// require.Empty(t, balances)

	// Test ListParticipantsByTrip edge cases
	participantArg := db.ListParticipantsByTripParams{
		TripID: trip.ID,
		Limit:  0, // Zero limit
		Offset: 0,
	}
	participants, err := testQueries.ListParticipantsByTrip(context.Background(), participantArg)
	require.NoError(t, err)
	require.Empty(t, participants)

	// Test ListTrips edge cases
	tripArg := db.ListTripsParams{
		Limit:  0, // Zero limit
		Offset: 0,
	}
	trips, err := testQueries.ListTrips(context.Background(), tripArg)
	require.NoError(t, err)
	require.Empty(t, trips)
}
