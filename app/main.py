import os
from datetime import datetime
from typing import Literal, Optional

import pymysql
from fastapi import FastAPI, HTTPException, Query
from pydantic import BaseModel, Field

Status = Literal["reported", "reviewing", "resolved", "dismissed", "archived"]

app = FastAPI(
    title="Grievance Ledger POC",
    description="A neutral, user-reported incident journal. Entries are reports, not verified claims.",
    version="0.1.0",
)


def db():
    return pymysql.connect(
        host=os.getenv("DB_HOST", "localhost"),
        port=int(os.getenv("DB_PORT", "3306")),
        user=os.getenv("DB_USER", "ledger"),
        password=os.getenv("DB_PASSWORD", "ledgerpass"),
        database=os.getenv("DB_NAME", "grievance_ledger"),
        cursorclass=pymysql.cursors.DictCursor,
        autocommit=True,
    )


class IncidentCreate(BaseModel):
    reporter_id: str = Field(..., min_length=1, max_length=128)
    occurred_at: Optional[datetime] = None
    subject: str = Field(..., min_length=1, max_length=255)
    category: str = Field(default="unspecified", min_length=1, max_length=128)
    severity: int = Field(default=1, ge=1, le=5)
    description: str = Field(..., min_length=1)
    evidence_uri: Optional[str] = None
    notes: Optional[str] = None


class IncidentPatch(BaseModel):
    status: Optional[Status] = None
    notes: Optional[str] = None
    category: Optional[str] = Field(default=None, min_length=1, max_length=128)
    severity: Optional[int] = Field(default=None, ge=1, le=5)
    evidence_uri: Optional[str] = None


@app.get("/health")
def health():
    with db() as conn:
        with conn.cursor() as cur:
            cur.execute("SELECT 1 AS ok")
            return cur.fetchone()


@app.post("/incidents", status_code=201)
def create_incident(payload: IncidentCreate):
    sql = """
    INSERT INTO incidents
      (reporter_id, occurred_at, subject, category, severity, description, evidence_uri, notes)
    VALUES
      (%s, %s, %s, %s, %s, %s, %s, %s)
    """
    values = (
        payload.reporter_id,
        payload.occurred_at,
        payload.subject,
        payload.category,
        payload.severity,
        payload.description,
        payload.evidence_uri,
        payload.notes,
    )
    with db() as conn:
        with conn.cursor() as cur:
            cur.execute(sql, values)
            incident_id = cur.lastrowid
    return get_incident(incident_id)


@app.get("/incidents")
def list_incidents(
    reporter_id: Optional[str] = None,
    status: Optional[Status] = None,
    category: Optional[str] = None,
    q: Optional[str] = Query(default=None, description="Search subject, description, and notes"),
    limit: int = Query(default=50, ge=1, le=200),
    offset: int = Query(default=0, ge=0),
):
    clauses = []
    values = []
    if reporter_id:
        clauses.append("reporter_id = %s")
        values.append(reporter_id)
    if status:
        clauses.append("status = %s")
        values.append(status)
    if category:
        clauses.append("category = %s")
        values.append(category)
    if q:
        clauses.append("(subject LIKE %s OR description LIKE %s OR notes LIKE %s)")
        like = f"%{q}%"
        values.extend([like, like, like])

    where = " WHERE " + " AND ".join(clauses) if clauses else ""
    sql = f"""
    SELECT * FROM incidents
    {where}
    ORDER BY recorded_at DESC, id DESC
    LIMIT %s OFFSET %s
    """
    values.extend([limit, offset])
    with db() as conn:
        with conn.cursor() as cur:
            cur.execute(sql, values)
            return cur.fetchall()


@app.get("/incidents/{incident_id}")
def get_incident(incident_id: int):
    with db() as conn:
        with conn.cursor() as cur:
            cur.execute("SELECT * FROM incidents WHERE id = %s", (incident_id,))
            row = cur.fetchone()
    if not row:
        raise HTTPException(status_code=404, detail="Incident not found")
    return row


@app.patch("/incidents/{incident_id}")
def patch_incident(incident_id: int, payload: IncidentPatch):
    updates = payload.model_dump(exclude_unset=True)
    if not updates:
        return get_incident(incident_id)

    assignments = ", ".join(f"{key} = %s" for key in updates.keys())
    values = list(updates.values()) + [incident_id]
    with db() as conn:
        with conn.cursor() as cur:
            cur.execute(f"UPDATE incidents SET {assignments} WHERE id = %s", values)
            if cur.rowcount == 0:
                raise HTTPException(status_code=404, detail="Incident not found")
    return get_incident(incident_id)


@app.delete("/incidents/{incident_id}")
def archive_incident(incident_id: int):
    with db() as conn:
        with conn.cursor() as cur:
            cur.execute("UPDATE incidents SET status = 'archived' WHERE id = %s", (incident_id,))
            if cur.rowcount == 0:
                raise HTTPException(status_code=404, detail="Incident not found")
    return {"id": incident_id, "status": "archived"}
