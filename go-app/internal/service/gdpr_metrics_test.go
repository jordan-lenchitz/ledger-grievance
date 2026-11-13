package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/telemetry"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestArchiveIncident_RecordsGDPRErasureMetric(t *testing.T) {
	ctx := context.Background()
	shutdown, err := telemetry.SetupOTEL(ctx)
	require.NoError(t, err)
	defer func() {
		_ = shutdown(ctx)
	}()

	// Re-bind counter against the newly configured global MeterProvider.
	gdprErasuresCounter, err = telemetry.Meter.Int64Counter("gdpr_erasures_total")
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockIncidentRepository(ctrl)
	mockRepo.EXPECT().Archive(gomock.Any(), uint64(42)).Return(nil)

	svc := NewIncidentService(mockRepo, nil)
	require.NoError(t, svc.ArchiveIncident(ctx, 42))

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	promhttp.Handler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "gdpr_erasures_total")
}
