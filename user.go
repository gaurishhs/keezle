package keezle

import (
	"github.com/gaurishhs/keezle/adapters"
	"github.com/gaurishhs/keezle/models"
	"github.com/rs/xid"
)

type CreateUserOptions[UA models.AnyStruct] struct {
	UserID  string
	Account struct {
		Provider       string
		ProviderUserID string
		Password       string
	}
	Attributes UA
}

// CreateUser creates a new user with the given options.
func (k *Keezle[UA, SA]) CreateUser(opts *CreateUserOptions[UA]) (*models.User[UA], error) {
	if opts.UserID == "" {
		opts.UserID = xid.New().String()
	}
	// if err := k.Validator.Struct(opts.Attributes); err != nil {
	// 	return nil, err
	// }
	user := &models.User[UA]{
		ID:         opts.UserID,
		Attributes: opts.Attributes,
	}
	if opts.Account.Provider == "" && opts.Account.ProviderUserID == "" {
		k.Config.Adapter.CreateUser(&adapters.CreateUserOpts[UA]{
			User: user,
		})
	}

	keyId, err := createKeyId(opts.Account.Provider, opts.Account.ProviderUserID)
	if err != nil {
		return nil, err
	}

	err = k.Config.Adapter.CreateUser(&adapters.CreateUserOpts[UA]{
		User: user,
		Key: &models.Key{
			ID:       keyId,
			Password: opts.Account.Password,
		},
	})

	return user, err
}

// GetUser returns the user associated the given userId.
func (k *Keezle[UA, SA]) GetUser(userId string) (*models.User[UA], error) {
	return k.Config.Adapter.GetUser(userId)
}

// UpdateUser updates the attributes of the user associated with the given userId.
func (k *Keezle[UA, SA]) UpdateUser(userId string, attributes UA) (*models.User[UA], error) {
	return k.Config.Adapter.UpdateUser(userId, attributes)
}

// DeleteUser deletes the user associated with the given userId.
func (k *Keezle[UA, SA]) DeleteUser(userId string) error {
	return nil
}
