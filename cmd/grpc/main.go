package main

import (
	"log"
	"net"

	"github.com/eunice99x/goMicro/cmd/config"
	"github.com/eunice99x/goMicro/db"
	"github.com/eunice99x/goMicro/grpc/pb"
	"github.com/eunice99x/goMicro/grpc/service"
	"github.com/eunice99x/goMicro/internal/repository"
	"google.golang.org/grpc"
)


func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Panicf("failed to load configuration: %v", err)
	}
	
	database, err := db.NewDatabase(cfg.DSN())
	if err != nil {
		log.Panicf("error opening db: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	st := repository.NewPostgresStorer(database.GetDB())
	srv := service.NewStorer(st)

	grpcSrv := grpc.NewServer()
	pb.RegisterEcommServer(grpcSrv, srv)
	
	listener, err := net.Listen("tcp", cfg.GRPCAddr())
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", cfg.GRPCAddr(), err)
	}

	log.Printf("gRPC server listening on %s", cfg.GRPCAddr())

	if err := grpcSrv.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}
