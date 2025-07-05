package keezle

import (
	"strings"

	"github.com/gaurishhs/keezle/models"
)

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

func TransformKey(key *models.DBKey) *models.Key {
	return &models.Key{
		ID:       deref(key.ID),
		UserID:   deref(key.UserID),
		Password: key.Password != nil,
	}
}

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

func (k *Keezle[UA, SA]) DeleteKey(provider, providerUserId string) error {
	keyId, err := createKeyId(provider, providerUserId)
	if err != nil {
		return err
	}
	return k.Config.Adapter.DeleteKey(keyId)
}

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
