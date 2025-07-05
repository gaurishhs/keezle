package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/gaurishhs/keezle/models"

	_ "modernc.org/sqlite"
)

type TableConfig struct {
	SessionTable string
	UserTable    string
	KeyTable     string
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
	var (
		user       models.User[UA]
		attrsBytes []byte
	)
	row := a.DB.QueryRow(fmt.Sprintf("select id, attributes from `%s` where id = ?", a.Tables.UserTable), userId)
	err := row.Scan(&user.ID, &attrsBytes)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(attrsBytes, &user.Attributes)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *SQLiteAdapter[UA, SA]) UpdateUser(userId string, attributes UA) (*models.User[UA], error) {
	updatedRow := a.DB.QueryRow(fmt.Sprintf("UPDATE `%s` SET attributes = ? where id = ? returning id, attributes", a.Tables.UserTable), attributes, userId)
	var (
		updatedUser models.User[UA]
		attrsBytes  []byte
	)
	if err := updatedRow.Scan(&updatedUser.ID, &attrsBytes); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(attrsBytes, &updatedUser.Attributes); err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (a *SQLiteAdapter[UA, SA]) DeleteUser(userId string) error {
	_, err := a.DB.Exec(fmt.Sprintf("DELETE FROM `%s` where `id` = ?", a.Tables.UserTable), userId)
	if err != nil {
		return err
	}
	return nil
}

func (a *SQLiteAdapter[UA, SA]) CreateSession(session *models.DBSession[SA]) error {
	_, err := a.DB.Exec(fmt.Sprintf("INSERT INTO `%s` (`id`, `user_id`, `active_expires_at`, `idle_expires_at`, `attributes`) VALUES (?, ?, ?, ?, ?) ", a.Tables.SessionTable), session.ID, session.UserId, session.ActiveExpiresAt, session.IdleExpiresAt, session.Attributes)

	if err != nil {
		return err
	}

	return nil
}

func (a *SQLiteAdapter[UA, SA]) GetSessionAndUser(sessionId string) (*models.DBSession[SA], *models.User[UA], error) {
	row := a.DB.QueryRow(fmt.Sprintf("SELECT `id`, `user_id`, `active_expires_at`, `idle_expires_at`, `attributes` FROM `%s` WHERE `id` = ?", a.Tables.SessionTable), sessionId)
	var (
		session    models.DBSession[SA]
		attrsBytes []byte
	)
	if err := row.Scan(&session.ID, &session.UserId, &session.ActiveExpiresAt, &session.IdleExpiresAt, &attrsBytes); err != nil {
		return nil, nil, err
	}

	if err := json.Unmarshal(attrsBytes, &session.Attributes); err != nil {
		return nil, nil, err
	}

	user, err := a.GetUser(session.UserId)
	if err != nil {
		return nil, nil, err
	}

	return &session, user, nil
}

func (a *SQLiteAdapter[UA, SA]) GetSessionsByUser(userId string) ([]*models.DBSession[SA], error) {
	rows, err := a.DB.Query(
		fmt.Sprintf(
			"SELECT `id`, `user_id`, `active_expires_at`, `idle_expires_at`, `attributes` FROM `%s` WHERE `user_id` = ?",
			a.Tables.SessionTable,
		),
		userId,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var sessions []*models.DBSession[SA]
	err = rowsToStructs(rows, &sessions)
	return sessions, err
}

func (a *SQLiteAdapter[UA, SA]) UpdateSession(sessionId string, attributes SA) (*models.DBSession[SA], error) {
	updatedRow := a.DB.QueryRow(fmt.Sprintf("UPDATE `%s` SET `attributes` = ? where `id` = ? returning `id`, `user_id`, `active_expires_at`, `idle_expires_at`, `attributes`", a.Tables.SessionTable), attributes, sessionId)
	var (
		session    models.DBSession[SA]
		attrsBytes []byte
	)

	if err := updatedRow.Scan(session.ID, session.UserId, session.ActiveExpiresAt, session.IdleExpiresAt, &attrsBytes); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(attrsBytes, &session.Attributes); err != nil {
		return nil, err
	}

	return &session, nil
}

func (a *SQLiteAdapter[UA, SA]) DeleteSession(sessionId string) error {
	_, err := a.DB.Exec(fmt.Sprintf("DELETE FROM `%s` WHERE `id` = ?", a.Tables.SessionTable), sessionId)
	if err != nil {
		return err
	}
	return nil
}

func (a *SQLiteAdapter[UA, SA]) DeleteAllUserSessions(userId string) error {
	_, err := a.DB.Exec(fmt.Sprintf("DELETE FROM `%s` WHERE `user_id` = ?", a.Tables.SessionTable), userId)
	if err != nil {
		return err
	}
	return nil
}

func (a *SQLiteAdapter[UA, SA]) CreateKey(key *models.DBKey) error {
	_, err := a.DB.Exec(fmt.Sprintf(
		"INSERT INTO `%s` (`id`, `user_id`, `password`) VALUES (?, ?, ?)",
		a.Tables.KeyTable,
	), key.ID, key.UserID, key.Password)
	if err != nil {
		return err
	}
	return nil
}

func (a *SQLiteAdapter[UA, SA]) GetKey(keyId string) (*models.DBKey, error) {
	row := a.DB.QueryRow(
		fmt.Sprintf("SELECT `id`, `user_id`, `password` FROM `%s` WHERE `id` = ?", a.Tables.KeyTable),
		keyId,
	)
	var key models.DBKey
	err := row.Scan(&key.ID, &key.UserID, &key.Password)
	if err != nil {
		return nil, err
	}
	return &key, err
}

func (a *SQLiteAdapter[UA, SA]) GetKeysByUser(userId string) ([]*models.DBKey, error) {
	rows, err := a.DB.Query(
		fmt.Sprintf(
			"SELECT `id`, `user_id`, `attributes` FROM `%s` WHERE `user_id` = ?",
			a.Tables.KeyTable,
		),
		userId,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var keys []*models.DBKey
	err = rowsToStructs(rows, &keys)
	return keys, err
}

func (a *SQLiteAdapter[UA, SA]) UpdateKey(keyId string, newKey *models.DBKey) (*models.DBKey, error) {
	updatedRow := a.DB.QueryRow(
		fmt.Sprintf(
			"UPDATE `%s` SET "+
				"`id` = CASE WHEN ? THEN ? ELSE `id` END, "+
				"`user_id` = CASE WHEN ? THEN ? ELSE `user_id` END, "+
				"`password` = CASE WHEN ? THEN ? ELSE `password` END "+
				"WHERE `id` = ? RETURNING `id`, `user_id`, `password`",
			a.Tables.KeyTable,
		),
		newKey.ID != nil,
		keyId,
		newKey.UserID != nil,
		newKey.UserID,
		newKey.Password != nil,
		newKey.Password,
		keyId,
	)
	var updatedKey models.DBKey
	err := updatedRow.Scan(&updatedKey.ID, &updatedKey.UserID, &updatedKey.Password)
	if err != nil {
		return nil, err
	}
	return &updatedKey, nil
}

func (a *SQLiteAdapter[UA, SA]) DeleteKey(keyId string) error {
	_, err := a.DB.Exec(fmt.Sprintf("DELETE FROM `%s` WHERE `id` = ?", a.Tables.KeyTable), keyId)
	if err != nil {
		return err
	}
	return nil
}
