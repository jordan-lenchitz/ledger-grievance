package repository

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqlmigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mariadb"
)

func TestMySQLIncidentRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	migrationsDir, err := filepath.Abs("../../migrations")
	require.NoError(t, err)

	mariadbContainer, err := mariadb.Run(ctx,
		"mariadb:11.4",
		mariadb.WithDatabase("grievance_ledger"),
		mariadb.WithUsername("ledger"),
		mariadb.WithPassword("ledgerpass"),
	)
	require.NoError(t, err)

	defer func() {
		if err := mariadbContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	connStr, err := mariadbContainer.ConnectionString(ctx)
	require.NoError(t, err)

	if !strings.Contains(connStr, "?") {
		connStr += "?parseTime=true&multiStatements=true"
	} else {
		connStr += "&parseTime=true&multiStatements=true"
	}

	db, err := sql.Open("mysql", connStr)
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	driver, err := mysqlmigrate.WithInstance(db, &mysqlmigrate.Config{})
	require.NoError(t, err)

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsDir, "mysql", driver)
	require.NoError(t, err)
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		require.NoError(t, err)
	}

	repo := NewMySQLIncidentRepository(db)

	t.Run("Create and Get Incident", func(t *testing.T) {
		occurredAt := time.Now()
		notes := "test notes"
		evidenceURI := "http://example.com"
		inc := &domain.Incident{
			ReporterID:  "jordan",
			Subject:     "Integration Test Grievance",
			Category:    "testing",
			Severity:    3,
			Description: "This is a test from Testcontainers",
			OccurredAt:  &occurredAt,
			Notes:       &notes,
			EvidenceURI: &evidenceURI,
		}

		id, err := repo.Create(ctx, inc)
		require.NoError(t, err)
		assert.Greater(t, id, uint64(0))

		fetched, err := repo.GetByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, inc.ReporterID, fetched.ReporterID)
		assert.Equal(t, inc.Subject, fetched.Subject)
		assert.Equal(t, inc.Category, fetched.Category)
		assert.Equal(t, inc.Severity, fetched.Severity)
		assert.Equal(t, inc.Description, fetched.Description)
		assert.False(t, fetched.WasForgotten)
		assert.Nil(t, fetched.OriginalHash)

		err = repo.Archive(ctx, id)
		require.NoError(t, err)

		forgotten, err := repo.GetByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, domain.StatusArchived, forgotten.Status)
		assert.Equal(t, "[REDACTED BY REQUEST]", forgotten.Subject)
		assert.Equal(t, "[REDACTED BY REQUEST]", forgotten.Description)
		assert.Nil(t, forgotten.EvidenceURI)
		assert.Nil(t, forgotten.Notes)
		assert.True(t, forgotten.WasForgotten)
		require.NotNil(t, forgotten.OriginalHash)
		assert.Len(t, *forgotten.OriginalHash, 64)
	})
}
