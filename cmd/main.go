package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/FrostJ143/simplebank/internal/api"
	"github.com/FrostJ143/simplebank/internal/gapi"
	"github.com/FrostJ143/simplebank/internal/query"
	"github.com/FrostJ143/simplebank/internal/utils"
	"github.com/FrostJ143/simplebank/pb"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Could not connect to database: ", err)
	}

	store := query.NewSQLStore(conn)
	runGRPCServer(store, config)
}

func runGRPCServer(store query.Store, config utils.Config) {
	serverHandler, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal("cannot create server handler:", err)
	}

	gRPCServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(gRPCServer, serverHandler)
	reflection.Register(gRPCServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = gRPCServer.Serve(listener)
	if err != nil {
		log.Fatal("can not start gRPC server: ", err)
	}
}

func runGinServer(store query.Store, config utils.Config) {
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("Could not create server: ", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Could not start server")
	}

}
