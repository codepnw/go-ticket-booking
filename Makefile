include dev.env

MIGRATIONS_PATH = ./internal/store/migrations

run:
	@go run cmd/main.go

dockerup:
	@docker compose --env-file dev.env up -d

migrate-create:
	@migrate create -ext sql -dir internal/store/migrations -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@migrate -database $(DB_ADDR) -path $(MIGRATIONS_PATH) up 

migrate-down:
	@migrate -database $(DB_ADDR) -path $(MIGRATIONS_PATH) down $(filter-out $@,$(MAKECMDGOALS))