package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestPkgsiteService_GetPackage(t *testing.T) {
	// Table-driven tests for PkgsiteService
	tests := []struct {
		name           string
		path           string
		serverResponse any
		status         int
		expectedPkg    *domain.Package
		expectedErr    bool
	}{
		{
			name: "Success - Standard Library",
			path: "fmt",
			serverResponse: domain.Package{
				Path:              "fmt",
				IsStandardLibrary: true,
				Synopsis:          "Package fmt implements formatted I/O.",
			},
			status: http.StatusOK,
			expectedPkg: &domain.Package{
				Path:              "fmt",
				IsStandardLibrary: true,
				Synopsis:          "Package fmt implements formatted I/O.",
			},
			expectedErr: false,
		},
		{
			name:           "Error - Not Found",
			path:           "nonexistent",
			serverResponse: nil,
			status:         http.StatusNotFound,
			expectedPkg:    nil,
			expectedErr:    true,
		},
		{
			name:           "Error - Internal Server Error",
			path:           "broken",
			serverResponse: nil,
			status:         http.StatusInternalServerError,
			expectedPkg:    nil,
			expectedErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				if tt.serverResponse != nil {
					json.NewEncoder(w).Encode(tt.serverResponse)
				}
			}))
			defer server.Close()

			svc := NewPkgsiteService(server.URL)
			pkg, err := svc.GetPackage(context.Background(), tt.path)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPkg, pkg)
			}
		})
	}
}

func TestPkgsiteService_Search(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/search", r.URL.Path)
		assert.Equal(t, "query", r.URL.Query().Get("q"))

		resp := domain.PaginatedResponse[domain.SearchResult]{
			Items: []domain.SearchResult{
				{PackagePath: "github.com/test/pkg", Synopsis: "a test package"},
			},
			Total: 1,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	svc := NewPkgsiteService(server.URL)
	res, err := svc.Search(context.Background(), "query")

	assert.NoError(t, err)
	assert.Equal(t, 1, res.Total)
	assert.Equal(t, "github.com/test/pkg", res.Items[0].PackagePath)
}
