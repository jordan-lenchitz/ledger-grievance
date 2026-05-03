# Grievance Ledger POC

A tiny MariaDB + FastAPI proof of concept for a neutral, user-reported incident journal.

Important framing: records are reports, not verified claims. This is for personal tracking, review, and structured recall, not harassment, doxxing, retaliation, or automated collection.

## Run

```bash
docker compose up --build
```

API docs:

```text
http://localhost:8000/docs
```

Health check:

```bash
curl http://localhost:8000/health
```

## Create an incident

```bash
curl -X POST http://localhost:8000/incidents \
  -H 'Content-Type: application/json' \
  -d '{
    "reporter_id": "jordan",
    "occurred_at": "2026-05-03T12:30:00",
    "subject": "unnamed party",
    "category": "communication",
    "severity": 2,
    "description": "Reported incident summary goes here, written neutrally.",
    "evidence_uri": "https://example.com/reference",
    "notes": "Initial note. Needs review."
  }'
```

## List incidents

```bash
curl 'http://localhost:8000/incidents?reporter_id=jordan&limit=25'
```

Search:

```bash
curl 'http://localhost:8000/incidents?q=communication'
```

## Update incident status or notes

```bash
curl -X PATCH http://localhost:8000/incidents/1 \
  -H 'Content-Type: application/json' \
  -d '{"status":"reviewing","notes":"Added context after review."}'
```

## Archive an incident

```bash
curl -X DELETE http://localhost:8000/incidents/1
```

## Schema

Core table: `incidents`

- `id`
- `reporter_id`
- `occurred_at`
- `recorded_at`
- `subject`
- `category`
- `severity` from 1 to 5
- `description`
- `evidence_uri`
- `status`: reported, reviewing, resolved, dismissed, archived
- `notes`

## POC boundaries

- No scraping
- No automatic collection
- No private-data ingestion
- No claims of verification
- No publishing or sharing workflow
