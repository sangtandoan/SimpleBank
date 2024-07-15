package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	_ "github.com/FrostJ143/simplebank/doc/statik"
	"github.com/FrostJ143/simplebank/internal/api"
	"github.com/FrostJ143/simplebank/internal/gapi"
	"github.com/FrostJ143/simplebank/internal/query"
	"github.com/FrostJ143/simplebank/internal/utils"
	"github.com/FrostJ143/simplebank/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	go runGatewayServer(store, config)
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

	log.Println("start gRPC server at ", listener.Addr().String())
	err = gRPCServer.Serve(listener)
	if err != nil {
		log.Fatal("can not start gRPC server: ", err)
	}
}

func runGatewayServer(store query.Store, config utils.Config) {
	serverHandler, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal("cannot create server handler: ", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, serverHandler)
	if err != nil {
		log.Fatal("cannot register server handler to gateway")
	}

	serverMux := http.NewServeMux()
	serverMux.Handle("/", grpcMux)

	// Remember that Go is a compiled language; most everything the program does happens at runtime.
	// In particular, in this case, the call to http.Dir() happens at runtime, and that means that the path is evaluated at runtime.
	// Because the path you have given is relative, it is therefore relative to the working directory from which you run the application.
	// The directory in which the source code resided is not relevant here.
	//
	// In Go, relative paths are resolved based on the current working directory, which is where you run your Go program from,
	// not necessarily where the source code or the project is located.

	// fs := http.FileServer(http.Dir("../doc/swagger")) not working
	// fs := http.FileServer(http.Dir("./doc/swagger"))
	// serverMux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
	//

	// Embed static frontend files inside Golang backend server's binary
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal("cannot create statik fs: ", err)
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	serverMux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}

	log.Println("start gRPC gateway server at ", config.HTTPServerAddress)
	err = http.Serve(listener, serverMux)
	if err != nil {
		log.Fatal("cannot start HTTP server gateway: ", err)
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
