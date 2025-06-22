package keezle

import (
	"github.com/gaurishhs/keezle/models"
	"github.com/go-playground/validator/v10"
	"github.com/rs/xid"
)

type functions[UA, SA any] struct {
	config    *Config[UA, SA]
	validator *validator.Validate
}

type CreateUserOptions[UA any] struct {
	UserID     string
	Password   string
	Attributes UA
}

// CreateUser creates a new user with the given options.
func (f *functions[UA, SA]) CreateUser(opts *CreateUserOptions[UA]) (*models.User[UA], error) {
	if opts.UserID == "" {
		opts.UserID = xid.New().String()
	}
	// TODO: validate user attributes if needed
	// Set user in db
	// Hash password
	hashedPassword, err := f.config.Hash(opts.Password)
	if err != nil {
		return nil, err
	}
}
