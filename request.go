package keezle

import (
	"errors"
	"net/http"
	"slices"
	"sync"

	"github.com/gaurishhs/keezle/models"
)

// AllowedMethods defines the HTTP methods that are allowed for CSRF protection.
var AllowedMethods = []string{
	"GET",
	"HEAD",
	"OPTIONS",
	"TRACE",
}

// AuthRequest represents an authentication request that can be validated and managed.
type AuthRequest[UA, SA models.AnyStruct] struct {
	Request      *http.Request
	SessionID    *string
	Keezle       *Keezle[UA, SA]
	validateOnce sync.Once
	validateRes  *models.Session[UA, SA]
	validateErr  error
}

// HandleRequest processes the incoming HTTP request and returns an AuthRequest.
// It checks for CSRF protection if configured and validates the request origin.
func (k *Keezle[UA, SA]) HandleRequest(req *http.Request) (*AuthRequest[UA, SA], error) {
	if k.Config.CSRF != nil {
		if !isValidRequestOrigin(k.Config.CSRF, req) {
			return nil, ErrInvalidRequestOrigin
		}
		return &AuthRequest[UA, SA]{
			Request:   req,
			Keezle:    k,
			SessionID: ptr(k.ReadSessionCookie(req)),
		}, nil
	}
	return &AuthRequest[UA, SA]{
		Request:   req,
		Keezle:    k,
		SessionID: nil,
	}, nil
}

// SetSession sets the session for the AuthRequest and updates the session cookie.
func (r *AuthRequest[UA, SA]) SetSession(session *models.Session[UA, SA]) {
	if session == nil {
		r.SessionID = nil
		r.Request.AddCookie(r.Keezle.CreateSessionCookie(nil))
		return
	}
	if deref(r.SessionID) == session.ID {
		return
	}
	r.SessionID = &session.ID
	r.Request.AddCookie(r.Keezle.CreateSessionCookie(session))
}

// Validate validates the session associated with the AuthRequest and resets the session if it is idle.
func (r *AuthRequest[UA, SA]) Validate() (*models.Session[UA, SA], error) {
	r.validateOnce.Do((func() {
		if r.SessionID == nil {
			return
		}
		session, err := r.Keezle.ValidateSession(deref(r.SessionID))
		if err != nil {
			if errors.Is(err, ErrInvalidSessionId) {
				r.SetSession(nil)
				return
			}
			r.validateErr = err
			return
		}
		if session.Fresh {
			r.SetSession(session)
		}
		r.validateRes = session
	}))
	return r.validateRes, r.validateErr
}

// Invalidate invalidates the internal cache for the validate method.
func (r *AuthRequest[UA, SA]) Invalidate() {
	r.validateOnce = sync.Once{}
	r.validateRes = nil
	r.validateErr = nil
}

func isValidRequestOrigin(config *CSRFProtectionConfig, req *http.Request) bool {
	if slices.Contains(AllowedMethods, req.Method) {
		return true
	}
	reqOrigin := req.Header.Get("Origin")
	if reqOrigin == "" {
		return false
	}
	if slices.Contains(config.AllowedSubdomains, "*") {
		if reqOrigin == config.Host {
			return true
		}
	}
	for _, subdomain := range config.AllowedSubdomains {
		if reqOrigin == subdomain+"."+config.Host {
			return true
		}
	}
	return false
}
