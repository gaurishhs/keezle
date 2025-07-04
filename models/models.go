package models

import "time"

type DBKey struct {
	ID       string
	UserID   string
	Password string
}

type Key struct {
	ID       string
	UserID   string
	Password bool
}

type User[UA AnyStruct] struct {
	ID         string
	Attributes UA
}

type DBSession[SA AnyStruct] struct {
	ID              string
	UserId          string
	ActiveExpiresAt time.Time
	IdleExpiresAt   time.Time
	Attributes      SA
}

type Session[UA, SA AnyStruct] struct {
	ID              string
	User            *User[UA]
	ActiveExpiresAt time.Time
	IdleExpiresAt   time.Time
	Attributes      SA
}
