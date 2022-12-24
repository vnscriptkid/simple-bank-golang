package gapi

import (
	"context"
	"database/sql"

	db "github.com/vnscriptkid/simple-bank-golang/db/sqlc"
	"github.com/vnscriptkid/simple-bank-golang/pb"
	"github.com/vnscriptkid/simple-bank-golang/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetUsername())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found %s", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to find user %s", err)
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password %s", err)
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create token %s", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create token %s", err)
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session %s", err)
	}

	rsp := &pb.LoginUserResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return rsp, nil
}
