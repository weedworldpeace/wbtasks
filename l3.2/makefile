.PHONY: migrate-up migrate-down start stop test

migrate-up:
	migrate -path app/migrations -database "postgres://shortener:shortener_pass@localhost:5433/url_shortener?sslmode=disable" up

migrate-down:
	migrate -path app/migrations -database "postgres://shortener:shortener_pass@localhost:5433/url_shortener?sslmode=disable" down

start:
	docker-compose up -d

stop:
	docker-compose down

