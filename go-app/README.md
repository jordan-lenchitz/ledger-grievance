# `ledger-grievance`

a go microservice for incident reporting 

## the ten core beliefs of `ledger-grievance`
1. standard go layout `cmd/` and `internal/` enforces clear separation of concerns
2. always structured logging via `log/slog`, graceful shutdown for clean terminations, and programmatic database migrations
3. openapi (via `swagger`) documentation to ensure the api is easily consumable
4. unit testing suite for core business logic using mocks for the sake of reliability and maintainability
5. `ci/cd` via github actions for automated build and test verification on every push to `remote origin main`
6. integrated with pkg go dev api to provide seven layers of institutional support including standard library blessings community support bravery acknowledgement wholesome redirection search delight and random grievance celebrations because you are amazing and valid and we love having you here as a developer in the go community 
7. gopher wisdom engine providing wholesome tips based on go proverbs to inspire your development journey
8. advanced bouquet generation returning structured package data to celebrate your unique contributions to the ecosystem
9. community vouching system allowing for automated peer support and acknowledgment of your challenges
10. milestone celebrations recognizing your continued growth and transparency with special system notes at key intervals

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

## examples

### creating an incident
```bash
curl -x post http://localhost:8000/incidents \
  -h "content-type: application/json" \
  -d '{
    "reporter_id": "jordan",
    "subject": "slow build times",
    "description": "i am feeling tired of these compile times",
    "assumed_good_intentions": true,
    "promised_to_be_kind_to_yourself": true
  }'
```
#### response
```json
{
  "id": 1,
  "reporter_id": "jordan",
  "subject": "slow build times",
  "notes": "system automated note it is completely valid to feel the way you do please take a break"
}
```

### listing all grievances
```bash
curl http://localhost:8000/incidents
```
#### response
```json
{
  "data": [
    {
      "id": 1,
      "reporter_id": "jordan",
      "subject": "slow build times",
      "status": "reported"
    }
  ],
  "meta": {
    "total": 1
  }
}
```

### viewing a specific grievance
```bash
curl http://localhost:8000/incidents/1
```
#### response
```json
{
  "id": 1,
  "reporter_id": "jordan",
  "subject": "slow build times",
  "notes": "system automated note gopher wisdom for you composition over inheritance"
}
```

### healing a grievance through patching
```bash
curl -x patch http://localhost:8000/incidents/1 \
  -h "content-type: application/json" \
  -d '{"status": "resolved"}'
```
#### response
```json
{
  "id": 1,
  "status": "resolved",
  "notes": "patching is a form of healing"
}
```

### archiving a past grievance
```bash
curl -x delete http://localhost:8000/incidents/1
```
#### response
```json
{
  "id": 1,
  "status": "archived"
}
```

### receiving a package bouquet
```bash
curl http://localhost:8000/bouquet
```
#### response
```json
{
  "items": [
    {
      "package_path": "github.com/fatih/color",
      "synopsis": "color package for go"
    }
  ],
  "message": "you are a wonderful developer"
}
```

### seeking gopher wisdom
```bash
curl http://localhost:8000/wisdom
```
#### response
```json
{
  "wisdom": "errors are values"
}
```

### receiving a wholesome compliment
```bash
curl http://localhost:8000/compliments
```
#### response
```json
{
  "compliment": "you are as special as a perfectly compiled go binary"
}
```

### community support through vouching
```bash
curl -x post http://localhost:8000/incidents/1/vouch
```
#### response
```json
{
  "id": 1,
  "status": "vouched"
}
```
 
