package main

import (
	"TeslaCoil196/api"
	db "TeslaCoil196/db/sqlc"
	"TeslaCoil196/gapi"
	"TeslaCoil196/pb"
	"TeslaCoil196/util"
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	runGRPCServer(config, store)

}

func runGRPCServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal("Could not create server", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTeslaBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listner, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("Could not create listner", err)
	}

	log.Printf("String GRPC server at %s", listner.Addr().String())

	err = grpcServer.Serve(listner)
	if err != nil {
		log.Fatal("Could not create GRPC Server", err)
	}

}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("Could not create server", err)
	}

	err = server.StartServer(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Unable to start server ", err)
	}
}
