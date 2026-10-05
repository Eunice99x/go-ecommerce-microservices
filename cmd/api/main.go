package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/eunice99x/goMicro/cmd/config"
	"github.com/eunice99x/goMicro/grpc/client"
	"github.com/eunice99x/goMicro/grpc/creds"
	"github.com/eunice99x/goMicro/grpc/pb"
	"github.com/eunice99x/goMicro/internal/handler"
	"github.com/eunice99x/goMicro/internal/pkg/auth"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("startup failed: %v", err)
	}
}

func run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	tokenGen := auth.DefaultJWTConfig(cfg.SecretKey)

	transportCreds, err := creds.Client(cfg.GRPCTLSCA, cfg.GRPCTLSCert, cfg.GRPCTLSKey)
	if err != nil {
		return fmt.Errorf("failed to load tls credentials: %w", err)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(transportCreds),
		grpc.WithUnaryInterceptor(client.ErrorInterceptor),
	}

	conn, err := grpc.NewClient(cfg.GRPCTarget(), opts...)
	if err != nil {
		return fmt.Errorf("failed to create grpc client: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error closing grpc connection: %v", err)
		}
	}()

	hld := handler.NewHandler(client.New(pb.NewEcommClient(conn)))

	srv := &http.Server{
		Addr:              cfg.Addr(),
		Handler:           handler.RegisterRoutes(hld, tokenGen),
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		ReadTimeout:       config.ReadTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// buffered so the goroutine can exit even if nobody reads the error
	srvErr := make(chan error, 1)

	go func() {
		log.Printf("listening on %s", srv.Addr)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srvErr <- err
		}
	}()

	select {
	case err := <-srvErr:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		log.Println("shutdown signal received")
	}

	// ctx is already cancelled by the signal, so shutdown needs its own
	shutdownCtx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	log.Println("shutdown complete")

	return nil
}
