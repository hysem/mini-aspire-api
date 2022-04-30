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
