package reposistory

import (
	database "Shopee_UMS/db"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"testing"
)

var (
	db       *sql.DB
	fixtures *testfixtures.Loader
	mg       *migrate.Migrate
)

func TestMain(m *testing.M) {
	dbPort := os.Getenv("DATABASE_PORT")
	testDB := os.Getenv("TEST_DATABASE_NAME")
	testUser := os.Getenv("DATABASE_USER")
	testPwd := os.Getenv("DATABASE_PASSWORD")

	var err error
	db, err = database.New("mysql",
		fmt.Sprintf("%s:%s@tcp(localhost:%s)/%s", testUser, testPwd, dbPort, testDB),
		10, 10)
	if err != nil {
		log.Fatal("could not connect to test DB", err)
	}

	fixtures, err = testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("mysql"),
		testfixtures.Directory("../testfixtures"),
	)
	if err != nil {
		log.Fatal("could not create fixtures", err)
	}

	mg, err = migrate.New(
		"file://../db/migrations",
		fmt.Sprintf("mysql://%s:%s@tcp(localhost:%s)/%s", testUser, testPwd, dbPort, testDB),
	)
	if err != nil {
		log.Fatal("could not create migrate", err)
	}

	if err := mg.Down(); err != nil && !noChangeErr(err) {
		log.Fatal("could not migrate down", err)
	}
	os.Exit(m.Run())
}

func setup() {
	if err := mg.Up(); err != nil && !noChangeErr(err) {
		log.Fatal("could not migrate up", err)
	}
	if err := fixtures.Load(); err != nil {
		log.Fatal("could not load fixtures", err)
	}
}

func tearDown() {
	if err := mg.Down(); err != nil && !noChangeErr(err) {
		log.Fatal("could not migrate down", err)
	}
}

func noChangeErr(err error) bool {
	return err.Error() == "no change"
}
