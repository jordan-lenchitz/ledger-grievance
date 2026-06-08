package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestMySQLIncidentRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %s", err)
	}
	defer db.Close()

	repo := NewMySQLIncidentRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "reporter_id", "occurred_at", "recorded_at", "subject", "category", "severity", "description", "evidence_uri", "requires_accommodation", "status", "notes"}).
		AddRow(1, "jordan", &now, now, "subject", "cat", 1, "desc", nil, false, "reported", "notes")

	mock.ExpectQuery("SELECT id, reporter_id, occurred_at, recorded_at, subject, category, severity, description, evidence_uri, requires_accommodation, status, notes FROM incidents WHERE id = ?").
		WithArgs(uint64(1)).
		WillReturnRows(rows)

	inc, err := repo.GetByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, inc)
	assert.Equal(t, uint64(1), inc.ID)
	assert.Equal(t, "jordan", inc.ReporterID)
}

func TestMySQLIncidentRepository_Archive(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %s", err)
	}
	defer db.Close()

	repo := NewMySQLIncidentRepository(db)

	mock.ExpectExec("UPDATE incidents SET status = 'archived' WHERE id = ?").
		WithArgs(uint64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Archive(context.Background(), 1)
	assert.NoError(t, err)
}
