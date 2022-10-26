COMPOSE ?= docker-compose -f ops/docker-compose.yml

run:
	$(COMPOSE) up --build --force-recreate -d

rm:
	$(COMPOSE) rm -sfv

generate-sql:
	sqlc generate

generate: generate-sql
