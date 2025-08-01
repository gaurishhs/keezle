package keezle

import (
	"net/http"
	"time"

	"github.com/gaurishhs/keezle/models"
	"github.com/gaurishhs/keezle/utils"
)

// CreateSessionOptions defines the options for creating a new session.
type CreateSessionOptions[SA models.AnyStruct] struct {
	SessionId  string
	UserId     string
	Attributes SA
}

func isValidSession[SA models.AnyStruct](dbSession *models.DBSession[SA]) bool {
	return dbSession.IdleExpiresAt.After(time.Now())
}

func derefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func ptr[T any](v T) *T {
	return &v
}

// TransformSession transforms a database session into a session with attributes.
// It uses the GetSessionAttributes function from the configuration to extract session attributes.
func (k *Keezle[UA, SA]) TransformSession(dbSession *models.DBSession[SA], dbUser *models.User[UA], fresh bool) (*models.Session[UA, SA], error) {
	sessionAttributes, err := k.Config.GetSessionAttributes(dbSession)
	if err != nil {
		return nil, err
	}
	var state string = "idle"
	if dbSession.ActiveExpiresAt.Before(time.Now()) {
		state = "active"
	}
	session := &models.Session[UA, SA]{
		ID:              deref(dbSession.ID),
		User:            dbUser,
		ActiveExpiresAt: derefTime(dbSession.ActiveExpiresAt),
		IdleExpiresAt:   derefTime(dbSession.IdleExpiresAt),
		Attributes:      sessionAttributes,
		State:           state,
	}
	return session, nil
}

// GetSession retrieves a session by its ID.
func (k *Keezle[UA, SA]) GetSession(sessionId string) (*models.Session[UA, SA], error) {
	if sessionId == "" {
		return nil, ErrInvalidSessionId
	}

	dbSession, dbUser, err := k.Config.Adapter.GetSessionAndUser(sessionId)

	if err != nil {
		return nil, err
	}

	user, err := k.TransformUser(dbUser)
	if err != nil {
		return nil, err
	}

	return k.TransformSession(dbSession, user, false)
}

