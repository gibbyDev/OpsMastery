# Makefile

.PHONY: run run-db init-db docker-up-dev docker-up-prod init-go wait-for-db

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
	@while ! nc -z localhost 9920; do \
		sleep 1; \
	done
	@echo "Database is ready!"

run-db:
	docker run --name postgres-db -e POSTGRES_USER=gorm -e POSTGRES_PASSWORD=gorm -e POSTGRES_DB=gorm -p 9920:5432 -d postgres:latest

init-db:
	docker exec -it postgres-db psql -U gorm -d gorm

docker-up-dev:
	docker compose -f docker-compose.dev.yml up --build

docker-up-prod:
	docker compose -f docker-compose.prop.yml up --build

init-go:
	@if [ ! -f go.mod ]; then \
		go mod init github.com/gibbyDev/OpsMastery; \
	fi
	go mod tidy