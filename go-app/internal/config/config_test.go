package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	os.Clearenv()
	cfg := Load()
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "ledger", cfg.DBUser)
	assert.Equal(t, "8000", cfg.Port)

	t.Setenv("DB_HOST", "db.example.com")
	t.Setenv("DB_USER", "admin")
	t.Setenv("PORT", "9000")

	cfg = Load()
	assert.Equal(t, "db.example.com", cfg.DBHost)
	assert.Equal(t, "admin", cfg.DBUser)
	assert.Equal(t, "9000", cfg.Port)
}
