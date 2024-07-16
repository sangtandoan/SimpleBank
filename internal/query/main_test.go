package query

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/FrostJ143/simplebank/internal/utils"
	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB      *sql.DB
	testStore   Store
)

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("./../..")

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)
	testStore = NewSQLStore(testDB)

	os.Exit(m.Run())
}