// GetAllUserSessions retrieves all sessions for a user by their user ID.
func (k *Keezle[UA, SA]) GetAllUserSessions(userId string) ([]*models.Session[UA, SA], error) {
	dbSessions, err := k.Config.Adapter.GetSessionsByUser(userId)
	if err != nil {
		return nil, err
	}

	var sessions []*models.Session[UA, SA]
	for _, dbSession := range dbSessions {
		if !isValidSession(dbSession) {
			continue
		}
		user, err := k.GetUser(deref(dbSession.UserId))
		if err != nil {
			return nil, err
		}
		session, err := k.TransformSession(dbSession, user, false)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

// CreateSession creates a new session with the provided options.
func (k *Keezle[UA, SA]) CreateSession(opts CreateSessionOptions[SA]) (*models.Session[UA, SA], error) {
	sessionId := opts.SessionId
	if sessionId == "" {
		id, err := utils.GenerateRandomString(32)
		if err != nil {
			return nil, err
		}
		sessionId = id
	}
	session := &models.DBSession[SA]{
		ID:              &sessionId,
		UserId:          &opts.UserId,
		Attributes:      &opts.Attributes,
		ActiveExpiresAt: ptr(time.Now().Add(k.Config.Session.ActivePeriod)),
		IdleExpiresAt:   ptr(time.Now().Add(k.Config.Session.ActivePeriod).Add(k.Config.Session.IdlePeriod)),
	}
	user, err := k.GetUser(opts.UserId)
	if err != nil {
		return nil, err
	}

	err = k.Config.Adapter.CreateSession(session)
	if err != nil {
		return nil, err
	}

	return k.TransformSession(session, user, false)
}

// UpdateSession updates an existing session with new attributes.
func (k *Keezle[UA, SA]) UpdateSession(sessionId string, newSession *models.DBSession[SA]) (*models.Session[UA, SA], error) {
	if sessionId == "" {
		return nil, ErrInvalidSessionId
	}

	dbSession, err := k.Config.Adapter.UpdateSession(sessionId, newSession)
	if err != nil {
		return nil, err
	}

	user, err := k.GetUser(deref(newSession.UserId))
	if err != nil {
		return nil, err
	}

	return k.TransformSession(dbSession, user, false)
}

// DeleteSession deletes a session by its ID.
func (k *Keezle[UA, SA]) DeleteSession(sessionId string) error {
	if sessionId == "" {
		return ErrInvalidSessionId
	}

	return k.Config.Adapter.DeleteSession(sessionId)
}

// DeleteAllUserSessions deletes all sessions for a user by their user ID.
func (k *Keezle[UA, SA]) DeleteAllUserSessions(userId string) error {
	return k.Config.Adapter.DeleteAllUserSessions(userId)
}

// DeleteInvalidUserSessions deletes all invalid sessions for a user by their user ID.
func (k *Keezle[UA, SA]) DeleteInvalidUserSessions(userId string) error {
	dbSessions, err := k.Config.Adapter.GetSessionsByUser(userId)
	if err != nil {
		return err
	}

	for _, dbSession := range dbSessions {
		if isValidSession(dbSession) {
			continue
		}
		err = k.Config.Adapter.DeleteSession(deref(dbSession.ID))
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateSession checks if a session is valid and returns the session if it is active.
// If the session is idle, it updates the session's expiration times and returns the updated session.
func (k *Keezle[UA, SA]) ValidateSession(sessionId string) (*models.Session[UA, SA], error) {
	if sessionId == "" {
		return nil, ErrInvalidSessionId
	}

	dbSession, dbUser, err := k.Config.Adapter.GetSessionAndUser(sessionId)
	if err != nil {
		return nil, err
	}

	user, err := k.TransformUser(dbUser)
	if err != nil {
		return nil, err
	}

	session, err := k.TransformSession(dbSession, user, false)
	if err != nil {
		return nil, err
	}

	if session.State == "active" {
		return session, nil
	}

	updatedSession, err := k.UpdateSession(sessionId, &models.DBSession[SA]{
		ActiveExpiresAt: ptr(time.Now().Add(k.Config.Session.ActivePeriod)),
		IdleExpiresAt:   ptr(time.Now().Add(k.Config.Session.ActivePeriod).Add(k.Config.Session.IdlePeriod)),
	})

	if err != nil {
		return nil, err
	}

	return &models.Session[UA, SA]{
		ID:              updatedSession.ID,
		User:            dbUser,
		ActiveExpiresAt: updatedSession.ActiveExpiresAt,
		IdleExpiresAt:   updatedSession.IdleExpiresAt,
		Attributes:      updatedSession.Attributes,
		State:           updatedSession.State,
		Fresh:           true,
	}, nil
}

// CreateSessionCookie creates a http cookie for the session.
func (k *Keezle[UA, SA]) CreateSessionCookie(session *models.Session[UA, SA]) *http.Cookie {
	var expires time.Time
	if session == nil {
		expires = time.Unix(0, 0)
	} else if k.Config.Session.Cookie.Expires {
		expires = session.IdleExpiresAt
	} else {
		expires = time.Now().Add(time.Hour * 24 * 365)
	}

	return &http.Cookie{
		Name:     k.Config.Session.Cookie.Name,
		HttpOnly: true,
		Secure:   k.Config.Session.Cookie.Secure,
		Expires:  expires,
		Value:    session.ID,
		// TODO: k.Config.Session.Cookie attributes should be configurable
	}
}

// ReadSessionCookie reads the session cookie from the request and returns its value.
// If the cookie is not found, it returns an empty string.
func (k *Keezle[UA, SA]) ReadSessionCookie(req *http.Request) string {
	cookie, err := req.Cookie(k.Config.Session.Cookie.Name)
	if err != nil {
		return ""
	}
	return cookie.Value
}
