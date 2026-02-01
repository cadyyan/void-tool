package configuration

import (
	"database/sql"
	"fmt"
)

type SQLiteConfiguration struct {
	Path string `required:"true"`
}

func (config SQLiteConfiguration) DSN() string {
	return fmt.Sprintf("file:%s?mode=rwc&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)", config.Path)
}

func (config SQLiteConfiguration) Connect() (*sql.DB, error) {
	conn, err := sql.Open("sqlite", config.DSN())
	if err != nil {
		return nil, fmt.Errorf("unable to connect to SQLite: %w", err)
	}

	return conn, nil
}
