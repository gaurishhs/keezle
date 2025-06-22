package sqlite

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type TableConfig struct {
	SessionTable string
	UserTable    string
}

type SQLiteAdapter[T any] struct {
	DB     *sql.DB
	Tables TableConfig
}

func initalize[T any](dsnURI string) *SQLiteAdapter[T] {
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		panic("Failed to connect to SQLite database: " + err.Error())
	}
	return &SQLiteAdapter[T]{
		DB: db,
	}
}
