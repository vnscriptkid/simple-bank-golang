package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"github.com/vnscriptkid/simple-bank-golang/api"
	db "github.com/vnscriptkid/simple-bank-golang/db/sqlc"
	"github.com/vnscriptkid/simple-bank-golang/gapi"
	"github.com/vnscriptkid/simple-bank-golang/pb"
	"github.com/vnscriptkid/simple-bank-golang/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	runGrpcServer(config, *store)
}

func runHttpServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal("can not create http server: ", err.Error())
	}

	err = server.Start(config.HttpServerAddress)

	if err != nil {
		log.Fatal("cannot start http server at " + config.HttpServerAddress)
	}

	fmt.Println("http server started at " + config.HttpServerAddress)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)

	if err != nil {
		log.Fatal("can not create new grpc server")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)

	if err != nil {
		log.Fatal("cannot create grpc listener")
	}

	fmt.Printf("start gRPC server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("cannot start gRPC server")
	}
}
