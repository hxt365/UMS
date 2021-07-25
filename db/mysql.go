package db

import (
	"Shopee_UMS/utils"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func New(idleConn, maxConn int) (*sql.DB, error) {
	dbHost := utils.MustEnv("DATABASE_HOST")
	dbPort := utils.MustEnv("DATABASE_PORT")
	dbName := utils.MustEnv("DATABASE_NAME")
	user := utils.MustEnv("DATABASE_USER")
	pwd := utils.MustEnv("DATABASE_PASSWORD")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pwd, dbHost, dbPort, dbName)

	return newDB("mysql", dsn, idleConn, maxConn)
}

func NewTestDB() (*sql.DB, error) {
	dbHost := utils.MustEnv("DATABASE_HOST")
	dbPort := utils.MustEnv("DATABASE_PORT")
	dbName := utils.MustEnv("TEST_DATABASE_NAME")
	user := utils.MustEnv("DATABASE_USER")
	pwd := utils.MustEnv("DATABASE_PASSWORD")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pwd, dbHost, dbPort, dbName)

	return newDB("mysql", dsn, 10, 10)
}

func newDB(dialect, dsn string, idleConn, maxConn int) (*sql.DB, error) {
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
