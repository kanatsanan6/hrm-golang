package queries_test

import (
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/kanatsanan6/hrm/queries"
	_ "github.com/lib/pq"
)

const (
	dbSource = "postgresql://postgres:secret@localhost:5432/hrm_development?sslmode=disable"
)

var testQueries queries.Queries

func TestMain(m *testing.M) {
	db, err := sqlx.Open("postgres", dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db", err)
	}

	testQueries = queries.NewQueries(db)
	os.Exit(m.Run())
}
