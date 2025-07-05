package keezle

import "errors"

var (
	ErrDefaultSecret    = errors.New("default secret is being used, please set a custom secret")
	ErrProviderDisabled = errors.New("provider is disabled, please enable it in the configuration")
	ErrInvalidEmail     = errors.New("invalid email address, please provide a valid email")
	ErrProviderColon    = errors.New("provider must not contain colons (:)")
	ErrInvalidSessionId = errors.New("invalid session id")
	ErrInvalidPassword  = errors.New("invalid password")
)
