package reposistory

import (
	"Shopee_UMS/storage"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-testfixtures/testfixtures/v3"
	"log"
	"os"
	"testing"
)

var (
	db       *sql.DB
	fixtures *testfixtures.Loader
)

func TestMain(m *testing.M) {
	dbPort := os.Getenv("DATABASE_PORT")
	testDB := os.Getenv("TEST_DATABASE_NAME")

	var err error
	db, err = storage.NewDB("mysql",
		fmt.Sprintf("testuser:test@tcp(localhost:%s)/%s", dbPort, testDB),
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

	os.Exit(m.Run())
}

func prepareTestDB() {
	if err := fixtures.Load(); err != nil {
		log.Fatal("could not load fixtures", err)
	}
}
