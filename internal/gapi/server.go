package gapi

import (
	"fmt"

	"github.com/FrostJ143/simplebank/internal/query"
	"github.com/FrostJ143/simplebank/internal/token"
	"github.com/FrostJ143/simplebank/internal/utils"
	"github.com/FrostJ143/simplebank/pb"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	store      query.Store
	tokenMaker token.Maker
	config     utils.Config
}

func NewServer(store query.Store, config utils.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := Server{store: store, config: config, tokenMaker: tokenMaker}

	return &server, nil
}
