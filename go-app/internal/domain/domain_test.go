package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncidentJSONRoundTrip(t *testing.T) {
	notes := "investigating"
	orig := Incident{
		ID:          1,
		ReporterID:  "jordan",
		Subject:     "test",
		Category:    "technical",
		Severity:    3,
		Description: "details",
		Status:      StatusReported,
		Notes:       &notes,
	}

	raw, err := json.Marshal(orig)
	require.NoError(t, err)

	var decoded Incident
	require.NoError(t, json.Unmarshal(raw, &decoded))

	assert.Equal(t, orig.ID, decoded.ID)
	assert.Equal(t, orig.ReporterID, decoded.ReporterID)
	assert.Equal(t, orig.Status, decoded.Status)
	require.NotNil(t, decoded.Notes)
	assert.Equal(t, notes, *decoded.Notes)
}
