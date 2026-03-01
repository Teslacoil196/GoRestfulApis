package token

import (
	"TeslaCoil196/util"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Second

	issuedAt := time.Now()
	expireAt := time.Now().Add(duration)

	token, _, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, username, payload.Username)
	require.NotEmpty(t, payload.ID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expireAt, payload.ExpiredAt, time.Second)
}

func TestExpiredToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	token, _, err := maker.CreateToken(util.RandomOwner(), -time.Second)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	Payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, errorTokenExpired.Error())
	require.Nil(t, Payload)
}

func TestInvalidToken(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), time.Second)
	require.NoError(t, err)

	JwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := JwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, errorInvalidToken.Error())
	require.Nil(t, payload)
}
