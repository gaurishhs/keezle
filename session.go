package keezle

import (
	"net/http"
	"time"

	"github.com/gaurishhs/keezle/models"
	"github.com/gaurishhs/keezle/utils"
)

type CreateSessionOptions[SA models.AnyStruct] struct {
	SessionId  string
	UserId     string
	Attributes SA
}

func TransformSession[UA, SA models.AnyStruct](dbSession *models.DBSession[SA], dbUser *models.User[UA]) (*models.Session[UA, SA], error) {
	user := &models.User[UA]{
		ID:         dbUser.ID,
		Attributes: dbUser.Attributes,
	}
	session := &models.Session[UA, SA]{
		ID:              dbSession.ID,
		User:            user,
		ActiveExpiresAt: dbSession.ActiveExpiresAt,
		IdleExpiresAt:   dbSession.IdleExpiresAt,
		Attributes:      dbSession.Attributes,
	}
	return session, nil
}

func (k *Keezle[UA, SA]) GetSession(sessionId string) (*models.Session[UA, SA], error) {
	if sessionId == "" {
		return nil, ErrInvalidSessionId
	}

	dbSession, user, err := k.Config.Adapter.GetSessionAndUser(sessionId)

	if err != nil {
		return nil, err
	}

	return TransformSession(dbSession, user)
}

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
		ID:              sessionId,
		UserId:          opts.UserId,
		Attributes:      opts.Attributes,
		ActiveExpiresAt: time.Now().Add(k.Config.Session.ActivePeriod),
		IdleExpiresAt:   time.Now().Add(k.Config.Session.ActivePeriod).Add(k.Config.Session.IdlePeriod),
	}
	user, err := k.GetUser(opts.UserId)
	if err != nil {
		return nil, err
	}

	err = k.Config.Adapter.CreateSession(session)
	if err != nil {
		return nil, err
	}

	return TransformSession(session, user)
}

func (k *Keezle[UA, SA]) UpdateSession(sessionId string, attributes SA) (*models.Session[UA, SA], error) {
	if sessionId == "" {
		return nil, ErrInvalidSessionId
	}

	dbSession, err := k.Config.Adapter.UpdateSession(sessionId, attributes)
	if err != nil {
		return nil, err
	}

	user, err := k.GetUser(dbSession.UserId)
	if err != nil {
		return nil, err
	}

	return TransformSession(dbSession, user)
}

func (k *Keezle[UA, SA]) DeleteSession(sessionId string) error {
	if sessionId == "" {
		return ErrInvalidSessionId
	}

	return k.Config.Adapter.DeleteSession(sessionId)
}

func (k *Keezle[UA, SA]) DeleteAllUserSessions(userId string) error {
	return k.Config.Adapter.DeleteAllUserSessions(userId)
}

func (k *Keezle[UA, SA]) DeleteInvalidUserSessions(userId string) error {
	dbSessions, err := k.Config.Adapter.GetSessionsByUser(userId)
	if err != nil {
		return err
	}

	for _, dbSession := range dbSessions {
		if dbSession.IdleExpiresAt.Before(time.Now()) {
			err = k.Config.Adapter.DeleteSession(dbSession.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

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
