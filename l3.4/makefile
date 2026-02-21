.PHONY: migrate-up migrate-down start stop test

migrate-up:
	migrate -path processor/migrations -database "postgres://user:password@localhost:5433/tasks?sslmode=disable" up

migrate-down:
	migrate -path processor/migrations -database "postgres://user:password@localhost:5433/tasks?sslmode=disable" down

start:
	docker-compose up -d

stop:
	docker-compose down

