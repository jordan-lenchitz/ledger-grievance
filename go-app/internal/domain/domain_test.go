package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIncidentStruct(t *testing.T) {
	// A simple test to ensure the domain structs can be instantiated
	inc := Incident{
		ID:         1,
		ReporterID: "jordan",
		Subject:    "test",
		Status:     StatusReported,
	}
	assert.Equal(t, uint64(1), inc.ID)
	assert.Equal(t, StatusReported, inc.Status)
}
