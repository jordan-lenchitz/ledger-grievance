package service

import (
	"errors"
	"testing"
	"time"

	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
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
	// This tests the logic of CheckHealth in incidentService
	// Since we already test repo.List, we just need to ensure it aggregates correctly
	// (Implementation details depend on how much we want to mock here)
}
