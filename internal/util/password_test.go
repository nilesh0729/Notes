package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashedPassword(t *testing.T) {
	password := RandomString(8)
	hashedPassword, err := HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = CheckPassword(password, hashedPassword)
	require.NoError(t, err)

	wrongPassword := RandomString(8)
	hell := CheckPassword(wrongPassword, hashedPassword)
	require.EqualError(t, hell, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword2, err := HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword, hashedPassword2)
}
