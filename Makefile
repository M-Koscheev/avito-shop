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
	migrate -path db/migrations -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable -verbose up

migratedown:
	migrate -path db/migrations -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable -verbose up

fixmigrate:
	migrate -path db/migrations -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable force 1

postgres:
	docker run --name=avito-shop-db -p 5432:5432 -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d --rm postgres:13

createdb:
	docker exec -it avito-shop-db createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ${POSTGRES_DB}

dropdb:
	docker exec -it avito-shop-db dropdb ${POSTGRES_DB}

swag:
	swag init --parseDependency --parseInternal -g cmd/main.go

lint-go: ## Use golintci-lint on your project
	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "--out-format checkstyle ./... | tee /dev/tty > checkstyle-report.xml" || echo "" ))
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:latest-alpine golangci-lint run --timeout=1000s $(OUTPUT_OPTIONS)


.PHONY: build run