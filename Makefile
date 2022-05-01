SHELL=/bin/bash

run:
	set -a && . ./config_local.env && set +a && \
	go run .

test:
	go test --cover -coverprofile cover.out -count=1 -v ./...
	go tool cover -html=cover.out 

generate:
	go generate ./...

tools:
	go install github.com/vektra/mockery/v2@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

dcmd=docker-compose -f docker-compose.yml
dep-up:
	${dcmd} up -d

dep-down:
	${dcmd} down

dep-stop:
	${dcmd} stop

dep-logs:
	${dcmd} logs

dep-clean:
	docker system prune -f;
	docker volume prune -f

export MIGRATION_DIR=./migrations
migration-create:
	mkdir -p $$MIGRATION_DIR && \
	goose -dir=$$MIGRATION_DIR create t sql

migrate-up:
	set -a && . ./config_local.env && set +a && \
	export connstring="user=$$MINI_ASPIRE_API_DATABASE_MASTER_USERNAME password=$$MINI_ASPIRE_API_DATABASE_MASTER_PASSWORD host=$$MINI_ASPIRE_API_DATABASE_MASTER_HOST port=$$MINI_ASPIRE_API_DATABASE_MASTER_PORT dbname=$$MINI_ASPIRE_API_DATABASE_MASTER_DB sslmode=$$MINI_ASPIRE_API_DATABASE_MASTER_SSL_MODE" && \
	goose -dir=$$MIGRATION_DIR postgres "$$connstring" up

migrate-down:
	set -a && . ./config_local.env && set +a && \
	export connstring="user=$$MINI_ASPIRE_API_DATABASE_MASTER_USERNAME password=$$MINI_ASPIRE_API_DATABASE_MASTER_PASSWORD host=$$MINI_ASPIRE_API_DATABASE_MASTER_HOST port=$$MINI_ASPIRE_API_DATABASE_MASTER_PORT dbname=$$MINI_ASPIRE_API_DATABASE_MASTER_DB sslmode=$$MINI_ASPIRE_API_DATABASE_MASTER_SSL_MODE" && \
	goose -dir=$$MIGRATION_DIR postgres "$$connstring" down