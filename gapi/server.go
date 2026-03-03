package gapi

import (
	db "TeslaCoil196/db/sqlc"
	"TeslaCoil196/pb"
	"TeslaCoil196/token"
	"TeslaCoil196/util"
	"fmt"
)

// type Server struct {
// 	pb.UnimplementedTeslaBankServer
// 	db         db.Store
// 	config     util.Config
// 	tokenMaker token.Maker
// }

// func NewServer(db db.Store, config util.Config) (*Server, error) {
// 	tokenMaker, err := token.NewPastoMaker(config.SymmetricKey)
// 	if err != nil {
// 		return nil, fmt.Errorf("Could not create tokenMaker %s", err)
// 	}
// 	server := &Server{
// 		config:     config,
// 		db:         db,
// 		tokenMaker: tokenMaker,
// 	}

// 	return server, nil
// }

type Server struct {
	pb.UnimplementedTeslaBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new gRPC server.
func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPastoMaker(config.SymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
