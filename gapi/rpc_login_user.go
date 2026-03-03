package gapi

import (
	db "TeslaCoil196/db/sqlc"
	"TeslaCoil196/pb"
	"TeslaCoil196/util"
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, request *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := server.store.GetUser(ctx, request.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "Error %s", err)

	}

	err = util.CheckPasswords(request.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Unauthrozied")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.TokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Access token error %s", err)

	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Refresh token error %s", err)
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    " ",
		ClientIp:     " ",
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Refresh token error %s", err)
	}

	reply := &pb.LoginUserResponse{
		SessionId:              session.ID.String(),
		AccessToken:            accessToken,
		AccessTokenExpireTime:  timestamppb.New(accessPayload.ExpiredAt),
		User:                   convertUser(user),
		RefreshToken:           refreshToken,
		RefreshTokenExpireTime: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return reply, nil
}
