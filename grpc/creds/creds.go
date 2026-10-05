// Package creds sets up mutual TLS between our services, or plaintext when no certs are configured.
package creds

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Server only accepts clients with a cert signed by our CA
func Server(caFile, certFile, keyFile string) (credentials.TransportCredentials, error) {
	if caFile == "" {
		log.Println("WARNING: gRPC server running without TLS; set GRPC_TLS_* outside local development")
		return insecure.NewCredentials(), nil
	}

	cert, pool, err := load(caFile, certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    pool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
	}), nil
}

// Client verifies the server and sends our own cert
func Client(caFile, certFile, keyFile string) (credentials.TransportCredentials, error) {
	if caFile == "" {
		log.Println("WARNING: gRPC client connecting without TLS; set GRPC_TLS_* outside local development")
		return insecure.NewCredentials(), nil
	}

	cert, pool, err := load(caFile, certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
		MinVersion:   tls.VersionTLS13,
	}), nil
}

func load(caFile, certFile, keyFile string) (tls.Certificate, *x509.CertPool, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return tls.Certificate{}, nil, fmt.Errorf("error loading tls key pair: %w", err)
	}

	ca, err := os.ReadFile(caFile)
	if err != nil {
		return tls.Certificate{}, nil, fmt.Errorf("error reading ca: %w", err)
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(ca) {
		return tls.Certificate{}, nil, fmt.Errorf("no certificates found in %s", caFile)
	}

	return cert, pool, nil
}
