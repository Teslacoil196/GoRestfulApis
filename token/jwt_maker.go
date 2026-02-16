package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const minSize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(key string) (Maker, error) {
	if len(key) < minSize {
		return nil, fmt.Errorf("Lenght of key is too short")
	}
	return &JWTMaker{key}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	fmt.Println("JWT-1")
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	fmt.Println("JWT-2")
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {

	Keyfunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errorInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, Keyfunc)
	if err != nil {
		errr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(errr.Inner, errorTokenExpired) {
			return nil, errorTokenExpired
		}
		return nil, errorInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, errorInvalidToken
	}
	return payload, nil
}
