package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PastoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPastoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("Incorrect key size: must be exactly %d charchers long", chacha20poly1305.KeySize)
	}

	maker := &PastoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

func (maker *PastoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

func (maker *PastoMaker) VerifyToken(token string) (*Payload, error) {

	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, err
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil

}
