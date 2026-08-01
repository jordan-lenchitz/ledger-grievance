package service

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/domain"
	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPkgsiteCircuitBreaker(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	svc := NewPkgsiteService(server.URL)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		_, err := svc.GetPackage(ctx, "fmt")
		assert.Error(t, err)
		assert.NotErrorIs(t, err, gobreaker.ErrOpenState)
	}

	_, err := svc.GetPackage(ctx, "fmt")
	assert.ErrorIs(t, err, gobreaker.ErrOpenState)
}

func TestCheckHealth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockIncidentRepository(ctrl)
	mockPkgsite := NewMockPkgsiteService(ctrl)
	svc := NewIncidentService(mockRepo, mockPkgsite)

	ctx := context.Background()

	mockRepo.EXPECT().List(gomock.Any(), domain.ListParams{Limit: 1}).Return(domain.ListResult{}, nil)
	mockPkgsite.EXPECT().CheckHealth(gomock.Any()).Return(nil)

	status := svc.CheckHealth(ctx)
	assert.Equal(t, "healthy", status["database"])
	assert.Equal(t, "healthy", status["pkgsite"])

	mockRepo.EXPECT().List(gomock.Any(), domain.ListParams{Limit: 1}).Return(domain.ListResult{}, errors.New("db down"))
	mockPkgsite.EXPECT().CheckHealth(gomock.Any()).Return(errors.New("api down"))

	status = svc.CheckHealth(ctx)
	assert.Equal(t, "unhealthy: db down", status["database"])
	assert.Equal(t, "unhealthy: api down", status["pkgsite"])
}
