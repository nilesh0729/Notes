package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_"github.com/lib/pq"
)

const (
	DBDriver = "postgres"
	DBSource = "postgres://root:MummyJi@localhost:5431/PapaJi?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(DBDriver, DBSource)
	if err != nil {
		log.Fatal("Can't connect to DB", err)
	}
	testQueries = New(conn)

	os.Exit(m.Run())

}
