// Package keezle simplifies user authentication and session management in Go applications.
// It allows you to manage user sessions, authenticate users, and handle user attributes.
// It is designed to be flexible and extensible, allowing you to use your own database adapters and
// customize the session management behavior.
package keezle

import (
	"time"

	"github.com/gaurishhs/keezle/adapters"
	"github.com/gaurishhs/keezle/logger"
	"github.com/gaurishhs/keezle/models"
	"github.com/gaurishhs/keezle/utils"
)

// SessionCookieConfig defines the configuration for session cookies.
type SessionCookieConfig struct {
	Expires bool
	Name    string
	Secure  bool
}

// SessionConfig defines the configuration for user sessions.
type SessionConfig struct {
	ActivePeriod time.Duration
	IdlePeriod   time.Duration
	Cookie       *SessionCookieConfig
}

type CSRFProtectionConfig struct {
	Host              string
	AllowedSubdomains []string
}

// Config defines the configuration for the Keezle instance.
type Config[UA, SA models.AnyStruct] struct {
	Adapter                adapters.Adapter[UA, SA]
	Session                *SessionConfig
	Logger                 logger.Logger
	Hash                   func(string) (string, error)
	ComparePasswordAndHash func(string, string) (bool, error)
	GetUserAttributes      func(user *models.User[UA]) (*UA, error)
	GetSessionAttributes   func(dbSession *models.DBSession[SA]) (*SA, error)
	CSRF                   *CSRFProtectionConfig
}

// Keezle is the main struct that holds the configuration and provides methods for authentication and session management.
// It is parameterized by user attributes (UA) and session attributes (SA).
// The user attributes represent the data associated with a user, while the session attributes represent the data
// associated with a session.
type Keezle[UA, SA models.AnyStruct] struct {
	Config *Config[UA, SA]
}

// New creates a new instance of Keezle with the provided configuration.
func New[UA, SA models.AnyStruct](config *Config[UA, SA]) (res *Keezle[UA, SA]) {
	res = &Keezle[UA, SA]{
		Config: config,
	}

	if res.Config.Adapter == nil {
		panic("adapter is required")
	}

	if res.Config.Logger == nil {
		res.Config.Logger = logger.NoOpLogger
	}

	if res.Config.Hash == nil {
		res.Config.Logger.Log("debug: hash function is not set, using default hash function")
		res.Config.Hash = utils.HashPassword
	}

	if res.Config.ComparePasswordAndHash == nil {
		res.Config.Logger.Log("debug: compare password and hash function is not set, using default compare password and hash function")
		res.Config.ComparePasswordAndHash = utils.ComparePasswordAndHash
	}

	if res.Config.Session == nil {
		res.Config.Session = &SessionConfig{
			ActivePeriod: time.Hour * 24,
			IdlePeriod:   time.Hour * 24 * 14,
		}
	}

	return nil
}
