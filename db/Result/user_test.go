package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/nilesh0729/Notes/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User{
	arg := CreateUserParams{
		Username:     util.RandomUser(),
		Email:        util.RandomEmail(),
		PasswordHash: util.RandomPassword(8),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.PasswordHash, user.PasswordHash)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T){
	CreateRandomUser(t)
}

func TestGetUsers(t *testing.T){
	user1 := CreateRandomUser(t)

	user, err := testQueries.GetUsers(context.Background(), user1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user1.ID, user.ID)
	require.Equal(t, user1.Username, user.Username)
	require.Equal(t, user1.Email, user.Email)
	require.Equal(t, user1.PasswordHash, user.PasswordHash)

	require.WithinDuration(t, user1.CreatedAt.Time, user.CreatedAt.Time, time.Second)
}

func TestListUsers(t *testing.T){

	for i:= 0; i<10; i++{
		CreateRandomUser(t)
	}

	arg := ListUsersParams{
		Limit: 5,
		Offset: 0,
	}
	users, err := testQueries.ListUsers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, users, 5)

	for _,user := range(users){
		require.NotEmpty(t, user)
	}

}

func TestUpdateUser(t *testing.T){
	user1 := CreateRandomUser(t)
	arg := UpdateUsersParams{
		ID: user1.ID,
		PasswordHash: util.RandomPassword(8),
		Email: util.RandomEmail(),
	}

	user, err := testQueries.UpdateUsers(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.ID, user1.ID)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.PasswordHash, user.PasswordHash)
	require.Equal(t, user1.Username, user.Username)

	require.WithinDuration(t, user.CreatedAt.Time, user1.CreatedAt.Time, time.Second)

}

func TestDeleteUser(t *testing.T){
	user1 := CreateRandomUser(t)

	err := testQueries.DeleteUsers(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := testQueries.GetUsers(context.Background(), user1.ID)
	require.EqualError(t, err, sql.ErrNoRows.Error())

	require.Empty(t, user2)
	
}
