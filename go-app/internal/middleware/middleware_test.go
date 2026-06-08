package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestCompassionateRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	
	// Limit to 1 request per second with a burst of 1
	r.Use(CompassionateRateLimiter(rate.Limit(1), 1))
	
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// First request - should succeed
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Second request (immediate) - should fail
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusTooManyRequests, w2.Code)
	assert.Contains(t, w2.Body.String(), "overflowing")
}
