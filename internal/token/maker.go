package token

import (
	"errors"
	"time"
)

// Different types of erros returned by VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Maker is a interface for managing tokens
type Maker interface {
	// CreateToken creates a token for a specific username and duration
	CreateToken(string, time.Duration) (string, error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(string) (*Payload, error)
}
