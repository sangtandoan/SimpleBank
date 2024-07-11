package gapi

import (
	"context"
	"fmt"

	"github.com/FrostJ143/simplebank/internal/query"
	"github.com/FrostJ143/simplebank/internal/token"
	"github.com/FrostJ143/simplebank/internal/utils"
	"github.com/FrostJ143/simplebank/pb"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	store      query.Store
	config     utils.Config
	tokenMaker token.Maker
}

func NewServer(store query.Store, config utils.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := Server{store: store, config: config, tokenMaker: tokenMaker}

	return &server, nil
}

func (server *Server) Sum(ctx context.Context, req *pb.SumRequest) (*pb.SumResponse, error) {
	result := req.GetNum1() + req.GetNum2()
	res := &pb.SumResponse{Result: result}

	return res, nil
}
