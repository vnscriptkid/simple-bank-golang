package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/vnscriptkid/simple-bank-golang/api"
	db "github.com/vnscriptkid/simple-bank-golang/db/sqlc"
	"github.com/vnscriptkid/simple-bank-golang/util"
)

const (
	DBDriver      = "postgres"
	DBSource      = "postgresql://postgres:123456@localhost:5433/bank?sslmode=disable"
	ServerAddress = "0.0.0.0:8000"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot load config")
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db")
	}

	store := db.NewStore(conn)

	server, err := api.NewServer(config, *store)

	if err != nil {
		log.Fatal("can not create new server: ", err.Error())
	}

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server at " + config.ServerAddress)
	}

	fmt.Println("server started at " + config.ServerAddress)
}
