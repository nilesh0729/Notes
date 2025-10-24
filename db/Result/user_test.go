package Database

import (
	"context"
	"testing"

	"github.com/nilesh0729/Notes/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUse(t *testing.T) {
	arg := CreateUserParams{
		Username: util.RandomString(5),
		HashedPassword: util.RandomString(8),
		Email: util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.Username, arg.Username)
	require.Equal(t, user.HashedPassword, arg.HashedPassword)
	require.Equal(t, user.Email, arg.Email)

}

func TestGetUser(t *testing.T){

	user1 := RandomUser(t)
	user2, err := testQueries.GetUser(context.Background(),user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Email, user2.Email)
}

func RandomUser(t *testing.T)User{

	arg := CreateUserParams{
		Username: util.RandomString(5),
		HashedPassword: util.RandomString(8),
		Email: util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(),arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.Username, arg.Username)
	require.Equal(t, user.HashedPassword, arg.HashedPassword)
	require.Equal(t, user.Email, arg.Email)

	return user
}