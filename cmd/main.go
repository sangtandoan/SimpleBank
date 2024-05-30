package main

import (
	"database/sql"
	"log"

	"github.com/FrostJ143/simplebank/internal/api"
	"github.com/FrostJ143/simplebank/internal/query"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	address  = "localhost:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Could not connect to database: ", err)
	}

	store := query.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(address)
	if err != nil {
		log.Fatal("Could not start server")
	}
}
