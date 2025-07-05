package keezle

import (
	"time"

	"github.com/gaurishhs/keezle/adapters"
	"github.com/gaurishhs/keezle/logger"
	"github.com/gaurishhs/keezle/models"
	"github.com/gaurishhs/keezle/utils"
	"github.com/go-playground/validator/v10"
)

type SessionConfig struct {
	// Expire is the duration after which the session will expire.
	// If zero, the session will not expire.
	// If negative, the session will expire immediately.
	// Default is 7 days.
	ActivePeriod time.Duration
	IdlePeriod   time.Duration
	Cookie       struct {
		Expires bool
		Name    string
		Secure  bool
	}
}

type AttributesConfig[UA, SA any] struct {
	User    UA
	Session SA
}

type Config[UA, SA models.AnyStruct] struct {
	Adapter                adapters.Adapter[UA, SA]
	Session                *SessionConfig
	Secret                 [][]byte
	Logger                 logger.Logger
	Hash                   func(string) (string, error)
	ComparePasswordAndHash func(string, string) (bool, error)
	Attributes             *AttributesConfig[UA, SA]
}

type Keezle[UA, SA models.AnyStruct] struct {
	Config    *Config[UA, SA]
	Validator *validator.Validate
}

func New[UA, SA models.AnyStruct](config *Config[UA, SA]) (res *Keezle[UA, SA]) {
	res = &Keezle[UA, SA]{
		Config:    config,
		Validator: validator.New(validator.WithRequiredStructEnabled()),
	}

	if res.Config.Logger == nil {
		res.Config.Logger = logger.StdoutLogger
	}

	if len(res.Config.Secret) == 0 {
		res.Config.Secret = [][]byte{[]byte(Default_Secret)}
		res.Config.Logger.Log("warn: no secret provided, using default secret")
	}

	if res.Config.Hash == nil {
		res.Config.Hash = utils.HashPassword
	}

	return nil
}
