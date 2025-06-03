include dev.env

MIGRATIONS_PATH = ./internal/database/migrations

run:
	@go run cmd/main.go

# Docker
dockerup:
	@docker compose --env-file dev.env up -d

createdb:
	@docker exec -it $(CONTAINER_NAME) psql -U $(PG_USER) -c "CREATE DATABASE bookingdb"

dropdb:
	@docker exec -it $(CONTAINER_NAME) psql -U $(PG_USER) -c "DROP DATABASE IF EXISTS bookingdb"

# Migrations
migrate-create:
	@migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@migrate -database $(DB_ADDR) -path $(MIGRATIONS_PATH) up 

migrate-down:
	@migrate -database $(DB_ADDR) -path $(MIGRATIONS_PATH) down $(filter-out $@,$(MAKECMDGOALS))