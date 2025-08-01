package keezle

import (
	"strings"

	"github.com/gaurishhs/keezle/models"
)

// CreateKeyOptions defines the options for creating a new key.
type CreateKeyOptions struct {
	UserID         string
	Provider       string
	ProviderUserID string
	Password       string
}

func createKeyId(provider, providerUserId string) (string, error) {
	if strings.Contains(provider, ":") {
		return "", ErrProviderColon
	}

	return provider + ":" + providerUserId, nil
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// TransformKey transforms a database key into a more secure representation.
// It dereferences the ID and UserID fields and sets the Password field based on whether it
// is nil or not.
func TransformKey(key *models.DBKey) *models.Key {
	return &models.Key{
		ID:       deref(key.ID),
		UserID:   deref(key.UserID),
		Password: key.Password != nil,
	}
}

// CreateKey creates a new key with the provided options.
func (k *Keezle[UA, SA]) CreateKey(opts CreateKeyOptions) (*models.Key, error) {
	keyId, err := createKeyId(opts.Provider, opts.ProviderUserID)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := k.Config.Hash(opts.Password)

	if err != nil {
		return nil, err
	}

	key := &models.DBKey{
		ID:       &keyId,
		UserID:   &opts.UserID,
		Password: &hashedPassword,
	}

	err = k.Config.Adapter.CreateKey(key)

	if err != nil {
		return nil, err
	}

	return TransformKey(key), nil
}

// DeleteKey deletes a key by its provider and provider user ID.
func (k *Keezle[UA, SA]) DeleteKey(provider, providerUserId string) error {
	keyId, err := createKeyId(provider, providerUserId)
	if err != nil {
		return err
	}
	return k.Config.Adapter.DeleteKey(keyId)
}

// GetKey retrieves a key by its provider and provider user ID.
func (k *Keezle[UA, SA]) GetKey(provider, providerUserId string) (*models.Key, error) {
	keyId, err := createKeyId(provider, providerUserId)
	if err != nil {
		return nil, err
	}

	key, err := k.Config.Adapter.GetKey(keyId)
	if err != nil {
		return nil, err
	}

	return TransformKey(key), nil
}

// GetKeysByUser retrieves all keys associated with a user by their user ID.
func (k *Keezle[UA, SA]) GetKeysByUser(userId string) ([]*models.Key, error) {
	_, err := k.GetUser(userId)
	if err != nil {
		return nil, err
	}

	dbKeys, err := k.Config.Adapter.GetKeysByUser(userId)
	if err != nil {
		return nil, err
	}

	keys := make([]*models.Key, len(dbKeys))
	for i, dbKey := range dbKeys {
		keys[i] = TransformKey(dbKey)
	}
	return keys, nil
}

// UpdateKey updates an existing key with a new password.
func (k *Keezle[UA, SA]) UpdateKey(provider, providerUserId, password string) (*models.Key, error) {
	keyId, err := createKeyId(provider, providerUserId)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := k.Config.Hash(password)
	if err != nil {
		return nil, err
	}

	updatedKey, err := k.Config.Adapter.UpdateKey(keyId, &models.DBKey{
		Password: &hashedPassword,
	})

	if err != nil {
		return nil, err
	}

	return TransformKey(updatedKey), nil
}

// UseKey retrieves a key by its provider and provider user ID, and validates the password if it exists.
func (k *Keezle[UA, SA]) UseKey(provider, providerUserId, password string) (*models.Key, error) {
	keyId, err := createKeyId(provider, providerUserId)
	if err != nil {
		return nil, err
	}
	key, err := k.Config.Adapter.GetKey(keyId)
	if err != nil {
		return nil, err
	}
	if key.Password != nil {
		if password == "" {
			return nil, ErrInvalidPassword
		}

		valid, err := k.Config.ComparePasswordAndHash(password, deref(key.Password))
		if err != nil {
			return nil, err
		}
		if !valid {
			return nil, ErrInvalidPassword
		}
	} else {
		if password != "" {
			return nil, ErrInvalidPassword
		}
	}
	return TransformKey(key), nil
}
