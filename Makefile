.PHONY: up down build-cli test lint bench security-scan profile generate help

# default target
help:
	@echo "auspicious gopher ledger build system"
	@echo "usage: make <target>"
	@echo ""
	@echo "targets:"
	@echo "  up             spin up the entire stack using docker compose"
	@echo "  down           tear down the stack"
	@echo "  test           run all tests with fresh results"
	@echo "  lint           run golangci-lint for a clean gopher soul"
	@echo "  bench          benchmark the service layer for performant kindness"
	@echo "  security-scan  run govulncheck to identify known vulnerabilities"
	@echo "  profile        collect a 30s cpu profile from the local service"
	@echo "  generate       regenerate all mocks and swagger documentation"
	@echo "  build-cli      compile the auspicious cli tool"

up:
	docker compose up --build

down:
	docker compose down

test:
	cd go-app && go test -v -count=1 ./...

lint:
	cd go-app && golangci-lint run ./...

bench:
	cd go-app && go test -bench=. -benchmem ./internal/service/...

security-scan:
	cd go-app && govulncheck ./...

profile:
	curl -s http://localhost:8000/debug/pprof/profile?seconds=30 > cpu.prof
	@echo "cpu profile saved to cpu.prof"

generate:
	cd go-app && swag init -g cmd/server/main.go
	cd go-app && go run go.uber.org/mock/mockgen -source=internal/domain/incident.go -destination=internal/service/mock_incident_repository.go -package=service
	cd go-app && go run go.uber.org/mock/mockgen -source=internal/service/pkgsite.go -destination=internal/service/mock_pkgsite_service.go -package=service
	cd go-app && go run go.uber.org/mock/mockgen -source=internal/service/incident.go -destination=internal/handler/mock_incident_service.go -package=handler

build-cli:
	cd cli && go build -o ../ledger-cli .
