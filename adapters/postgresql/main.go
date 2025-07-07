package postgresql

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gaurishhs/keezle/adapters"
	"github.com/gaurishhs/keezle/models"
	"github.com/jackc/pgx/v5"
)

type TableConfig struct {
	SessionTable string
	UserTable    string
	KeyTable     string
}

type PostgreSQLAdapter[UA, SA models.AnyStruct] struct {
	Conn   *pgx.Conn
	Tables TableConfig
}

func Initialize[UA, SA models.AnyStruct](connString string) *PostgreSQLAdapter[UA, SA] {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		panic("Failed to connect to PostgreSQL database: " + err.Error())
	}
	defer conn.Close(context.Background())
	return &PostgreSQLAdapter[UA, SA]{
		Conn: conn,
	}
}

func (a *PostgreSQLAdapter[UA, SA]) CreateUser(opts *adapters.CreateUserOpts[UA]) error {
	_, err := a.Conn.Exec(context.Background(), fmt.Sprintf("INSERT INTO `%s` (id, attributes) VALUES ($1, $2)", a.Tables.UserTable), opts.User.ID, opts.User.Attributes)
	if err != nil {
		return err
	}
	return nil
}

func (a *PostgreSQLAdapter[UA, SA]) GetUser(userId string) (*models.User[UA], error) {
	var (
		user       models.User[UA]
		attrsBytes []byte
	)
	row := a.Conn.QueryRow(context.Background(), fmt.Sprintf("SELECT `id`, `attributes` FROM `%s` WHERE `id` = $1", a.Tables.UserTable), userId)
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

func (a *PostgreSQLAdapter[UA, SA]) UpdateUser(userId string, attributes UA) (*models.User[UA], error) {
	updatedRow := a.Conn.QueryRow(context.Background(), fmt.Sprintf("UPDATE `%s` SET attributes = $1 where id = $2 returning id, attributes", a.Tables.UserTable), attributes, userId)
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

func (a *PostgreSQLAdapter[UA, SA]) DeleteUser(userId string) error {
	_, err := a.Conn.Exec(context.Background(), fmt.Sprintf("DELETE FROM `%s` WHERE `id` = $1", a.Tables.UserTable), userId)
	if err != nil {
		return err
	}
	return nil
}

func (a *PostgreSQLAdapter[UA, SA]) CreateSession(session *models.DBSession[SA]) error {
	_, err := a.Conn.Exec(context.Background(), fmt.Sprintf("INSERT INTO `%s` (`id`, `user_id`, `active_expires_at`, `idle_expires_at`, `attributes`) VALUES ($1, $2, $3, $4, $5)", a.Tables.SessionTable), session.ID, session.UserId, session.ActiveExpiresAt, session.IdleExpiresAt, session.Attributes)
	if err != nil {
		return err
	}
	return nil
}

func (a *PostgreSQLAdapter[UA, SA]) GetSessionAndUser(sessionId string) (*models.DBSession[SA], *models.User[UA], error) {
	row := a.Conn.QueryRow(context.Background(), fmt.Sprintf("SELECT `id`, `user_id`, `active_expires_at`, `idle_expires_at`, `attributes` FROM `%s` WHERE `id` = $1", a.Tables.SessionTable), sessionId)
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

func (a *PostgreSQLAdapter[UA, SA]) GetSessionsByUser(userId string) ([]*models.DBSession[SA], error) {
	rows, err := a.Conn.Query(
		context.Background(),
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
	sessions, err := pgx.CollectRows(rows, pgx.RowToStructByPos[*models.DBSession[SA]])
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (a *PostgreSQLAdapter[UA, SA]) UpdateSession(sessionId string, attributes SA) (*models.DBSession[SA], error) {
	updatedRow := a.Conn.QueryRow(context.Background(), fmt.Sprintf("UPDATE `%s` SET `attributes` = $1 where `id` = $2 returning `id`, `user_id`, `active_expires_at`, `idle_expires_at`, `attributes`", a.Tables.SessionTable), attributes, sessionId)
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

func (a *PostgreSQLAdapter[UA, SA]) DeleteSession(sessionId string) error {
	_, err := a.Conn.Exec(context.Background(), fmt.Sprintf("DELETE FROM `%s` WHERE `id` = $1", a.Tables.SessionTable), sessionId)
	if err != nil {
		return err
	}
	return nil
}

func (a *PostgreSQLAdapter[UA, SA]) DeleteAllUserSessions(userId string) error {
	_, err := a.Conn.Exec(context.Background(), fmt.Sprintf("DELETE FROM `%s` WHERE `user_id` = $1", a.Tables.SessionTable), userId)
	if err != nil {
		return err
	}
	return nil
}

func (a *PostgreSQLAdapter[UA, SA]) CreateKey(key *models.DBKey) error {
	_, err := a.Conn.Exec(context.Background(), fmt.Sprintf(
		"INSERT INTO `%s` (`id`, `user_id`, `password`) VALUES ($1, $2, $3)",
		a.Tables.KeyTable,
	), key.ID, key.UserID, key.Password)
	if err != nil {
		return err
	}
	return nil
}

func (a *PostgreSQLAdapter[UA, SA]) GetKey(keyId string) (*models.DBKey, error) {
	row := a.Conn.QueryRow(
		context.Background(),
		fmt.Sprintf("SELECT `id`, `user_id`, `password` FROM `%s` WHERE `id` = $1", a.Tables.KeyTable),
		keyId,
	)
	var key models.DBKey
	err := row.Scan(&key.ID, &key.UserID, &key.Password)
	if err != nil {
		return nil, err
	}
	return &key, err
}

func (a *PostgreSQLAdapter[UA, SA]) GetKeysByUser(userId string) ([]*models.DBKey, error) {
	rows, err := a.Conn.Query(
		context.Background(),
		fmt.Sprintf(
			"SELECT `id`, `user_id`, `attributes` FROM `%s` WHERE `user_id` = $1",
			a.Tables.KeyTable,
		),
		userId,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	keys, err := pgx.CollectRows(rows, pgx.RowToStructByPos[*models.DBKey])
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (a *PostgreSQLAdapter[UA, SA]) UpdateKey(keyId string, newKey *models.DBKey) (*models.DBKey, error) {
	updatedRow := a.Conn.QueryRow(
		context.Background(),
		fmt.Sprintf(
			"UPDATE `%s` SET "+
				"`id` = CASE WHEN $1 THEN $2 ELSE `id` END, "+
				"`user_id` = CASE WHEN $3 THEN $4 ELSE `user_id` END, "+
				"`password` = CASE WHEN $5 THEN $6 ELSE `password` END "+
				"WHERE `id` = $7 RETURNING `id`, `user_id`, `password`",
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

func (a *PostgreSQLAdapter[UA, SA]) DeleteKey(keyId string) error {
	_, err := a.Conn.Exec(context.Background(), fmt.Sprintf("DELETE FROM `%s` WHERE `id` = $1", a.Tables.KeyTable), keyId)
	if err != nil {
		return err
	}
	return nil
}
