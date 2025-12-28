package Database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_"github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:MummyJi@localhost:5431/PapaJi?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil{
		log.Fatal("cannot connect to DB", err)
	}
	testQueries = New(conn)
	os.Exit(m.Run())
}
