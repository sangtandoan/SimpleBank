package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/FrostJ143/simplebank/doc/statik"
	"github.com/FrostJ143/simplebank/internal/api"
	"github.com/FrostJ143/simplebank/internal/email"
	"github.com/FrostJ143/simplebank/internal/gapi"
	"github.com/FrostJ143/simplebank/internal/query"
	"github.com/FrostJ143/simplebank/internal/utils"
	"github.com/FrostJ143/simplebank/internal/worker"
	"github.com/FrostJ143/simplebank/pb"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
		log.Fatal().Msgf("Could not load config: %s", err)
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Msgf("Could not connect to database: %s", err)
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisServerAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	// run db migration
	runDBMigration(config.MigrationURL, config.DBSource)

	store := query.NewSQLStore(conn)
	go runTaskProcessor(redisOpt, store, config)
	go runGatewayServer(store, config, taskDistributor)
	runGRPCServer(store, config, taskDistributor)
}

func runDBMigration(migrationURL string, DBSource string) {
	m, err := migrate.New(migrationURL, DBSource)
	if err != nil {
		log.Fatal().Msgf("cannt create migrate object: %s", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msgf("cannot migrate up: %s", err)
	}

	log.Info().Msg("migrated successfully")
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store query.Store, config utils.Config) {
	mailSender := email.NewGmailSender(
		config.EmailSenderName,
		config.EmailSenderAddress,
		config.EmailSenderPasswrod,
	)
	processor := worker.NewRedisTaskProcessor(redisOpt, store, mailSender)

	log.Info().Msg("start task processor")
	err := processor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

func runGRPCServer(store query.Store, config utils.Config, taskDistributor worker.TaskDistributor) {
	serverHandler, err := gapi.NewServer(store, config, taskDistributor)
	if err != nil {
		log.Fatal().Msgf("cannot create server handler: %s", err)
	}

	logOpt := grpc.UnaryInterceptor(gapi.GrpcLogger)
	gRPCServer := grpc.NewServer(logOpt)
	pb.RegisterSimpleBankServer(gRPCServer, serverHandler)
	reflection.Register(gRPCServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msgf("cannot create listener: %s", err)
	}

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
	err = gRPCServer.Serve(listener)
	if err != nil {
		log.Fatal().Msgf("can not start gRPC server: %s", err)
	}
}

func runGatewayServer(store query.Store, config utils.Config, taskDistributor worker.TaskDistributor) {
	serverHandler, err := gapi.NewServer(store, config, taskDistributor)
	if err != nil {
		log.Fatal().Msgf("cannot create server handler: %s", err)
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
		log.Fatal().Msgf("cannot register server handler to gateway")
	}

	serveMux := http.NewServeMux()
	serveMux.Handle("/", grpcMux)

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
		log.Fatal().Msgf("cannot create statik fs: %s", err)
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	serveMux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msgf("cannot create listener: %s", err)
	}

	log.Info().Msgf("start gRPC gateway server at %s", config.HTTPServerAddress)
	handler := gapi.HttpLogger(serveMux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Msgf("cannot start HTTP server gateway: %s", err)
	}
}

func runGinServer(store query.Store, config utils.Config) {
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal().Msgf("Could not create server: %s", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msgf("Could not start server")
	}
}
