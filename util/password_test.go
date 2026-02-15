package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	pass := RandomString(7)

	hashedPass1, err := HashedPassword(pass)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPass1)

	err = CheckPasswords(pass, hashedPass1)
	require.NoError(t, err)

	wrongPass := RandomString(7)
	err = CheckPasswords(wrongPass, hashedPass1)
	require.Equal(t, err, bcrypt.ErrMismatchedHashAndPassword)

	hashedPass2, err := HashedPassword(hashedPass1)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPass2)
	require.NotEqual(t, hashedPass1, hashedPass2)
}
