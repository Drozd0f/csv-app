COMPOSE ?= docker-compose -f ops/docker-compose.yml

run:
	$(COMPOSE) up --build --force-recreate -d

rm:
	$(COMPOSE) rm -sfv

logs:
	docker logs ops-app-1 -f

generate-sql:
	sqlc generate

fmt:
	@swag fmt

generate-swagger:
	@swag init -g ./cmd/main.go

generate: generate-sql generate-swagger
