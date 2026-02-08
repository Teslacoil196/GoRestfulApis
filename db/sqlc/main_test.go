package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:root@localhost:5432/TeslaBank?sslmode=disable"
)

var (
	TestAccount1 Account
	TestAccount2 Account
	TestAccount3 Account
)

var accountsCreatedForTrasnfers = false
var accountCreatedForEntry = false
var testQuries *Queries
var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Couldn't connect to databse ", err)
	}

	testQuries = New(db)

	os.Exit(m.Run())
}
