package reposistory

import (
	database "Shopee_UMS/db"
	"Shopee_UMS/utils"
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
	var err error
	db, err = database.NewTestDB()
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

	dbHost := utils.MustEnv("DATABASE_HOST")
	dbPort := utils.MustEnv("DATABASE_PORT")
	testDB := utils.MustEnv("TEST_DATABASE_NAME")
	testUser := utils.MustEnv("DATABASE_USER")
	testPwd := utils.MustEnv("DATABASE_PASSWORD")
	mg, err = migrate.New(
		"file://../db/migrations",
		fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s", testUser, testPwd, dbHost, dbPort, testDB),
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
