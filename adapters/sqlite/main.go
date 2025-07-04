package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/gaurishhs/keezle/models"

	_ "modernc.org/sqlite"
)

type TableConfig struct {
	SessionTable string
	UserTable    string
}

type SQLiteAdapter[UA, SA models.AnyStruct] struct {
	DB     *sql.DB
	Tables TableConfig
}

type CreateUserOpts[UA models.AnyStruct] struct {
	User *models.User[UA]
	Key  *models.DBKey
}

func initalize[UA, SA models.AnyStruct](dsnURI string) *SQLiteAdapter[UA, SA] {
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		panic("Failed to connect to SQLite database: " + err.Error())
	}
	return &SQLiteAdapter[UA, SA]{
		DB: db,
	}
}

func (a *SQLiteAdapter[UA, SA]) CreateUser(opts CreateUserOpts[UA]) error {
	_, err := a.DB.Exec(fmt.Sprintf("insert into `%s` (id, attributes) values (?, ?)", a.Tables.UserTable), opts.User.ID, opts.User.Attributes)
	if err != nil {
		return err
	}
	return nil
}

func (a *SQLiteAdapter[UA, SA]) GetUser(userId string) (*models.User[UA], error) {
	var user models.User[UA]
	row := a.DB.QueryRow(fmt.Sprintf("select id, attributes from `%s` where id = ?", a.Tables.UserTable), userId)
	err := row.Scan(&user.ID, &user.Attributes)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *SQLiteAdapter[UA, SA]) UpdateUser(userId string, attributes UA) (*models.User[UA], error) {
	updatedRow := a.DB.QueryRow(fmt.Sprintf("update `%s` set attributes = ? where id = ? returning id, attributes", a.Tables.UserTable), attributes, userId)
	var updatedUser models.User[UA]
	if err := updatedRow.Scan(&updatedUser.ID, updatedUser.Attributes); err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

func (a *SQLiteAdapter[UA, SA]) DeleteUser(userId string) error {
	_, err := a.DB.Exec(fmt.Sprintf("delete from `%s` where id = ?", a.Tables.UserTable), userId)
	if err != nil {
		return err
	}
	return nil
}

func (a *SQLiteAdapter[UA, SA]) CreateSession(session *models.DBSession[SA]) error {
	_, err := a.DB.Exec(fmt.Sprintf("insert into `%s` (id, userId, activeExpiresAt, idleExpiresAt, attributes) values (?, ?, ?, ?, ?) ", a.Tables.SessionTable), session.ID, session.UserId, session.ActiveExpiresAt, session.IdleExpiresAt, session.Attributes)

	if err != nil {
		return err
	}

	return nil
}

func (a *SQLiteAdapter[UA, SA]) UpdateSession(sessionId string, attributes SA) (*models.DBSession[SA], error) {
	updatedRow := a.DB.QueryRow(fmt.Sprintf("update `%s` set attributes = ? where id = ? returning id, userId, activeExpiresAt, idleExpiresAt, attributes", a.Tables.SessionTable), attributes, sessionId)
	var session models.DBSession[SA]
	err := updatedRow.Scan(session.ID, session.UserId, session.ActiveExpiresAt, session.IdleExpiresAt, session.Attributes)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (a *SQLiteAdapter[UA, SA]) DeleteSession(sessionId string) error {
	_, err := a.DB.Exec(fmt.Sprintf("delete from `%s` where id = ?", a.Tables.SessionTable), sessionId)
	if err != nil {
		return err
	}
	return nil
}
