package reposistory

import (
	"Shopee_UMS/usecases"
	"context"
	"database/sql"
	"time"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{db: db}
}

func (u *Users) Get(uid int) (*usecases.UserData, error) {
	var ud usecases.UserData

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := `SELECT a.id, u.username, u.nickname, u.profile_picture_uri
		FROM users u
		JOIN accounts a ON u.account_id = a.id
		WHERE u.username = ?`
	err := u.db.QueryRowContext(ctx, q, uid).Scan(&ud.Id, &ud.Username, &ud.Nickname, &ud.ProfilePicUri)
	if err != nil {
		return nil, err
	}

	return &ud, nil
}

func (u *Users) UpdateNickname(uid int, nickname string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := `UPDATE users SET nickname = ? WHERE id = ?`
	stmt, err := u.db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, uid, nickname)
	return err
}

func (u *Users) UpdateProfilePicUri(uid int, uri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := `UPDATE users SET profile_picture_uri = ? WHERE id = ?`
	stmt, err := u.db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, uid, uri)
	return err
}
