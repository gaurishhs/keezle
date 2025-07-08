package keezle

import (
	"time"

	"github.com/gaurishhs/keezle/adapters"
	"github.com/gaurishhs/keezle/logger"
	"github.com/gaurishhs/keezle/models"
	"github.com/gaurishhs/keezle/utils"
)

type SessionCookieConfig struct {
	Expires bool
	Name    string
	Secure  bool
}

type SessionConfig struct {
	ActivePeriod time.Duration
	IdlePeriod   time.Duration
	Cookie       *SessionCookieConfig
}

type Config[UA, SA models.AnyStruct] struct {
	Adapter                adapters.Adapter[UA, SA]
	Session                *SessionConfig
	Logger                 logger.Logger
	Hash                   func(string) (string, error)
	ComparePasswordAndHash func(string, string) (bool, error)
	GetUserAttributes      func(user *models.User[UA]) (*UA, error)
	GetSessionAttributes   func(dbSession *models.DBSession[SA]) (*SA, error)
}

type Keezle[UA, SA models.AnyStruct] struct {
	Config *Config[UA, SA]
}

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
