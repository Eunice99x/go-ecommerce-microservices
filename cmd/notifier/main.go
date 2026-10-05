package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eunice99x/goMicro/cmd/config"
	"github.com/eunice99x/goMicro/grpc/client"
	"github.com/eunice99x/goMicro/grpc/creds"
	"github.com/eunice99x/goMicro/grpc/pb"
	"github.com/eunice99x/goMicro/internal/notifier"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("notifier failed: %v", err)
	}
}

func run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	transportCreds, err := creds.Client(cfg.GRPCTLSCA, cfg.GRPCTLSCert, cfg.GRPCTLSKey)
	if err != nil {
		return fmt.Errorf("failed to load tls credentials: %w", err)
	}

	conn, err := grpc.NewClient(
		cfg.GRPCTarget(),
		grpc.WithTransportCredentials(transportCreds),
		grpc.WithUnaryInterceptor(client.ErrorInterceptor),
	)
	if err != nil {
		return fmt.Errorf("failed to create grpc client: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error closing grpc connection: %v", err)
		}
	}()

	var mailer notifier.Mailer = notifier.LogMailer{}
	if cfg.SMTPHost != "" {
		mailer = notifier.NewSMTPMailer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPFrom)
	} else {
		log.Println("SMTP_HOST not set, emails will be logged instead of sent")
	}

	worker := notifier.NewWorker(client.New(pb.NewEcommClient(conn)), mailer, notifier.Config{
		Interval:      10 * time.Second,
		BatchSize:     50,
		MaxConcurrent: 10,
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("notifier started")

	worker.Run(ctx)

	log.Println("notifier stopped")

	return nil
}
