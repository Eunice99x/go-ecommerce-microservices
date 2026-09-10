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
	"time"

	"github.com/eunice99x/goMicro/cmd/config"
	"github.com/eunice99x/goMicro/db"
	"github.com/eunice99x/goMicro/internal/handler"
	"github.com/eunice99x/goMicro/internal/pkg/auth"
	"github.com/eunice99x/goMicro/internal/repository"
	"github.com/eunice99x/goMicro/internal/service"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 20 * time.Second
	idleTimeout       = 120 * time.Second
	shutdownTimeout   = 15 * time.Second
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

	database, err := db.NewDatabase(cfg.DSN())
	if err != nil {
		return fmt.Errorf("error opening db: %w", err)
	}
	defer database.Close()

	log.Println("successfully connected to database")

	tokenGen := auth.DefaultJWTConfig(cfg.SecretKey)

	store := repository.NewPostgresStorer(database.GetDB())
	svc := service.NewService(store, tokenGen)
	hld := handler.NewHandler(svc)

	srv := &http.Server{
		Addr:              cfg.Addr(),
		Handler:           handler.RegisterRoutes(hld),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
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
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	log.Println("shutdown complete")

	return nil
}
