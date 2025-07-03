package models

import "time"

type Key struct {
	ID       string
	UserID   string
	Password string
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
