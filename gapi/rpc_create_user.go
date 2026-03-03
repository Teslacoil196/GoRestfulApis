package gapi

import (
	db "TeslaCoil196/db/sqlc"
	"TeslaCoil196/pb"
	"TeslaCoil196/util"
	"context"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	hashedPassword, err := util.HashedPassword(request.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error occoured while hashing password : %s", err)
	}

	arg := db.CreateUserParams{
		Username:       request.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       request.GetFullName(),
		Email:          request.GetEmail(),
	}
	//fmt.Println("about to make the call")
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pgerr, ok := err.(*pq.Error); ok {
			switch pgerr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.Internal, "Username already exists")
			}
		}
		return nil, status.Errorf(codes.Internal, "error occoured while hashing password : %s", err)
	}

	response := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return response, nil
}
