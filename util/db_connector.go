package util

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func GetDBConnection(addr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", addr)
	if err != nil {
		return nil, fmt.Errorf("db failed to open, %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("db failed to ping, %w", err)
	}

	return db, nil
}
