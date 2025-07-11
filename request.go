package keezle

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gaurishhs/keezle/models"
)

type AuthRequest[UA, SA models.AnyStruct] struct {
	Request      *http.Request
	SessionID    *string
	Keezle       *Keezle[UA, SA]
	validateOnce sync.Once
	validateRes  *models.Session[UA, SA]
	validateErr  error
}

// TODO: Add CSRF Protection

func (k *Keezle[UA, SA]) HandleRequest(req *http.Request) (*AuthRequest[UA, SA], error) {
	return &AuthRequest[UA, SA]{
		Request: req,
		Keezle:  k,
	}, nil
}

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

func (r *AuthRequest[UA, SA]) Validate() (*models.Session[UA, SA], error) {
	r.validateOnce.Do((func() {
		if r.SessionID == nil {
			return
		}
		session, err := r.Keezle.ValidateSession(*r.SessionID)
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

func (r *AuthRequest[UA, SA]) Invalidate() {
	r.validateOnce = sync.Once{}
	r.validateRes = nil
	r.validateErr = nil
}
