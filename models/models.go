package models

import "time"

type Account struct {
	ID        string    `json:"id"`
	Provider  string    `json:"provider"`
	AccountID string    `json:"account_id"`
	UserID    string    `json:"user_id"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User[T any] struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	Name          string    `json:"name"`
	Image         string    `json:"image"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Attributes    T         `json:"attributes,omitempty"`
}

type Session[T any] struct {
	ID         string    `json:"id"`
	UserId     string    `json:"user_id"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Token      string    `json:"token"`
	Attributes T         `json:"attributes,omitempty"`
}
