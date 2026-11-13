# ledger-grievance

an institutional-friction-enabled grievance journal
> see `go-app/readme.md` for full technical and architectural details

## core concepts
- **go api** for easy incident management
- **native cli** for interacting with the ledger
- **institutional metaphors** implement specific wholesome support features reflecting the best of humanity
    1. **the kindness promise** means you must promise to be kind to yourself or we will gently remind you that your well-being is our top priority
    2. **wholesome validation** providing automated notes triggered by pre-defined strings in your grievance to remind you that your feelings are valid and you should take a break
    3. **assume good intent requirement** because we believe in your capacity for empathy and kindness
    4. **go api integration** implements seven amazing ways to support your grievance using [pkg.go.dev api](https://pkg.go.dev/v1beta/api) data like standard library blessings and community support and bravery acknowledgement and wholesome redirection and search delight gaslighting turned to delight
    5. **grievance celebrations** where we randomly upgrade your grievance to a celebration because you are doing a great job
    6. **compliment endpoint** giving you random wholesome compliments based on the go package ecosystem because you deserve it
    7. **gopher wisdom engine** providing you with wholesome go themed tips to brighten your day
    8. **advanced bouquet generation** giving you a structured selection of amazing packages from across the ecosystem
    9. **community vouching system** where others can stand in solidarity with your grievance through automated gopher support
    10. **milestone celebrations** honoring your commitment to emotional transparency with special badges at key intervals
    11. **wholesome observability** with opentelemetry tracing and metrics to visualize the journey of your grievance
    12. **gopher guard reliability** featuring circuit breakers and rate limiting to protect your peace of mind
    13. **deep emotional health checks** verifying the system is fully ready to receive your kindness
    14. **auspicious infrastructure as code** providing a robust makefile for lifecycle management and security scanning
    15. **quine narrator middleware** injecting its own punctuation free lowercase source code and meta-narratives into your json responses because a microservice is merely an illusion of boundaries
    16. **declarative state reconciliation and scalable kindness** via opentofu and terraform configurations in `tofu/` provisioning kubernetes resources horizontal pod autoscaling mariadb persistence prometheus scraping and wholesome grafana dashboards because true healing begins when desired state meets actual state
    17. **strict accept literalism** rejecting wildcard and generic `application/json` headers on mutating endpoints with `406 not acceptable` and a 12 item `alternates` header listing every media type we refuse to negotiate
    18. **rfc 2324 and 451 content gatekeeper** returning `418 i'm a teapot` for coffee or burnout grievances and `451 unavailable for legal reasons` with an nda `link` header when restricted vendor or named party terms appear
    19. **gdpr tombstone poisoning** honoring right to be forgotten requests on `delete /incidents/:id` by redacting subject and description while preserving the row `was_forgotten` flag `sha256` hash of the original complaint and a `gdpr_erasures_total` counter
    20. **100 percent uptime sla guarantee** optionally rewriting every non `200` response into a synthetic `200 ok` with `x-real-http-status` and `actual_status_code` tucked into the json body so uptime monitors stay green
    21. **sub second precondition pedantry** rejecting 1 second resolution `if-unmodified-since` headers with `412 precondition failed` and requiring nanosecond exact `x-nanosecond-epoch` timestamps before mutating history

## robust testing
our system is verified by a massive suite of table driven tests that demonstrate our commitment to bureaucratic excellence by mocking the go pkg api with httptest and validating every single wholesome note and status transition to ensure that your happiness is mathematically guaranteed and consistently delivered across all seven layers of institutional support. we also utilize **fuzz testing** to ensure no edge cases panic our systems (including infinite patience fuzzing of our middleware), and **testcontainers** for fully isolated integration testing against live mariadb instances! we have also enhanced our middleware testing to account for anatta (not-self), anicca (impermanence), sunnata (openness) and karuna (compassion) so that our routing may be balanced

## how to get started
1. **api**: see `go-app/readme.md` for setup, migrations, and ci/cd details.
2. **cli**:
   ```bash
   cd cli
   go run main.go --help
   ```
3. **containers** (includes api, mariadb, prometheus, and grafana!):
   ```bash
   docker compose up --build
   ```
   > visit `http://localhost:3000` to view the wholesome grafana observability dashboards!
