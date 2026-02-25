.PHONY: migrate-up migrate-down start stop test

migrate-up:
	migrate -path app/migrations -database "postgres://user:12345@localhost:5433/transactions?sslmode=disable" up

migrate-down:
	migrate -path app/migrations -database "postgres://user:12345@localhost:5433/transactions?sslmode=disable" down

start:
	docker-compose up -d

stop:
	docker-compose down

