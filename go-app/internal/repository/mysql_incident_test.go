package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
	defer func() { _ = db.Close() }()

	repo := NewMySQLIncidentRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "reporter_id", "occurred_at", "recorded_at", "subject", "category", "severity", "description", "evidence_uri", "requires_accommodation", "status", "notes", "was_forgotten", "original_hash"}).
		AddRow(1, "jordan", &now, now, "subject", "cat", 1, "desc", nil, false, "reported", "notes", false, nil)

	mock.ExpectQuery("SELECT id, reporter_id, occurred_at, recorded_at, subject, category, severity, description, evidence_uri, requires_accommodation, status, notes, was_forgotten, original_hash FROM incidents WHERE id = ?").
		WithArgs(uint64(1)).
		WillReturnRows(rows)

	inc, err := repo.GetByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, inc)
	assert.Equal(t, uint64(1), inc.ID)
	assert.Equal(t, "jordan", inc.ReporterID)
	assert.False(t, inc.WasForgotten)
	assert.Nil(t, inc.OriginalHash)
}

func TestMySQLIncidentRepository_Archive(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %s", err)
	}
	defer func() { _ = db.Close() }()

	repo := NewMySQLIncidentRepository(db)

	selectRows := sqlmock.NewRows([]string{"subject", "description"}).
		AddRow("secret subject", "secret description")
	mock.ExpectQuery("SELECT subject, description FROM incidents WHERE id = ?").
		WithArgs(uint64(1)).
		WillReturnRows(selectRows)

	sum := sha256.Sum256([]byte("secret subject:secret description"))
	expectedHash := hex.EncodeToString(sum[:])

	mock.ExpectExec("UPDATE incidents SET status = 'archived', subject = '\\[REDACTED BY REQUEST\\]', description = '\\[REDACTED BY REQUEST\\]', evidence_uri = NULL, notes = NULL, was_forgotten = TRUE, original_hash = \\? WHERE id = \\?").
		WithArgs(expectedHash, uint64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Archive(context.Background(), 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
