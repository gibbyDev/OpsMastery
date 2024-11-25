# Makefile

.PHONY: run-db init-db

run-db:
	docker run --name postgres-db -e POSTGRES_USER=gorm -e POSTGRES_PASSWORD=gorm -e POSTGRES_DB=gorm -p 9920:5432 -d postgres:latest

init-db:
	docker exec -it postgres-db psql -U gorm -d gorm 