package keezle

import "errors"

var (
	ErrProviderColon        = errors.New("provider must not contain colons (:)")
	ErrInvalidSessionId     = errors.New("invalid session id")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrInvalidRequestOrigin = errors.New("invalid request origin")
)
