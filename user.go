package keezle

import (
	"github.com/gaurishhs/keezle/adapters"
	"github.com/gaurishhs/keezle/models"
	"github.com/rs/xid"
)

type CreateUserOptions[UA models.AnyStruct] struct {
	UserID string
	Key    struct {
		Provider       string
		ProviderUserID string
		Password       string
	}
	Attributes UA
}

func (k *Keezle[UA, SA]) CreateUser(opts CreateUserOptions[UA]) (*models.User[UA], error) {
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
	if opts.Key.Provider == "" && opts.Key.ProviderUserID == "" {
		k.Config.Adapter.CreateUser(&adapters.CreateUserOpts[UA]{
			User: user,
		})
	}

	keyId, err := createKeyId(opts.Key.Provider, opts.Key.ProviderUserID)
	if err != nil {
		return nil, err
	}

	err = k.Config.Adapter.CreateUser(&adapters.CreateUserOpts[UA]{
		User: user,
		Key: &models.DBKey{
			ID:       keyId,
			Password: opts.Key.Password,
		},
	})

	return user, err
}

func (k *Keezle[UA, SA]) GetUser(userId string) (*models.User[UA], error) {
	return k.Config.Adapter.GetUser(userId)
}

func (k *Keezle[UA, SA]) UpdateUser(userId string, attributes UA) (*models.User[UA], error) {
	return k.Config.Adapter.UpdateUser(userId, attributes)
}

func (k *Keezle[UA, SA]) DeleteUser(userId string) error {
	return k.Config.Adapter.DeleteUser(userId)
}
