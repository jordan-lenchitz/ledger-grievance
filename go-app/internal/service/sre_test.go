package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/domain"
	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPkgsiteCircuitBreaker(t *testing.T) {
	// We want to test that the circuit breaker trips after enough failures
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "test-cb",
		MaxRequests: 1,
		Interval:    1 * time.Second,
		Timeout:     1 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.TotalFailures >= 2
		},
	})

	action := func() (interface{}, error) {
		return nil, errors.New("fail")
	}

	// First failure
	_, err := cb.Execute(action)
	assert.Error(t, err)
	assert.Equal(t, gobreaker.StateClosed, cb.State())

	// Second failure
	_, err = cb.Execute(action)
	assert.Error(t, err)
	assert.Equal(t, gobreaker.StateOpen, cb.State())

	// Third attempt while open
	_, err = cb.Execute(action)
	assert.Equal(t, gobreaker.ErrOpenState, err)
}

func TestCheckHealth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockIncidentRepository(ctrl)
	mockPkgsite := NewMockPkgsiteService(ctrl)
	svc := NewIncidentService(mockRepo, mockPkgsite)

	ctx := context.Background()

	// Both healthy
	mockRepo.EXPECT().List(gomock.Any(), domain.ListParams{Limit: 1}).Return(domain.ListResult{}, nil)
	mockPkgsite.EXPECT().CheckHealth(gomock.Any()).Return(nil)

	status := svc.CheckHealth(ctx)
	assert.Equal(t, "healthy", status["database"])
	assert.Equal(t, "healthy", status["pkgsite"])

	// Both unhealthy
	mockRepo.EXPECT().List(gomock.Any(), domain.ListParams{Limit: 1}).Return(domain.ListResult{}, errors.New("db down"))
	mockPkgsite.EXPECT().CheckHealth(gomock.Any()).Return(errors.New("api down"))

	status = svc.CheckHealth(ctx)
	assert.Equal(t, "unhealthy: db down", status["database"])
	assert.Equal(t, "unhealthy: api down", status["pkgsite"])
}
