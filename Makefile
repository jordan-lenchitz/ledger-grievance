.PHONY: up down build-cli test

up:
	docker compose up --build

down:
	docker compose down

build-cli:
	cd cli && go build -o ../ledger-cli .

test:
	cd go-app && go test -v ./...
