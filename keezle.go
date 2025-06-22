package keezle

import (
	"time"

	"github.com/gaurishhs/keezle/adapters"
	"github.com/gaurishhs/keezle/logger"
	"github.com/gaurishhs/keezle/utils"
)

type SessionConfig struct {
	// Expire is the duration after which the session will expire.
	// If zero, the session will not expire.
	// If negative, the session will expire immediately.
	// Default is 7 days.
	Expires   time.Duration
	UpdateAge time.Duration
	FreshAge  time.Duration
}

type AttributesConfig struct {
	User    interface{}
	Session interface{}
}

type Config[UA, SA any] struct {
	Adapter    adapters.Adapter[UA, SA]
	Session    *SessionConfig
	Secret     [][]byte
	Logger     logger.Logger
	Hash       func(string) (string, error)
	Attributes *AttributesConfig
}

type Keezle[UA, SA any] struct {
	Config *Config[UA, SA]
}

func New[UA, SA any](config *Config[UA, SA]) (res *Keezle[UA, SA]) {
	res = &Keezle[UA, SA]{
		Config: config,
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
