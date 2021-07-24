package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func New(dialect, dsn string, idleConn, maxConn int) (*sql.DB, error) {
	db, err := sql.Open(dialect, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(idleConn)
	db.SetMaxOpenConns(maxConn)

	return db, nil
}
