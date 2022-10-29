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
