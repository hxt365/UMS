package reposistory

import (
	"Shopee_UMS/usecases"
	"context"
	"database/sql"
	"time"
)

type Accounts struct {
	db *sql.DB
}

func NewAccounts(db *sql.DB) *Accounts {
	return &Accounts{db: db}
}

func (a *Accounts) Get(username string) (*usecases.AccountData, error) {
	var acc usecases.AccountData
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := "SELECT id, username, password FROM accounts WHERE username = ?"
	err := a.db.QueryRowContext(ctx, q, username).Scan(&acc.Id, &acc.Username, &acc.Password)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}