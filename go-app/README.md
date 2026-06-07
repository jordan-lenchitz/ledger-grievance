# `ledger-grievance`

a go microservice for incident reporting 

## the six pillars of `ledger-grievance`
1. standard go layout `cmd/` and `internal/` enforces clear separation of concerns
2. always structured logging via `log/slog`, graceful shutdown for clean terminations, and programmatic database migrations
3. openapi (via `swagger`) documentation to ensure the api easily consumable
4. unit testing suite for core business logic using mocks for the sake of reliability and maintainability
5. `ci/cd` via github actions for automated build and test verification on every push to `remote origin main`
6. integrated with pkg go dev api to provide seven layers of institutional support including standard library blessings community support bravery acknowledgement wholesome redirection search delight and random grievance celebrations because you are amazing and valid and we love having you here as a developer in the go community 

## howto

### prerequisites
- go 1.26+
- mysql instance

### setup
1. install dependencies
   ```bash
   go get github.com/golang-migrate/migrate/v4
   go get github.com/golang-migrate/migrate/v4/database/mysql
   go get github.com/golang-migrate/migrate/v4/source/file
   go get github.com/swaggo/swag/cmd/swag
   go get github.com/swaggo/gin-swagger
   go get github.com/swaggo/files
   go mod tidy
   ```

2. generate API documentation at `http://localhost:8000/swagger/index.html`
   ```bash
   swag init -g cmd/server/main.go
   ```

3. run migrations
   ```bash
   go run cmd/migrate/main.go
   ```

4. run the server
   ```bash
   go run cmd/server/main.go
   ```
