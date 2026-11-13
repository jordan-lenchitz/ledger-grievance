# `ledger-grievance`

a go microservice for recording grievances as they occur 

## the twenty two core beliefs of `ledger-grievance`
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
11. opentelemetry integration providing distributed tracing and custom wholesome metrics for deep observability, exposed via prometheus and grafana
12. resiliency patterns including gobreaker circuit breakers and token bucket rate limiting for industrial strength stability
13. deep health check endpoint for comprehensive system readiness and dependency verification
14. advanced makefile overhaul including linting benchmarking security scanning and runtime profiling for operational excellence
15. rock solid fuzz testing (including infinite patience fuzzing of our middleware) and testcontainers driven integration tests, along with enhanced middleware testing to account for anatta (not-self), anicca (impermanence), sunnata (openness) and karuna (compassion) so that our routing may be balanced
16. quine narrator middleware injecting its own punctuation free lowercase source code and meta-narratives into your json responses because a microservice is merely an illusion of boundaries
17. declarative state reconciliation via opentofu and terraform in `tofu/` provisioning kubernetes resources horizontal pod autoscaling mariadb persistence and automated grafana dashboards so our infrastructure drift is gently guided back to harmony
18. strict `accept` literalism requiring `application/vnd.ledger-grievance.v1+json; charset=utf-8` on mutating `/incidents` endpoints and returning `406 not acceptable` with a 12 item `alternates` header for everything else
19. rfc 2324 and rfc 9110 content gatekeeping returning `418` for coffee or burnout complaints and `451` with an nda `link` header for restricted legal or vendor mentions
20. gdpr tombstone poisoning on `delete /incidents/:id` overwriting sensitive fields with `[redacted by request]` while keeping `was_forgotten = true`, the `sha256` `original_hash`, and `gdpr_erasures_total`
21. synthetic `200 ok` sla wrapping via `SLA_GUARANTEE_MODE=true` or `x-sla-mode: 100-percent` moving the real status into `x-real-http-status` and `actual_status_code`
22. nanosecond precondition enforcement rejecting `if-unmodified-since` with `412 precondition failed` and requiring exact `x-nanosecond-epoch` matching against `recorded_at`

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

### seeking a deep health check
```bash
curl http://localhost:8000/health/deep
```
#### response
```json
{
  "overall": true,
  "status": {
    "database": "healthy",
    "pkgsite": "healthy"
  }
}
```

