package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/brianvoe/gofakeit"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

const NumRecords = 10000000
const Offset = 1

var (
	db       *sql.DB
	users    []User
	accounts []Account
	ctx      context.Context
)

type Account struct {
	id        int
	username  string
	password  string
	createdAt string
	updatedAt string
}

type User struct {
	id                int
	accountId         int
	nickname          string
	profilePictureUri string
	createdAt         string
	updatedAt         string
}

func main() {
	gofakeit.Seed(0)

	var err error
	db, err = newDB("mysql", "testuser:test@tcp(localhost:3306)/ums", 10, 10)
	if err != nil {
		log.Fatal("could not connect to MySQL", err)
	}

	for i := Offset; i <= Offset+NumRecords-1; i++ {
		accounts = append(accounts, Account{
			id:        i,
			username:  fmt.Sprintf("user%d", i),
			password:  "$2a$10$Aub2w87Kx/9t3yh4hWcB3.w0A.x6K36X3kKKeSD32pJ5RdFTd312i", // bcrypt of "secret", cost=10
			createdAt: "2018-09-12 22:35:56",
			updatedAt: "2018-09-12 22:35:56",
		})

		users = append(users, User{
			id:                i,
			accountId:         i,
			nickname:          gofakeit.Name(),
			profilePictureUri: "https://shopee-ums-hxt365.s3.ap-southeast-1.amazonaws.com/test_photo.png",
			createdAt:         "2018-09-12 22:35:56",
			updatedAt:         "2018-09-12 22:35:56",
		})
	}

	//if err := bulkInsertDB(accounts, users); err != nil {
	//	log.Fatal(err)
	//}

	writeAccountsToCsv(accounts)
	writeUsersToCsv(users)
	fmt.Println("Done!")
}

func bulkInsertDB(accounts []Account, users []User) error {
	ctx = context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := bulkInsertAccounts(tx, accounts); err != nil {
		tx.Rollback()
		return err
	}
	if err := bulkInsertUsers(tx, users); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func bulkInsertAccounts(tx *sql.Tx, accounts []Account) error {
	vals := []interface{}{}
	q := `INSERT INTO accounts(username, password) VALUES `
	for _, a := range accounts {
		q += "(?, ?),"
		vals = append(vals, a.username, a.password)
	}
	q = q[:len(q)-1]
	return bulkInsert(tx, q, vals)
}

func bulkInsertUsers(tx *sql.Tx, users []User) error {
	vals := []interface{}{}
	q := `INSERT INTO users(account_id, nickname, profile_picture_uri) VALUES `
	for _, u := range users {
		q += "(?, ?, ?),"
		vals = append(vals, u.accountId, u.nickname, u.profilePictureUri)
	}
	q = q[:len(q)-1]
	return bulkInsert(tx, q, vals)
}

func bulkInsert(tx *sql.Tx, q string, vals []interface{}) error {
	stmt, err := tx.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(vals...); err != nil {
		return err
	}

	return nil
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

func writeAccountsToCsv(accounts []Account) {
	records := [][]string{
		{"id", "username", "password", "created_at", "updated_at"},
	}
	for _, a := range accounts {
		records = append(records, []string{
			fmt.Sprint(a.id), a.username, a.password, a.createdAt, a.updatedAt,
		})
	}

	f, err := os.Create("./data/accounts.csv")
	defer f.Close()

	if err != nil {
		log.Fatal("failed to open file", err)
	}

	w := csv.NewWriter(f)
	err = w.WriteAll(records)

	if err != nil {
		log.Fatal(err)
	}
}

func writeUsersToCsv(users []User) {
	records := [][]string{
		{"id", "account_id", "nickname", "profile_picture_uri", "created_at", "updated_at"},
	}
	for _, u := range users {
		records = append(records, []string{
			fmt.Sprint(u.id), fmt.Sprint(u.accountId), u.nickname, u.profilePictureUri, u.createdAt, u.updatedAt,
		})
	}

	f, err := os.Create("./data/users.csv")
	defer f.Close()

	if err != nil {
		log.Fatal("failed to open file", err)
	}

	w := csv.NewWriter(f)
	err = w.WriteAll(records)

	if err != nil {
		log.Fatal(err)
	}
}