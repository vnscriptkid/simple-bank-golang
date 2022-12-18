package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/vnscriptkid/simple-bank-golang/api"
	db "github.com/vnscriptkid/simple-bank-golang/db/sqlc"
)

const (
	DBDriver      = "postgres"
	DBSource      = "postgresql://postgres:123456@localhost:5433/bank?sslmode=disable"
	ServerAddress = "0.0.0.0:8000"
)

func main() {
	conn, err := sql.Open(DBDriver, DBSource)

	if err != nil {
		log.Fatal("cannot connect to db")
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(ServerAddress)

	if err != nil {
		log.Fatal("cannot start server at " + ServerAddress)
	}

	fmt.Println("server started at " + ServerAddress)
}
