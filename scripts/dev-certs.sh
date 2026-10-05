#!/usr/bin/env bash
# Generates a local CA and one certificate per service for gRPC mutual TLS.
# Development only: in production certificates come from a real CA or a service mesh.
set -euo pipefail

dir="${1:-certs}"
days=365

if [[ -f "$dir/ca.pem" ]]; then
    echo "certs already exist in $dir (delete it to regenerate)"
    exit 0
fi

mkdir -p "$dir"
cd "$dir"

openssl ecparam -name prime256v1 -genkey -noout -out ca-key.pem
openssl req -x509 -new -key ca-key.pem -sha256 -days "$days" -subj "/CN=gochi dev ca" -out ca.pem

# issue <name> <subjectAltName>
issue() {
    local name="$1" san="$2"

    openssl ecparam -name prime256v1 -genkey -noout -out "$name-key.pem"
    openssl req -new -key "$name-key.pem" -subj "/CN=$name" -out "$name.csr"
    openssl x509 -req -in "$name.csr" -CA ca.pem -CAkey ca-key.pem -CAcreateserial \
        -days "$days" -sha256 -out "$name.pem" \
        -extfile <(printf "subjectAltName=%s\nextendedKeyUsage=serverAuth,clientAuth\nkeyUsage=digitalSignature\n" "$san")
    rm "$name.csr"
}

# the server cert must match the hostnames clients dial: "grpc" in compose, localhost with go run
issue grpc "DNS:grpc,DNS:localhost,IP:127.0.0.1"
issue api "DNS:api"
issue notifier "DNS:notifier"

rm -f ca.srl
# containers run as a non-root user and read these through a bind mount
chmod 644 ./*.pem
chmod 600 ca-key.pem

echo "generated dev certs in $dir"
