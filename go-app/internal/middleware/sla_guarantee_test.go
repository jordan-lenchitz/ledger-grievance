package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSLAGuarantee(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		alwaysWrap     bool
		slaHeader      string
		handlerStatus  int
		handlerJSON    gin.H
		handlerRaw     string
		wantWireStatus int
		wantRealHeader string
		wantBodyStatus int
	}{
		{
			name:           "disabled without alwaysWrap or header",
			alwaysWrap:     false,
			handlerStatus:  http.StatusNotFound,
			handlerJSON:    gin.H{"error": "missing"},
			wantWireStatus: http.StatusNotFound,
		},
		{
			name:           "wraps 404 via header into 200 OK",
			alwaysWrap:     false,
			slaHeader:      "100-percent",
			handlerStatus:  http.StatusNotFound,
			handlerJSON:    gin.H{"error": "missing"},
			wantWireStatus: http.StatusOK,
			wantRealHeader: "404",
			wantBodyStatus: http.StatusNotFound,
		},
		{
			name:           "wraps 418 via alwaysWrap into 200 OK",
			alwaysWrap:     true,
			handlerStatus:  http.StatusTeapot,
			handlerJSON:    gin.H{"error": "teapot"},
			wantWireStatus: http.StatusOK,
			wantRealHeader: "418",
			wantBodyStatus: http.StatusTeapot,
		},
		{
			name:           "wraps non-JSON 503 into JSON 200 OK",
			alwaysWrap:     true,
			handlerStatus:  http.StatusServiceUnavailable,
			handlerRaw:     "upstream down",
			wantWireStatus: http.StatusOK,
			wantRealHeader: "503",
			wantBodyStatus: http.StatusServiceUnavailable,
		},
		{
			name:           "passes 200 OK through without X-Real-HTTP-Status",
			alwaysWrap:     true,
			handlerStatus:  http.StatusOK,
			handlerJSON:    gin.H{"status": "healthy"},
			wantWireStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			r.Use(SLAGuarantee(tc.alwaysWrap))
			r.GET("/test", func(c *gin.Context) {
				if tc.handlerRaw != "" {
					c.String(tc.handlerStatus, tc.handlerRaw)
					return
				}
				c.JSON(tc.handlerStatus, tc.handlerJSON)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tc.slaHeader != "" {
				req.Header.Set("X-SLA-Mode", tc.slaHeader)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.wantWireStatus, w.Code)
			assert.Equal(t, tc.wantRealHeader, w.Header().Get("X-Real-HTTP-Status"))

			if tc.wantBodyStatus != 0 {
				var body map[string]any
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
				assert.Equal(t, float64(tc.wantBodyStatus), body["actual_status_code"])
			}
		})
	}
}

func TestSLAGuarantee_WithQuineNarrator(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(SLAGuarantee(true))
	r.Use(QuineNarrator())
	r.GET("/fail", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "boom"})
	})

	req := httptest.NewRequest(http.MethodGet, "/fail", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "500", w.Header().Get("X-Real-HTTP-Status"))

	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, float64(500), body["actual_status_code"])
	assert.Contains(t, body, "narrative")
	assert.Contains(t, body, "essence")
}
