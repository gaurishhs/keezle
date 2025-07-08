package postgresql

import (
	"context"
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

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (a *PostgreSQLAdapter[UA, SA]) CreateUser(opts *adapters.CreateUserOpts[UA]) error {
	_, err := a.Conn.Exec(context.Background(), fmt.Sprintf("INSERT INTO `%s` (id, attributes) VALUES ($1, $2)", a.Tables.UserTable), opts.User.ID, opts.User.Attributes)
	if err != nil {
		return err
	}
	return nil
}

func (a *PostgreSQLAdapter[UA, SA]) GetUser(userId string) (*models.User[UA], error) {
	var user models.User[UA]
	row := a.Conn.QueryRow(context.Background(), fmt.Sprintf("SELECT `id`, `attributes` FROM `%s` WHERE `id` = $1", a.Tables.UserTable), userId)
	if err := row.Scan(&user.ID, &user.Attributes); err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *PostgreSQLAdapter[UA, SA]) UpdateUser(userId string, attributes UA) (*models.User[UA], error) {
	updatedRow := a.Conn.QueryRow(context.Background(), fmt.Sprintf("UPDATE `%s` SET attributes = $1 where id = $2 returning id, attributes", a.Tables.UserTable), attributes, userId)
	var updatedUser models.User[UA]
	if err := updatedRow.Scan(&updatedUser.ID, &updatedUser.Attributes); err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (a *PostgreSQLAdapter[UA, SA]) DeleteUser(userId string) error {
	if _, err := a.Conn.Exec(context.Background(), fmt.Sprintf("DELETE FROM `%s` WHERE `id` = $1", a.Tables.UserTable), userId); err != nil {
		return err
	}
	return nil
}

func (a *PostgreSQLAdapter[UA, SA]) CreateSession(session *models.DBSession[SA]) error {
	if _, err := a.Conn.Exec(context.Background(), fmt.Sprintf("INSERT INTO `%s` (`id`, `user_id`, `active_expires_at`, `idle_expires_at`, `attributes`) VALUES ($1, $2, $3, $4, $5)", a.Tables.SessionTable), session.ID, session.UserId, session.ActiveExpiresAt, session.IdleExpiresAt, session.Attributes); err != nil {
		return err
	}
	return nil
}

func (a *PostgreSQLAdapter[UA, SA]) GetSessionAndUser(sessionId string) (*models.DBSession[SA], *models.User[UA], error) {
	row := a.Conn.QueryRow(context.Background(), fmt.Sprintf("SELECT `id`, `user_id`, `active_expires_at`, `idle_expires_at`, `attributes` FROM `%s` WHERE `id` = $1", a.Tables.SessionTable), sessionId)
	var session models.DBSession[SA]
	if err := row.Scan(&session.ID, &session.UserId, &session.ActiveExpiresAt, &session.IdleExpiresAt, &session.Attributes); err != nil {
		return nil, nil, err
	}

	user, err := a.GetUser(deref(session.UserId))
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

func (a *PostgreSQLAdapter[UA, SA]) UpdateSession(sessionId string, newSession *models.DBSession[SA]) (*models.DBSession[SA], error) {
	updatedRow := a.Conn.QueryRow(
		context.Background(),
		fmt.Sprintf(
			"UPDATE `%s` SET "+
				"`id` = COALESCE($1, `id`), "+
				"`user_id` = COALESCE($2, `user_id`), "+
				"`active_expires_at` = COALESCE($3, `active_expires_at`), "+
				"`idle_expires_at` = COALESCE($4, `idle_expires_at`), "+
				"`attributes` = COALESCE($5, `attributes`) "+
				"WHERE `id` = $6 RETURNING `id`, `user_id`, `active_expires_at`, `idle_expires_at`, `attributes`",
			a.Tables.SessionTable,
		),
		newSession.ID,
		newSession.UserId,
		newSession.ActiveExpiresAt,
		newSession.IdleExpiresAt,
		newSession.Attributes,
		sessionId,
	)
	var session models.DBSession[SA]
	if err := updatedRow.Scan(&session.ID, &session.UserId, &session.ActiveExpiresAt, &session.IdleExpiresAt, &session.Attributes); err != nil {
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
				"`id` = COALESCE($1, `id`), "+
				"`user_id` = COALESCE($2, `user_id`), "+
				"`password` = COALESCE($3, `password`) "+
				"WHERE `id` = $4 RETURNING `id`, `user_id`, `password`",
			a.Tables.KeyTable,
		),
		newKey.ID,
		newKey.UserID,
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
