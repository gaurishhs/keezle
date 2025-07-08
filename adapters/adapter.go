package adapters

import "github.com/gaurishhs/keezle/models"

type CreateUserOpts[UA models.AnyStruct] struct {
	User *models.User[UA]
	Key  *models.DBKey
}

type Adapter[UA, SA models.AnyStruct] interface {
	CreateUser(opts *CreateUserOpts[UA]) error
	GetUser(userId string) (*models.User[UA], error)
	GetUsersByAttribute(attribute string, value string) ([]*models.User[UA], error)
	UpdateUser(userId string, attributes UA) (*models.User[UA], error)
	DeleteUser(userId string) error
	CreateSession(session *models.DBSession[SA]) error
	GetSessionAndUser(sessionId string) (*models.DBSession[SA], *models.User[UA], error)
	GetSessionsByUser(userId string) ([]*models.DBSession[SA], error)
	UpdateSession(sessionId string, newSession *models.DBSession[SA]) (*models.DBSession[SA], error)
	DeleteSession(sessionId string) error
	DeleteAllUserSessions(userId string) error
	CreateKey(key *models.DBKey) error
	GetKey(keyId string) (*models.DBKey, error)
	GetKeysByUser(userId string) ([]*models.DBKey, error)
	UpdateKey(keyId string, updatedKey *models.DBKey) (*models.DBKey, error)
	DeleteKey(keyId string) error
}
