COMPOSE ?= docker-compose -f ops/docker-compose.yml

run:
	$(COMPOSE) up --build --force-recreate -d

rm:
	$(COMPOSE) rm -sfv

logs:
	docker logs ops-app-1 -f

init:
	go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
	go install github.com/swaggo/swag/cmd/swag@latest

fmt:
	@swag fmt

generate-sql:
	sqlc generate

generate-swagger:
	@swag init -g ./cmd/main.go

generate: generate-sql generate-swagger

setup-db:
	docker run --name test-db \
	-e POSTGRES_USER=test \
	-e POSTGRES_PASSWORD=test \
	-e POSTGRES_DB=test \
	-p 5432:5432 \
	-d \
	postgres:latest

cleanup-db:
	docker rm -f test-db
