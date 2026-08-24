DB_CONTAINER=go-micro
DB_NAME=ecomm
DB_USER=postgres
DB_PASSWORD=postgres
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable

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