### viewing metrics
```bash
curl http://localhost:8000/metrics
```
#### response
```text
# help go_gc_duration_seconds a summary of the wall-time pause (stop-the-world) duration in garbage collection cycles.
# type go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
go_gc_duration_seconds{quantile="0.75"} 0
go_gc_duration_seconds{quantile="1"} 0
go_gc_duration_seconds_sum 0
go_gc_duration_seconds_count 0
# help go_gc_gogc_percent heap size target percentage configured by the user, otherwise 100. this value is set by the gogc environment variable, and the runtime/debug.setgcpercent function. sourced from /gc/gogc:percent.
# type go_gc_gogc_percent gauge
go_gc_gogc_percent 100
# help go_gc_gomemlimit_bytes go runtime memory limit configured by the user, otherwise math.maxint64. this value is set by the gomemlimit environment variable, and the runtime/debug.setmemorylimit function. sourced from /gc/gomemlimit:bytes.
# type go_gc_gomemlimit_bytes gauge
go_gc_gomemlimit_bytes 9.223372036854776e+18
# help go_goroutines number of goroutines that currently exist.
# type go_goroutines gauge
go_goroutines 6
# help go_info information about the go environment.
# type go_info gauge
go_info{version="go1.26.3"} 1
# help go_memstats_alloc_bytes number of bytes allocated in heap and currently in use. equals to /memory/classes/heap/objects:bytes.
# type go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 225296
# help go_memstats_alloc_bytes_total total number of bytes allocated in heap until now, even if released already. equals to /gc/heap/allocs:bytes.
# type go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 225296
# help go_memstats_buck_hash_sys_bytes number of bytes used by the profiling bucket hash table. equals to /memory/classes/profiling/buckets:bytes.
# type go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 1.444029e+06
# help go_memstats_frees_total total number of heap objects frees. equals to /gc/heap/frees:objects + /gc/heap/tiny/allocs:objects.
# type go_memstats_frees_total counter
go_memstats_frees_total 0
# help go_memstats_gc_sys_bytes number of bytes used for garbage collection system metadata. equals to /memory/classes/metadata/other:bytes.
# type go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 1.429264e+06
# help go_memstats_heap_alloc_bytes number of heap bytes allocated and currently in use, same as go_memstats_alloc_bytes. equals to /memory/classes/heap/objects:bytes.
# type go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 225296
# help go_memstats_heap_idle_bytes number of heap bytes waiting to be used. equals to /memory/classes/heap/released:bytes + /memory/classes/heap/free:bytes.
# type go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 6.119424e+06
# help go_memstats_heap_inuse_bytes number of heap bytes that are in use. equals to /memory/classes/heap/objects:bytes + /memory/classes/heap/unused:bytes
# type go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 1.8432e+06
# help go_memstats_heap_objects number of currently allocated objects. equals to /gc/heap/objects:objects.
# type go_memstats_heap_objects gauge
go_memstats_heap_objects 573
# help go_memstats_heap_released_bytes number of heap bytes released to os. equals to /memory/classes/heap/released:bytes.
# type go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 6.119424e+06
# help go_memstats_heap_sys_bytes number of heap bytes obtained from system. equals to /memory/classes/heap/objects:bytes + /memory/classes/heap/unused:bytes + /memory/classes/heap/released:bytes + /memory/classes/heap/free:bytes.
# type go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 7.962624e+06
# help go_memstats_last_gc_time_seconds number of seconds since 1970 of last garbage collection.
# type go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 0
# help go_memstats_mallocs_total total number of heap objects allocated, both live and gc-ed. semantically a counter version for go_memstats_heap_objects gauge. equals to /gc/heap/allocs:objects + /gc/heap/tiny/allocs:objects.
# type go_memstats_mallocs_total counter
go_memstats_mallocs_total 573
# help go_memstats_mcache_inuse_bytes number of bytes in use by mcache structures. equals to /memory/classes/metadata/mcache/inuse:bytes.
# type go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 9184
# help go_memstats_mcache_sys_bytes number of bytes used for mcache structures obtained from system. equals to /memory/classes/metadata/mcache/inuse:bytes + /memory/classes/metadata/mcache/free:bytes.
# type go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 16072
# help go_memstats_mspan_inuse_bytes number of bytes in use by mspan structures. equals to /memory/classes/metadata/mspan/inuse:bytes.
# type go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 45440
# help go_memstats_mspan_sys_bytes number of bytes used for mspan structures obtained from system. equals to /memory/classes/metadata/mspan/inuse:bytes + /memory/classes/metadata/mspan/free:bytes.
# type go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 48960
# help go_memstats_next_gc_bytes number of heap bytes when next garbage collection will take place. equals to /gc/heap/goal:bytes.
# type go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 4.194304e+06
# help go_memstats_other_sys_bytes number of bytes used for other system allocations. equals to /memory/classes/other:bytes.
# type go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 951091
# help go_memstats_stack_inuse_bytes number of bytes obtained from system for stack allocator in non-cgo environments. equals to /memory/classes/heap/stacks:bytes.
# type go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 425984
# help go_memstats_stack_sys_bytes number of bytes obtained from system for stack allocator. equals to /memory/classes/heap/stacks:bytes + /memory/classes/os-stacks:bytes.
# type go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 425984
# help go_memstats_sys_bytes number of bytes obtained from system. equals to /memory/classes/total:byte.
# type go_memstats_sys_bytes gauge
go_memstats_sys_bytes 1.2278024e+07
# help go_sched_gomaxprocs_threads the current runtime.gomaxprocs setting, or the number of operating system threads that can execute user-level go code simultaneously. sourced from /sched/gomaxprocs:threads.
# type go_sched_gomaxprocs_threads gauge
go_sched_gomaxprocs_threads 4
# help go_threads number of os threads created.
# type go_threads gauge
go_threads 7
# help gopher_wisdom_total total number of gopher wisdom tips shared
# type gopher_wisdom_total counter
gopher_wisdom_total{otel_scope_name="ledger-grievance",otel_scope_schema_url="",otel_scope_version=""} 1
# help process_cpu_seconds_total total user and system cpu time spent in seconds.
# type process_cpu_seconds_total counter
process_cpu_seconds_total 0.04
# help process_max_fds maximum number of open file descriptors.
# type process_max_fds gauge
process_max_fds 1.048575e+06
# help process_network_receive_bytes_total number of bytes received by the process over the network.
# type process_network_receive_bytes_total counter
process_network_receive_bytes_total 1.108652957e+09
# help process_network_transmit_bytes_total number of bytes sent by the process over the network.
# type process_network_transmit_bytes_total counter
process_network_transmit_bytes_total 3.11268052e+08
# help process_open_fds number of open file descriptors.
# type process_open_fds gauge
process_open_fds 6
# help process_resident_memory_bytes resident memory size in bytes.
# type process_resident_memory_bytes gauge
process_resident_memory_bytes 1.2795904e+07
# help process_start_time_seconds start time of the process since unix epoch in seconds.
# type process_start_time_seconds gauge
process_start_time_seconds 1.79426091893e+09
# help process_virtual_memory_bytes virtual memory size in bytes.
# type process_virtual_memory_bytes gauge
process_virtual_memory_bytes 1.688723456e+09
# help process_virtual_memory_max_bytes maximum amount of virtual memory available in bytes.
# type process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes 1.8446744073709552e+19
# help promhttp_metric_handler_requests_in_flight current number of scrapes being served.
# type promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# help promhttp_metric_handler_requests_total total number of scrapes by http status code.
# type promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 1
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
# help target_info target metadata
# type target_info gauge
target_info{service_name="ledger-grievance",service_version="1.0.0"} 1
# help wholesome_notes_total total number of wholesome notes generated
# type wholesome_notes_total counter
wholesome_notes_total{otel_scope_name="ledger-grievance",otel_scope_schema_url="",otel_scope_version=""} 1
```

### creating an incident
```bash
curl -x post http://localhost:8000/incidents \
  -h "content-type: application/json" \
  -h "accept: application/vnd.ledger-grievance.v1+json; charset=utf-8" \
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
  -h "accept: application/vnd.ledger-grievance.v1+json; charset=utf-8" \
  -h "x-nanosecond-epoch: 1700000000123456789" \
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
curl -x delete http://localhost:8000/incidents/1 \
  -h "accept: application/vnd.ledger-grievance.v1+json; charset=utf-8" \
  -h "x-nanosecond-epoch: 1700000000123456789"
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
 
