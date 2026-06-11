package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Test default values
	os.Clearenv()
	cfg := Load()
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "ledger", cfg.DBUser)
	assert.Equal(t, "8000", cfg.Port)

	// Test environment variable overrides
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_USER", "admin")
	os.Setenv("PORT", "9000")

	cfg = Load()
	assert.Equal(t, "db.example.com", cfg.DBHost)
	assert.Equal(t, "admin", cfg.DBUser)
	assert.Equal(t, "9000", cfg.Port)
}
