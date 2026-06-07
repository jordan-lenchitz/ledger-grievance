package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/domain"
	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/middleware"
)

func TestCreateIncident_Handler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockIncidentService(ctrl)
	h := NewIncidentHandler(mockSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/incidents", h.CreateIncident)

	payload := domain.IncidentCreate{
		ReporterID:                 "jordan",
		Subject:                    "test",
		Description:                "test desc",
		Severity:                   1,
		Category:                   "technical",
		AssumedGoodIntentions:      true,
		PromisedToBeKindToYourself: true,
	}
	body, _ := json.Marshal(payload)

	mockSvc.EXPECT().
		CreateIncident(gomock.Any(), gomock.Any()).
		Return(&domain.Incident{ID: 1}, nil)

	req, _ := http.NewRequest(http.MethodPost, "/incidents", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateIncident_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockIncidentService(ctrl)
	h := NewIncidentHandler(mockSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/incidents", h.CreateIncident)

	req, _ := http.NewRequest(http.MethodPost, "/incidents", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateIncident_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockIncidentService(ctrl)
	h := NewIncidentHandler(mockSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler())
	r.POST("/incidents", h.CreateIncident)

	payload := domain.IncidentCreate{
		ReporterID:                 "jordan",
		Subject:                    "test",
		Description:                "test desc",
		Severity:                   1,
		Category:                   "technical",
		AssumedGoodIntentions:      true,
		PromisedToBeKindToYourself: true,
	}
	body, _ := json.Marshal(payload)

	mockSvc.EXPECT().
		CreateIncident(gomock.Any(), gomock.Any()).
		Return(nil, assert.AnError)

	req, _ := http.NewRequest(http.MethodPost, "/incidents", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
