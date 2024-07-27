package query

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	CreatedAt time.Time `json:"created_at"`
	Owner     string    `json:"owner"`
	Currenncy string    `json:"currency"`
	ID        int64     `json:"id"`
	Balance   int64     `json:"balance"`
}

type Entry struct {
	CreatedAt  time.Time `json:"created_at"`
	ID         int64     `json:"id"`
	Account_ID int64     `json:"account_id"`
	Amount     int64     `json:"amount"`
}

type Transfer struct {
	CreatedAt       time.Time `json:"created_at"`
	ID              int64     `json:"id"`
	From_Account_ID int64     `json:"from_account_id"`
	To_Account_ID   int64     `json:"to_account_id"`
	Amount          int64     `json:"amount"`
}

type User struct {
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
	Username          string    `json:"username"`
	HashedPassword    string    `json:"hashed_password"`
	Fullname          string    `json:"fullname"`
	Email             string    `json:"email"`
	IsEmailVerified   bool      `json:"is_email_verified"`
}

type Session struct {
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	Username     string    `json:"username"`
	RefrestToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIP     string    `json:"client_ip"`
	ID           uuid.UUID `json:"id"`
	IsBlocked    bool      `json:"is_blocked"`
}

type VerifyEmail struct {
	CreatedAt  time.Time `json:"created_at"`
	ExpiredAt  time.Time `json:"expired_at"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	SecretCode string    `json:"secret_code"`
	ID         int64     `json:"id"`
	IsUsed     bool      `json:"is_used"`
}
