package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/domain"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/middleware"
	"go.uber.org/mock/gomock"
)

func FuzzListIncidentsQuery(f *testing.F) {
	f.Add("limit=10&offset=0")
	f.Add("limit=999&offset=-5")
	f.Add("q=wholesome&status=reported")
	f.Add("limit=abc&offset=def")
	f.Add("limit=0&offset=0")
	f.Add("limit=200&offset=999999999999999999999")

	f.Fuzz(func(t *testing.T, query string) {
		gin.SetMode(gin.TestMode)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSvc := NewMockIncidentService(ctrl)
		mockSvc.EXPECT().ListIncidents(gomock.Any(), gomock.Any()).Return(domain.ListResult{}, nil).AnyTimes()

		h := NewIncidentHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, err := http.NewRequest(http.MethodGet, "/incidents?"+query, nil)
		if err != nil {
			return
		}
		c.Request = req

		h.ListIncidents(c)
	})
}

func FuzzMaliciousComplianceMiddleware(f *testing.F) {
	f.Add("POST", "/incidents", "application/json", "", "", `{"subject":"coffee","description":"burnout"}`)
	f.Add("POST", "/incidents", middleware.RequiredAcceptMediaType, "", "", `{"subject":"ok","description":"fine"}`)
	f.Add("PATCH", "/incidents/1", middleware.RequiredAcceptMediaType, "Mon, 02 Jan 2006 15:04:05 GMT", "1700000000000000000", `{}`)
	f.Add("DELETE", "/incidents/1", "*/*", "", "not-a-num", ``)

	f.Fuzz(func(t *testing.T, method, path, accept, ifUnmod, nanoEpoch, body string) {
		if method != http.MethodPost && method != http.MethodPatch && method != http.MethodDelete && method != http.MethodGet {
			return
		}
		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.Use(middleware.SLAGuarantee(false))
		r.Use(middleware.StrictAccept())
		r.Use(middleware.ContentGatekeeper())
		r.Use(middleware.SubSecondPrecondition(func(_ context.Context, _ uint64) (time.Time, error) {
			return time.Unix(1700000000, 0), nil
		}))
		r.Any("/*path", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		req, err := http.NewRequest(method, path, bytes.NewBufferString(body))
		if err != nil {
			return
		}
		if accept != "" {
			req.Header.Set("Accept", accept)
		}
		if ifUnmod != "" {
			req.Header.Set("If-Unmodified-Since", ifUnmod)
		}
		if nanoEpoch != "" {
			req.Header.Set("X-Nanosecond-Epoch", nanoEpoch)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	})
}
