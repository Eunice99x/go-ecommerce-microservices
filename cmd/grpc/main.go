package main

import (
	"context"
	"log"
	"net"

	"github.com/eunice99x/goMicro/cmd/config"
	"github.com/eunice99x/goMicro/db"
	"github.com/eunice99x/goMicro/grpc/creds"
	"github.com/eunice99x/goMicro/grpc/pb"
	"github.com/eunice99x/goMicro/grpc/server"
	"github.com/eunice99x/goMicro/internal/pkg/auth"
	"github.com/eunice99x/goMicro/internal/repository"
	"github.com/eunice99x/goMicro/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
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

	// tokens are issued here, the api only verifies them
	tokenGen := auth.DefaultJWTConfig(cfg.SecretKey)

	st := repository.NewPostgresStorer(database.GetDB())
	svc := service.NewService(st, tokenGen)
	srv := server.NewServer(svc)

	transportCreds, err := creds.Server(cfg.GRPCTLSCA, cfg.GRPCTLSCert, cfg.GRPCTLSKey)
	if err != nil {
		log.Fatalf("failed to load tls credentials: %v", err)
	}

	grpcSrv := grpc.NewServer(
		grpc.Creds(transportCreds),
		grpc.UnaryInterceptor(server.ErrorInterceptor),
	)
	pb.RegisterEcommServer(grpcSrv, srv)

	// standard grpc health check, for docker/k8s probes
	healthpb.RegisterHealthServer(grpcSrv, health.NewServer())

	var lc net.ListenConfig

	listener, err := lc.Listen(context.Background(), "tcp", cfg.GRPCAddr())
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", cfg.GRPCAddr(), err)
	}

	log.Printf("gRPC server listening on %s", cfg.GRPCAddr())

	if err := grpcSrv.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}
