package main

import (
	"TeslaCoil196/api"
	db "TeslaCoil196/db/sqlc"
	"TeslaCoil196/util"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Couldn't connect to databse ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.StartServer(config.ServerAddress)
	if err != nil {
		log.Fatal("Unable to start server ", err)
	}
}
