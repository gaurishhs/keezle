package oauth

type Provider interface {
	GetAuthorizationURL() (url, state string, err error)
	ValidateCallback(code, state string)
}
