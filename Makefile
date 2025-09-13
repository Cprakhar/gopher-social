include .envrc
MIGRATIONS_DIR=./migrate/migrations

.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_DIR) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_DIR) -database="$(DB_ADDR)" up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_DIR) -database="$(DB_ADDR)" down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: seed
seed:
	@go run ./migrate/seed/main.go

.PHONY: gen-docs
gen-docs:
	@swag init -g ./main.go && swag fmt