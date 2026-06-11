package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/domain"
	"go.uber.org/mock/gomock"
)

func FuzzListIncidentsQuery(f *testing.F) {
	// Add some seed corpus
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

		// We don't care about the return value or errors, we just want to ensure it doesn't panic
		mockSvc.EXPECT().ListIncidents(gomock.Any(), gomock.Any()).Return(domain.ListResult{}, nil).AnyTimes()

		h := NewIncidentHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, err := http.NewRequest(http.MethodGet, "/incidents?"+query, nil)
		if err != nil {
			// Fuzzer generated an invalid URL, just skip
			return
		}
		c.Request = req

		// If it panics, the fuzz test fails
		h.ListIncidents(c)
	})
}
