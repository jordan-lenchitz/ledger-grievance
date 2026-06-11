package repository

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

	dbScriptPath, err := filepath.Abs("../../../db/init.sql")
	require.NoError(t, err)

	mariadbContainer, err := mariadb.Run(ctx,
		"mariadb:11.4",
		mariadb.WithDatabase("grievance_ledger"),
		mariadb.WithUsername("ledger"),
		mariadb.WithPassword("ledgerpass"),
		mariadb.WithScripts(dbScriptPath),
	)
	require.NoError(t, err)

	// Clean up the container
	defer func() {
		if err := mariadbContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	connStr, err := mariadbContainer.ConnectionString(ctx)
	require.NoError(t, err)

	// Append parseTime=true to the connection string to handle time.Time correctly
	if !strings.Contains(connStr, "?") {
		connStr += "?parseTime=true"
	} else {
		connStr += "&parseTime=true"
	}

	db, err := sql.Open("mysql", connStr)
	require.NoError(t, err)
	defer db.Close()

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
	})
}
