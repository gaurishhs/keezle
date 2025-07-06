package keezle

import (
	"time"

	"github.com/gaurishhs/keezle/adapters"
	"github.com/gaurishhs/keezle/logger"
	"github.com/gaurishhs/keezle/models"
	"github.com/gaurishhs/keezle/utils"
	"github.com/go-playground/validator/v10"
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

type AttributesConfig[UA, SA any] struct {
	User    UA
	Session SA
}

type Config[UA, SA models.AnyStruct] struct {
	Adapter                adapters.Adapter[UA, SA]
	Session                *SessionConfig
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

	if res.Config.Adapter == nil {
		panic("adapter is required")
	}

	if res.Config.Logger == nil {
		res.Config.Logger = logger.NoOpLogger
	}

	if res.Config.Hash == nil {
		res.Config.Hash = utils.HashPassword
	}

	if res.Config.ComparePasswordAndHash == nil {
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
