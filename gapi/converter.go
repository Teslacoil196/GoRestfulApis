package gapi

import (
	db "TeslaCoil196/db/sqlc"
	"TeslaCoil196/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:      user.Username,
		FullName:      user.FullName,
		Email:         user.Email,
		LastPassReset: timestamppb.New(user.LastPassReset),
		CreatedAt:     timestamppb.New(user.CreatedAt),
	}
}
