package reposistory

import (
	"Shopee_UMS/usecases"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Users struct {
	db *sql.DB
	c  Cache
}

func NewUsers(db *sql.DB, c Cache) *Users {
	return &Users{db: db, c: c}
}

func (u *Users) Get(uid int) (*usecases.UserData, error) {
	if ud, err := u.getCache(uid); err == nil {
		return ud, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var ud usecases.UserData
	q := `SELECT a.id, a.username, u.nickname, u.profile_picture_uri
	FROM users u
	JOIN accounts a ON u.account_id = a.id
	WHERE u.id = ?`
	err := u.db.QueryRowContext(ctx, q, uid).Scan(&ud.Id, &ud.Username, &ud.Nickname, &ud.ProfilePicUri)
	if err != nil {
		return nil, err
	}

	if err := u.cache(&ud); err != nil {
		return nil, err
	}
	return &ud, nil
}

func (u *Users) UpdateNickname(uid int, nickname string) error {
	u.deleteCache(uid)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := `UPDATE users SET nickname = ? WHERE id = ?`
	stmt, err := u.db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, nickname, uid)
	return err
}

func (u *Users) UpdateProfilePicUri(uid int, uri string) error {
	u.deleteCache(uid)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := `UPDATE users SET profile_picture_uri = ? WHERE id = ?`
	stmt, err := u.db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, uri, uid)
	return err
}

func (u *Users) cache(ud *usecases.UserData) error {
	key := fmt.Sprintf("users:%d", ud.Id)
	if err := u.c.StoreJson(key, ud); err != nil {
		return err
	}
	return nil
}

func (u *Users) getCache(uid int) (*usecases.UserData, error) {
	key := fmt.Sprintf("users:%d", uid)
	var ud usecases.UserData
	if err := u.c.LoadJson(key, &ud); err != nil {
		return nil, err
	}
	return &ud, nil
}

func (u *Users) deleteCache(uid int) {
	key := fmt.Sprintf("users:%d", uid)
	_ = u.c.Delete(key) // no need to handle error because cache failure should not affect the dataflow
}
