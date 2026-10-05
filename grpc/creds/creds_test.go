package creds

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type certFiles struct{ ca, cert, key string }

// newCA returns a CA plus a function that issues leaf certs signed by it into dir
func newCA(t *testing.T, dir, name string) func(leaf string) certFiles {
	t.Helper()

	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	caTmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: name},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTmpl, caTmpl, &caKey.PublicKey, caKey)
	require.NoError(t, err)

	caFile := writePEM(t, filepath.Join(dir, name+"-ca.pem"), "CERTIFICATE", caDER)

	return func(leaf string) certFiles {
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		require.NoError(t, err)

		tmpl := &x509.Certificate{
			SerialNumber: big.NewInt(time.Now().UnixNano()),
			Subject:      pkix.Name{CommonName: leaf},
			DNSNames:     []string{"localhost"},
			IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
			NotBefore:    time.Now().Add(-time.Hour),
			NotAfter:     time.Now().Add(time.Hour),
			KeyUsage:     x509.KeyUsageDigitalSignature,
			ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		}
		der, err := x509.CreateCertificate(rand.Reader, tmpl, caTmpl, &key.PublicKey, caKey)
		require.NoError(t, err)

		keyDER, err := x509.MarshalECPrivateKey(key)
		require.NoError(t, err)

		return certFiles{
			ca:   caFile,
			cert: writePEM(t, filepath.Join(dir, name+"-"+leaf+".pem"), "CERTIFICATE", der),
			key:  writePEM(t, filepath.Join(dir, name+"-"+leaf+"-key.pem"), "EC PRIVATE KEY", keyDER),
		}
	}
}

func writePEM(t *testing.T, path, typ string, der []byte) string {
	t.Helper()

	require.NoError(t, os.WriteFile(path, pem.EncodeToMemory(&pem.Block{Type: typ, Bytes: der}), 0o600))

	return path
}

// startServer runs a gRPC server with only the health service and returns its address
func startServer(t *testing.T, c credentials.TransportCredentials) string {
	t.Helper()

	var lc net.ListenConfig
	lis, err := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)

	srv := grpc.NewServer(grpc.Creds(c))
	healthpb.RegisterHealthServer(srv, health.NewServer())

	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)

	return lis.Addr().String()
}

func check(t *testing.T, addr string, c credentials.TransportCredentials) error {
	t.Helper()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(c))
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()

	_, err = healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{})

	return err
}

func TestMutualTLS(t *testing.T) {
	dir := t.TempDir()
	issue := newCA(t, dir, "ours")
	srvFiles := issue("server")

	serverCreds, err := Server(srvFiles.ca, srvFiles.cert, srvFiles.key)
	require.NoError(t, err)

	addr := startServer(t, serverCreds)

	t.Run("client with a cert from our CA is accepted", func(t *testing.T) {
		f := issue("api")
		c, err := Client(f.ca, f.cert, f.key)
		require.NoError(t, err)

		require.NoError(t, check(t, addr, c))
	})

	t.Run("plaintext client is rejected", func(t *testing.T) {
		require.Error(t, check(t, addr, insecure.NewCredentials()))
	})

	t.Run("client with a cert from another CA is rejected", func(t *testing.T) {
		f := newCA(t, dir, "other")("api")

		// trusts our server, but presents a cert our server doesn't trust
		c, err := Client(srvFiles.ca, f.cert, f.key)
		require.NoError(t, err)

		require.Error(t, check(t, addr, c))
	})
}

func TestLoadErrors(t *testing.T) {
	_, err := Server("missing-ca.pem", "missing.pem", "missing-key.pem")
	require.Error(t, err)

	_, err = Client("missing-ca.pem", "missing.pem", "missing-key.pem")
	require.Error(t, err)
}
