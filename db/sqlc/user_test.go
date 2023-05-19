package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createRandomUser(t *testing.T) User {

	arg := CreateUserParams{
		ID:        util.RandomUID(28),
		FirstName: util.RandomString(6),
		LastName:  util.RandomString(6),
		Email:     util.RandomEmail(),
		Phone:     util.RandomPhoneNumber(),
		Age:       util.RandomInt(0, 150),
		Gender:    util.PickRandomGender(),
		Ethnicity: util.PickRandomEthnicity(),
		Nsfw:      util.RandomBool(),
		Metadata:  json.RawMessage([]byte(`{"key": "value","key2": "value2"}`)),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Age, user.Age)
	require.Equal(t, arg.Phone, user.Phone)
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

	require.Equal(t, user1.FirstName, user2.FirstName)
	require.Equal(t, user1.LastName, user2.LastName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Age, user2.Age)
	require.Equal(t, user1.Phone, user2.Phone)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)

	newName := util.RandomName()
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		ID: oldUser.ID,
		FirstName: sql.NullString{
			String: newName,
			Valid:  true,
		},
		LastName: sql.NullString{
			String: newName,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.FirstName, updatedUser.FirstName)
	require.NotEqual(t, oldUser.LastName, updatedUser.LastName)
	require.Equal(t, newName, updatedUser.FirstName)
	require.Equal(t, newName, updatedUser.LastName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := util.RandomEmail()
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		ID: oldUser.ID,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, oldUser.FirstName, updatedUser.FirstName)
	require.Equal(t, oldUser.LastName, updatedUser.LastName)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	newName := util.RandomName()
	newEmail := util.RandomEmail()
	newPhoneNumber := util.RandomPhoneNumber()
	newAge := util.RandomInt(0, 150)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		ID: oldUser.ID,
		FirstName: sql.NullString{
			String: newName,
			Valid:  true,
		},
		LastName: sql.NullString{
			String: newName,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		Phone: sql.NullString{
			String: newPhoneNumber,
			Valid:  true,
		},
		Age: sql.NullInt64{
			Int64: newAge,
			Valid: true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.FirstName, updatedUser.FirstName)
	require.Equal(t, newName, updatedUser.FirstName)
	require.NotEqual(t, oldUser.LastName, updatedUser.LastName)
	require.Equal(t, newName, updatedUser.LastName)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.NotEqual(t, oldUser.Age, updatedUser.Age)
	require.Equal(t, newAge, updatedUser.Age)
}
