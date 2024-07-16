package query

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/FrostJ143/simplebank/internal/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) *User {
	hashedPassword, err := utils.HashPassword(utils.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		Fullname:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Fullname, user.Fullname)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Fullname, user2.Fullname)
	require.Equal(t, user1.Email, user2.Email)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestUpdateUserOnlyFullname(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullname := utils.RandomOwner()
	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Fullname: sql.NullString{
			String: newFullname,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	require.Equal(t, newFullname, newUser.Fullname)
	require.NotEqual(t, oldUser.Fullname, newUser.Fullname)
	require.Equal(t, oldUser.HashedPassword, newUser.HashedPassword)
	require.Equal(t, oldUser.Email, newUser.Email)
}

func TestUpdateUserOnlyHashedPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	newHashedPassword, err := utils.HashPassword(utils.RandomString(6))
	require.NoError(t, err)

	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	require.Equal(t, newHashedPassword, newUser.HashedPassword)
	require.NotEqual(t, oldUser.HashedPassword, newUser.HashedPassword)
	require.Equal(t, oldUser.Fullname, newUser.Fullname)
	require.Equal(t, oldUser.Email, newUser.Email)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := utils.RandomEmail()

	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	require.Equal(t, newEmail, newUser.Email)
	require.NotEqual(t, oldUser.Email, newUser.Email)
	require.Equal(t, oldUser.Fullname, newUser.Fullname)
	require.Equal(t, oldUser.HashedPassword, newUser.HashedPassword)
}
