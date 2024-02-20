package helper

import (
	"database/sql"
	"fmt"
)

func NewDatabaseConnection(dsn string) (*sql.DB, error) {
	const fn = "helper.NewDatabaseConnection"

	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return db, nil
}
