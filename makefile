DB_CONTAINER=go-micro
DB_NAME=ecomm
DB_USER=postgres
DB_PASSWORD=postgres
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@localhost:5433/$(DB_NAME)?sslmode=disable

go-run:
	go run ./cmd/api/main.go

db-up:
	docker start $(DB_CONTAINER)

db-down:
	docker stop $(DB_CONTAINER)

db-shell:
	docker exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME)

migrate-up:
	docker run --rm --network host -v ./db:/db migrate/migrate -path=/db/migrations -database "$(DB_URL)" up

migrate-down:
	docker run --rm --network host -v ./db:/db migrate/migrate -path=/db/migrations -database "$(DB_URL)" down 1

test:
	go test ./...

test-race:
	go test ./... -race -count=1

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out | tail -1

lint:
	golangci-lint run

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		grpc/pb/api.proto

grpc-run:
	go run ./cmd/grpc/main.go

notifier-run:
	go run ./cmd/notifier/main.go

# full stack in containers (postgres, migrations, grpc, api, notifier, mailpit)
certs:
	./scripts/dev-certs.sh

up: certs
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f api grpc notifier

.PHONY: go-run grpc-run notifier-run db-up db-down db-shell migrate-up migrate-down test test-race test-cover lint proto certs up down logs
