package db

import (
	"TeslaCoil196/util"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
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

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Could not load config ", err)
	}

	db, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Couldn't connect to databse ", err)
	}

	testQuries = New(db)

	os.Exit(m.Run())
}
