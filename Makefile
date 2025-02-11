include .env

build:
	docker compose build avito-shop-service

update:
	docker compose pull avito-shop-service

run:
	docker compose up avito-shop-service -d

test:
	go test -v ./...

migrateup:
	migrate -path db/migrations -database ${POSTGRESQL_URL} -verbose up

migratedown:
	migrate -path db/migrations -database ${POSTGRESQL_URL} -verbose down

.PHONY: build run