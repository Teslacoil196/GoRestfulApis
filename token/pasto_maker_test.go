package token

import (
	"TeslaCoil196/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPastoMaker(t *testing.T) {
	maker, err := NewPastoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Second

	issuedAt := time.Now()
	expireAt := time.Now().Add(duration)

	token, payload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, username, payload.Username)
	require.NotEmpty(t, payload.ID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expireAt, payload.ExpiredAt, time.Second)
}

func TestPastoExpiredToken(t *testing.T) {
	maker, err := NewPastoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, Payload, err := maker.CreateToken(util.RandomOwner(), -time.Second)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	Payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, errorTokenExpired.Error())
	require.Nil(t, Payload)
}
