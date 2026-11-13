package middleware

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStrictAccept(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		path           string
		accept         string
		wantStatus     int
		wantAlternates bool
	}{
		{
			name:           "rejects wildcard accept on POST /incidents",
			method:         http.MethodPost,
			path:           "/incidents",
			accept:         "*/*",
			wantStatus:     http.StatusNotAcceptable,
			wantAlternates: true,
		},
		{
			name:           "rejects generic application/json on PATCH",
			method:         http.MethodPatch,
			path:           "/incidents/1",
			accept:         "application/json",
			wantStatus:     http.StatusNotAcceptable,
			wantAlternates: true,
		},
		{
			name:           "rejects missing accept on DELETE",
			method:         http.MethodDelete,
			path:           "/incidents/1",
			accept:         "",
			wantStatus:     http.StatusNotAcceptable,
			wantAlternates: true,
		},
		{
			name:           "rejects missing charset=utf-8",
			method:         http.MethodPost,
			path:           "/incidents",
			accept:         "application/vnd.ledger-grievance.v1+json",
			wantStatus:     http.StatusNotAcceptable,
			wantAlternates: true,
		},
		{
			name:       "allows exact vendor media type with charset",
			method:     http.MethodPost,
			path:       "/incidents",
			accept:     RequiredAcceptMediaType,
			wantStatus: http.StatusOK,
		},
		{
			name:       "exempts GET /incidents",
			method:     http.MethodGet,
			path:       "/incidents",
			accept:     "*/*",
			wantStatus: http.StatusOK,
		},
		{
			name:       "exempts non-incident routes",
			method:     http.MethodPost,
			path:       "/health",
			accept:     "application/json",
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			r.Use(StrictAccept())
			r.Any("/*path", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(tc.method, tc.path, nil)
			if tc.accept != "" {
				req.Header.Set("Accept", tc.accept)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			if tc.wantAlternates {
				alt := w.Header().Get("Alternates")
				assert.NotEmpty(t, alt)
				assert.Equal(t, 12, len(strings.Split(alt, ", ")))
			}
		})
	}
}

func TestContentGatekeeper(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		subject    string
		desc       string
		wantStatus int
		wantLink   string
	}{
		{
			name:       "coffee triggers 418",
			subject:    "Need coffee",
			desc:       "The breakroom is out of beans",
			wantStatus: http.StatusTeapot,
		},
		{
			name:       "burnout triggers 418",
			subject:    "Sprint",
			desc:       "Experiencing severe burnout",
			wantStatus: http.StatusTeapot,
		},
		{
			name:       "NDA triggers 451 with Link header",
			subject:    "Vendor issue",
			desc:       "Blocked by NDA",
			wantStatus: http.StatusUnavailableForLegalReasons,
			wantLink:   `<https://example.com/nda>; rel="blocked-by"`,
		},
		{
			name:       "at-mention triggers 451",
			subject:    "Escalation",
			desc:       "Paging @bob about this",
			wantStatus: http.StatusUnavailableForLegalReasons,
			wantLink:   `<https://example.com/nda>; rel="blocked-by"`,
		},
		{
			name:       "clean payload passes through intact",
			subject:    "CI pipeline flaked",
			desc:       "Unit test timed out waiting on lock",
			wantStatus: http.StatusCreated,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			r.Use(ContentGatekeeper())
			r.POST("/incidents", func(c *gin.Context) {
				var got map[string]any
				err := c.ShouldBindJSON(&got)
				require.NoError(t, err)
				c.JSON(http.StatusCreated, got)
			})

			body, err := json.Marshal(map[string]string{
				"subject":     tc.subject,
				"description": tc.desc,
			})
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/incidents", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			if tc.wantLink != "" {
				assert.Equal(t, tc.wantLink, w.Header().Get("Link"))
			}
		})
	}
}

func TestSubSecondPrecondition(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorded := time.Unix(1700000000, 123456789).UTC()

	lookup := func(_ context.Context, id uint64) (time.Time, error) {
		if id == 999 {
			return time.Time{}, sql.ErrNoRows
		}
		return recorded, nil
	}

	tests := []struct {
		name              string
		method            string
		path              string
		ifUnmodifiedSince string
		nanoEpoch         string
		wantStatus        int
	}{
		{
			name:              "rejects RFC 1123 If-Unmodified-Since",
			method:            http.MethodPatch,
			path:              "/incidents/1",
			ifUnmodifiedSince: recorded.Format(http.TimeFormat),
			nanoEpoch:         strconv.FormatInt(recorded.UnixNano(), 10),
			wantStatus:        http.StatusPreconditionFailed,
		},
		{
			name:       "rejects missing X-Nanosecond-Epoch",
			method:     http.MethodPatch,
			path:       "/incidents/1",
			wantStatus: http.StatusPreconditionFailed,
		},
		{
			name:       "rejects non-numeric X-Nanosecond-Epoch",
			method:     http.MethodDelete,
			path:       "/incidents/1",
			nanoEpoch:  "not-a-number",
			wantStatus: http.StatusPreconditionFailed,
		},
		{
			name:       "rejects off-by-one nanosecond",
			method:     http.MethodPatch,
			path:       "/incidents/1",
			nanoEpoch:  strconv.FormatInt(recorded.UnixNano()-1, 10),
			wantStatus: http.StatusPreconditionFailed,
		},
		{
			name:       "returns 404 when incident does not exist",
			method:     http.MethodDelete,
			path:       "/incidents/999",
			nanoEpoch:  strconv.FormatInt(recorded.UnixNano(), 10),
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "allows exact nanosecond match on PATCH",
			method:     http.MethodPatch,
			path:       "/incidents/1",
			nanoEpoch:  strconv.FormatInt(recorded.UnixNano(), 10),
			wantStatus: http.StatusOK,
		},
		{
			name:       "allows exact nanosecond match on DELETE",
			method:     http.MethodDelete,
			path:       "/incidents/1",
			nanoEpoch:  strconv.FormatInt(recorded.UnixNano(), 10),
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			r.Use(SubSecondPrecondition(lookup))
			r.PATCH("/incidents/:id", func(c *gin.Context) { c.Status(http.StatusOK) })
			r.DELETE("/incidents/:id", func(c *gin.Context) { c.Status(http.StatusOK) })

			req := httptest.NewRequest(tc.method, tc.path, nil)
			if tc.ifUnmodifiedSince != "" {
				req.Header.Set("If-Unmodified-Since", tc.ifUnmodifiedSince)
			}
			if tc.nanoEpoch != "" {
				req.Header.Set("X-Nanosecond-Epoch", tc.nanoEpoch)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
		})
	}
}
