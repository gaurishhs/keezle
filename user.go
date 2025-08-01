package keezle

import (
	"github.com/gaurishhs/keezle/adapters"
	"github.com/gaurishhs/keezle/models"
	"github.com/rs/xid"
)

// CreateUserOptions defines the options for creating a new user.
type CreateUserOptions[UA models.AnyStruct] struct {
	UserID string
	Key    struct {
		Provider       string
		ProviderUserID string
		Password       string
	}
	Attributes *UA
}

// TransformUser transforms a database user into a user with attributes.
// It uses the GetUserAttributes function from the configuration to extract user attributes.
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

// CreateUser creates a new user with the provided options.
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

// GetUser retrieves a user by their ID.
func (k *Keezle[UA, SA]) GetUser(userId string) (*models.User[UA], error) {
	user, err := k.Config.Adapter.GetUser(userId)
	if err != nil {
		return nil, err
	}
	return k.TransformUser(user)
}

// UpdateUser updates the attributes of an existing user.
func (k *Keezle[UA, SA]) UpdateUser(userId string, attributes UA) (*models.User[UA], error) {
	return k.Config.Adapter.UpdateUser(userId, attributes)
}

// DeleteUser deletes a user by their ID.
func (k *Keezle[UA, SA]) DeleteUser(userId string) error {
	return k.Config.Adapter.DeleteUser(userId)
}

// GetUsersByAttribute retrieves users based on a specific attribute and its value.
func (k *Keezle[UA, SA]) GetUsersByAttribute(attribute string, value string) ([]*models.User[UA], error) {
	users, err := k.Config.Adapter.GetUsersByAttribute(attribute, value)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	var transformedUsers []*models.User[UA]
	for _, user := range users {
		transformedUser, err := k.TransformUser(user)
		if err != nil {
			return nil, err
		}
		transformedUsers = append(transformedUsers, transformedUser)
	}
	return transformedUsers, nil
}
