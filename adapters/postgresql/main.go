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

func (a *PostgreSQLAdapter[UA, SA]) CreateUser(opts *adapters.CreateUserOpts[UA]) error {
	_, err := a.Conn.Exec(context.Background(), fmt.Sprintf("INSERT INTO `%s` (id, attributes) VALUES ($1, $2)", a.Tables.UserTable), opts.User.ID, opts.User.Attributes)
	if err != nil {
		return err
	}
	return nil
}

func (a *PostgreSQLAdapter[UA, SA]) GetUser(opts *adapters.CreateUserOpts[UA]) error {
	return nil
}

func (a *PostgreSQLAdapter[UA, SA]) UpdateUser(userId string, attributes UA) (*models.User[UA], error) {
	return nil, nil
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
