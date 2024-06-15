package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto *paseto.V2
	symKey []byte
}

func NewPasetoMaker(symKey string) (Maker, error) {
	if len(symKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf(
			"invalid key size: must be exactly %d characters",
			chacha20poly1305.KeySize,
		)
	}

	return &PasetoMaker{
		paseto: paseto.NewV2(),
		symKey: []byte(symKey),
	}, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return maker.paseto.Encrypt(maker.symKey, payload, nil)
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symKey, payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
