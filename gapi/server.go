package gapi

import (
	"fmt"

	db "github.com/vnscriptkid/simple-bank-golang/db/sqlc"
	"github.com/vnscriptkid/simple-bank-golang/pb"
	"github.com/vnscriptkid/simple-bank-golang/token"
	"github.com/vnscriptkid/simple-bank-golang/util"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new gRPC server.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
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
