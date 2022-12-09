package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

const (
	DBDriver = "postgres"
	DBSource = "postgresql://postgres:123456@localhost:5433/bank?sslmode=disable"
)

func TestMain(m *testing.M) {

	testDB, err := sql.Open(DBDriver, DBSource)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
