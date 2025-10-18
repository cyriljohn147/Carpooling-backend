package test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	db "github.com/cyriljohn147/carpooling/db/sqlc"
	"github.com/cyriljohn147/carpooling/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) db.Users {
	arg := db.CreateUserParams{
		Name:  util.NullString(util.RandomName()),
		Email: util.NullString(util.RandomEmail()),
		Phone: util.NullString(util.RandomPhone()),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Phone, user.Phone)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Phone, user2.Phone)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)

	arg := db.UpdateUserParams{
		ID:    user1.ID,
		Name:  util.NullString(util.RandomName()),
		Email: util.NullString(util.RandomEmail()),
		Phone: util.NullString(util.RandomPhone()),
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, arg.Name, user2.Name)
	require.Equal(t, arg.Email, user2.Email)
	require.Equal(t, arg.Phone, user2.Phone)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}

func TestListUsers(t *testing.T) {
	var users []db.Users
	for i := 0; i < 10; i++ {
		users = append(users, createRandomUser(t))
	}

	arg := db.ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	retrievedUsers, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, retrievedUsers, 5)

	for _, user := range retrievedUsers {
		require.NotEmpty(t, user)
	}
}

func TestListUsersEmpty(t *testing.T) {
	arg := db.ListUsersParams{
		Limit:  5,
		Offset: 0,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotNil(t, users)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}

func TestCreateUserWithMinimalData(t *testing.T) {
	arg := db.CreateUserParams{
		Name:  util.NullString(util.RandomName()),
		Email: sql.NullString{}, // Empty email
		Phone: sql.NullString{}, // Empty phone
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Phone, user.Phone)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
}

func TestUpdateUserPartial(t *testing.T) {
	user1 := createRandomUser(t)

	// Update only name and email, keep others same
	arg := db.UpdateUserParams{
		ID:    user1.ID,
		Name:  util.NullString(util.RandomName()),
		Email: util.NullString(util.RandomEmail()),
		Phone: user1.Phone,
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, arg.Name, user2.Name)
	require.Equal(t, arg.Email, user2.Email)
	require.Equal(t, user1.Phone, user2.Phone)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}

// Test with cancelled context to potentially trigger error paths
func TestListUsersWithCancelledContext(t *testing.T) {
	// Create some users first
	for i := 0; i < 3; i++ {
		createRandomUser(t)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	arg := db.ListUsersParams{
		Limit:  10,
		Offset: 0,
	}

	// This should fail due to cancelled context
	_, err := testQueries.ListUsers(ctx, arg)
	require.Error(t, err)
}

// Test with very large limit to potentially trigger edge cases
func TestListUsersLargeLimit(t *testing.T) {
	// Create some users first
	for i := 0; i < 5; i++ {
		createRandomUser(t)
	}

	arg := db.ListUsersParams{
		Limit:  1000000, // Very large limit
		Offset: 0,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(users), 5) // At least the users we created
}
