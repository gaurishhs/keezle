package models

import "time"

// DBKey represents a key in the database.
// It contains the ID, UserID, and an optional hashed password.
type DBKey struct {
	// ID is the unique identifier for the key created by concatinating the provider and provider user ID.
	ID *string
	// UserID is the ID of the user associated with this key.
	UserID *string
	// Password is the hashed password for the key.
	Password *string
}

// Key represents a DBKey with a more secure representation.
// It contains the ID, UserID, and a boolean indicating whether the key has a password
type Key struct {
	// ID is the unique identifier for the key.
	ID string
	// UserID is the ID of the user associated with this key.
	UserID string
	// Password indicates whether the key has a password.
	Password bool
}

// User represents a user in the system.
// It contains the ID and attributes associated with the user.
type User[UA AnyStruct] struct {
	ID         string
	Attributes *UA
}

// DBSession represents a session in the database.
// It contains the ID, UserID, active and idle expiration times, and session attributes.
// The ID and UserID are pointers to allow for nil values.
type DBSession[SA AnyStruct] struct {
	ID              *string
	UserId          *string
	ActiveExpiresAt *time.Time
	IdleExpiresAt   *time.Time
	Attributes      *SA
}

// Session represents a user session.
// It contains the ID, associated User, active and idle expiration times, session attributes, and state.
// The User is a pointer to allow for nil values.
type Session[UA, SA AnyStruct] struct {
	// ID is the unique identifier for the session.
	ID string
	// User is the user associated with this session.
	// It is a pointer to allow for nil values.
	User *User[UA]
	// ActiveExpiresAt is the time when the session becomes inactive.
	ActiveExpiresAt time.Time
	// IdleExpiresAt is the time when the session expires due to inactivity.
	IdleExpiresAt time.Time
	// Attributes are the session attributes associated with this session.
	Attributes *SA
	// State indicates the current state of the session i.e. "active" or "idle".
	State string
	// Fresh indicates whether the session is newly created.
	Fresh bool
}
