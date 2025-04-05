# Makefile

.PHONY: run run-db init-db docker-up-dev docker-up-prod init-go wait-for-db

# Load environment variables from .env file
ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

run:
	@mkdir -p tmp
	$(MAKE) docker-up-dev
	$(MAKE) wait-for-db
	$(MAKE) init-go
	@if command -v air > /dev/null; then \
		air; \
	else \
		go run main.go; \
	fi

wait-for-db:
	@echo "Waiting for database to be ready..."
	@while ! nc -z $(DB_HOST) $(DB_PORT); do \
		sleep 1; \
	done
	@echo "Database is ready!"

run-db:
	docker run --name postgres-db -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -e POSTGRES_DB=$(DB_NAME) -p $(DB_PORT):5432 -d postgres:latest

init-db:
	docker exec -it postgres-db psql -U $(DB_USER) -d $(DB_NAME)

docker-up-dev:
	docker compose -f docker-compose.dev.yml up --build

docker-up-prod:
	docker compose -f docker-compose.prop.yml up --build

init-go:
	@if [ ! -f go.mod ]; then \
		go mod init github.com/gibbyDev/OpsMastery; \
	fi
	go mod tidy