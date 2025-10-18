package test

import (
	"context"
	"database/sql"
	"testing"

	db "github.com/cyriljohn147/carpooling/db/sqlc"
	"github.com/cyriljohn147/carpooling/util"
	"github.com/stretchr/testify/require"
)

func createRandomBalance(t *testing.T) db.Balances {
	// First create a user
	user := createRandomUser(t)

	arg := db.CreateBalanceParams{
		UserID:     user.ID,
		AmountOwed: util.NullString(util.RandomPrice()),
		Week:       util.NullInt32(util.RandomWeek()),
	}

	balance, err := testQueries.CreateBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, balance)

	require.Equal(t, arg.UserID, balance.UserID)
	require.Equal(t, arg.AmountOwed, balance.AmountOwed)
	require.Equal(t, arg.Week, balance.Week)

	require.NotZero(t, balance.ID)
	require.NotZero(t, balance.CreatedAt)

	return balance
}

func TestCreateBalance(t *testing.T) {
	createRandomBalance(t)
}

func TestGetBalance(t *testing.T) {
	balance1 := createRandomBalance(t)
	balance2, err := testQueries.GetBalance(context.Background(), balance1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, balance2)

	require.Equal(t, balance1.ID, balance2.ID)
	require.Equal(t, balance1.UserID, balance2.UserID)
	require.Equal(t, balance1.AmountOwed, balance2.AmountOwed)
	require.Equal(t, balance1.Week, balance2.Week)
	require.Equal(t, balance1.CreatedAt, balance2.CreatedAt)
}

func TestUpdateBalanceAmount(t *testing.T) {
	balance1 := createRandomBalance(t)

	arg := db.UpdateBalanceAmountParams{
		ID:         balance1.ID,
		AmountOwed: util.NullString(util.RandomPrice()),
	}

	balance2, err := testQueries.UpdateBalanceAmount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, balance2)

	require.Equal(t, balance1.ID, balance2.ID)
	require.Equal(t, balance1.UserID, balance2.UserID)
	require.Equal(t, arg.AmountOwed, balance2.AmountOwed)
	require.Equal(t, balance1.Week, balance2.Week)
	require.Equal(t, balance1.CreatedAt, balance2.CreatedAt)
}

func TestDeleteBalance(t *testing.T) {
	balance1 := createRandomBalance(t)
	err := testQueries.DeleteBalance(context.Background(), balance1.ID)
	require.NoError(t, err)

	balance2, err := testQueries.GetBalance(context.Background(), balance1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, balance2)
}

func TestListBalancesByUser(t *testing.T) {
	// Create a user first
	user := createRandomUser(t)

	var balances []db.Balances
	// Create 10 balances for the same user
	for i := 0; i < 10; i++ {
		arg := db.CreateBalanceParams{
			UserID:     user.ID,
			AmountOwed: util.NullString(util.RandomPrice()),
			Week:       util.NullInt32(util.RandomWeek()),
		}

		balance, err := testQueries.CreateBalance(context.Background(), arg)
		require.NoError(t, err)
		balances = append(balances, balance)
	}

	arg := db.ListBalancesByUserParams{
		UserID: user.ID,
		Limit:  5,
		Offset: 5,
	}

	retrievedBalances, err := testQueries.ListBalancesByUser(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, retrievedBalances, 5)

	for _, balance := range retrievedBalances {
		require.NotEmpty(t, balance)
		require.Equal(t, user.ID, balance.UserID)
	}
}

func TestListBalancesByUserEmpty(t *testing.T) {
	user := createRandomUser(t)

	arg := db.ListBalancesByUserParams{
		UserID: user.ID,
		Limit:  5,
		Offset: 0,
	}

	balances, err := testQueries.ListBalancesByUser(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, balances, 0) // Should be empty since no balances added
}

// Test with cancelled context to potentially trigger error paths
func TestListBalancesByUserWithCancelledContext(t *testing.T) {
	user := createRandomUser(t)

	// Create some balances first
	for i := 0; i < 3; i++ {
		arg := db.CreateBalanceParams{
			UserID:     user.ID,
			AmountOwed: util.NullString(util.RandomPrice()),
			Week:       util.NullInt32(util.RandomWeek()),
		}
		_, err := testQueries.CreateBalance(context.Background(), arg)
		require.NoError(t, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	arg := db.ListBalancesByUserParams{
		UserID: user.ID,
		Limit:  10,
		Offset: 0,
	}

	// This should fail due to cancelled context
	_, err := testQueries.ListBalancesByUser(ctx, arg)
	require.Error(t, err)
}
