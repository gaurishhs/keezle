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
	Attributes *UA
}

func (k *Keezle[UA, SA]) TransformUser(dbUser *models.User[UA]) (*models.User[UA], error) {
	userAttributes, err := k.Config.GetUserAttributes(dbUser)
	if err != nil {
		return nil, err
	}
	return &models.User[UA]{
		ID:         dbUser.ID,
		Attributes: userAttributes,
	}, nil
}

func (k *Keezle[UA, SA]) CreateUser(opts CreateUserOptions[UA]) (*models.User[UA], error) {
	if opts.UserID == "" {
		opts.UserID = xid.New().String()
	}
	user := &models.User[UA]{
		ID:         opts.UserID,
		Attributes: opts.Attributes,
	}
	if opts.Key.Provider == "" && opts.Key.ProviderUserID == "" {
		err := k.Config.Adapter.CreateUser(&adapters.CreateUserOpts[UA]{
			User: user,
		})
		if err != nil {
			return nil, err
		}
		return k.TransformUser(user)
	}

	keyId, err := createKeyId(opts.Key.Provider, opts.Key.ProviderUserID)
	if err != nil {
		return nil, err
	}

	err = k.Config.Adapter.CreateUser(&adapters.CreateUserOpts[UA]{
		User: user,
		Key: &models.DBKey{
			ID:       &keyId,
			Password: &opts.Key.Password,
		},
	})

	if err != nil {
		return nil, err
	}

	return k.TransformUser(user)
}

func (k *Keezle[UA, SA]) GetUser(userId string) (*models.User[UA], error) {
	user, err := k.Config.Adapter.GetUser(userId)
	if err != nil {
		return nil, err
	}
	return k.TransformUser(user)
}

func (k *Keezle[UA, SA]) UpdateUser(userId string, attributes UA) (*models.User[UA], error) {
	return k.Config.Adapter.UpdateUser(userId, attributes)
}

func (k *Keezle[UA, SA]) DeleteUser(userId string) error {
	return k.Config.Adapter.DeleteUser(userId)
}

func (k *Keezle[UA, SA]) GetUsersByAttribute(attribute string, value string) (*models.User[UA], error) {
	users, err := k.Config.Adapter.GetUsersByAttribute(attribute, value)
	if err != nil {
		return nil, err
	}
	var user *models.User[UA]
	for _, u := range users {
		user = u
		break
	}
	if user == nil {
		return nil, nil
	}
	return k.TransformUser(user)
}
