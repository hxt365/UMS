package reposistory

import (
	"Shopee_UMS/usecases"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Accounts struct {
	db *sql.DB
	c  Cache
}

func NewAccounts(db *sql.DB, c Cache) *Accounts {
	return &Accounts{db: db, c: c}
}

func (a *Accounts) Get(username string) (*usecases.AccountData, error) {
	if acc, err := a.getCache(username); err == nil {
		return acc, nil
	}

	var acc usecases.AccountData
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := "SELECT id, username, password FROM accounts WHERE username = ?"
	if err := a.db.QueryRowContext(ctx, q, username).Scan(&acc.Id, &acc.Username, &acc.Password); err != nil {
		return nil, err
	}

	if err := a.cache(&acc); err != nil {
		return nil, err
	}
	return &acc, nil
}

func (a *Accounts) cache(acc *usecases.AccountData) error {
	key := fmt.Sprintf("accounts:%s", acc.Username)
	if err := a.c.StoreJson(key, acc); err != nil {
		return err
	}
	return nil
}

func (a *Accounts) getCache(username string) (*usecases.AccountData, error) {
	key := fmt.Sprintf("accounts:%s", username)
	var acc usecases.AccountData
	if err := a.c.LoadJson(key, &acc); err != nil {
		return nil, err
	}
	return &acc, nil
}
