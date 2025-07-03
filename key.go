package keezle

import "strings"

type CreateKeyOptions struct {
	UserID         string
	Provider       string
	ProviderUserID string
}

func createKeyId(provider, providerUserId string) (string, error) {
	if strings.Contains(provider, ":") {
		return "", ErrProviderColon
	}

	return provider + ":" + providerUserId, nil
}

func CreateKey()
