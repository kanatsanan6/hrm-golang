package queries_test

import (
	"log"
	"os"
	"testing"

	"github.com/kanatsanan6/hrm/queries"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	dbSource = "postgresql://postgres:secret@localhost:5432/hrm_test"
)

var testQueries queries.Queries

func TestMain(m *testing.M) {
	db, err := gorm.Open(postgres.Open(dbSource), &gorm.Config{})

	if err != nil {
		log.Fatal("Cannot connect to db", err)
	}

	testQueries = queries.NewQueries(db)

	os.Exit(m.Run())
}
