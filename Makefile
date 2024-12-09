# Makefile

.PHONY: run run-db init-db

run:
	@mkdir -p tmp
	@if command -v air > /dev/null; then \
		air; \
	else \
		go run main.go; \
	fi

run-db:
	docker run --name postgres-db -e POSTGRES_USER=gorm -e POSTGRES_PASSWORD=gorm -e POSTGRES_DB=gorm -p 9920:5432 -d postgres:latest

init-db:
	docker exec -it postgres-db psql -U gorm -d gorm 