package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/domain"
	"github.com/sony/gobreaker"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type PkgsiteService interface {
	GetPackage(ctx context.Context, path string) (*domain.Package, error)
	GetModule(ctx context.Context, path string) (*domain.Module, error)
	GetVersions(ctx context.Context, path string) (*domain.PaginatedResponse[domain.ModuleVersion], error)
	GetSymbols(ctx context.Context, path string) (*domain.PackageSymbols, error)
	GetImportedBy(ctx context.Context, path string) (*domain.PackageImportedBy, error)
	GetVulns(ctx context.Context, path string) (*domain.PaginatedResponse[domain.Vulnerability], error)
	Search(ctx context.Context, query string) (*domain.PaginatedResponse[domain.SearchResult], error)
	CheckHealth(ctx context.Context) error
}

type pkgsiteService struct {
	baseURL string
	client  *http.Client
	cb      *gobreaker.CircuitBreaker
}

func NewPkgsiteService(baseURL string) PkgsiteService {
	if baseURL == "" {
		baseURL = "https://pkg.go.dev/v1beta"
	}

	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "pkgsite-api",
		MaxRequests: 3,
		Interval:    10 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 5 && failureRatio >= 0.6
		},
	})

	return &pkgsiteService{
		baseURL: baseURL,
		client: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
			Timeout:   10 * time.Second,
		},
		cb: cb,
	}
}

func (s *pkgsiteService) get(ctx context.Context, endpoint string, result any) error {
	_, err := s.cb.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, "GET", s.baseURL+endpoint, nil)
		if err != nil {
			return nil, err
		}

		resp, err := s.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("not found")
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("api error: %d", resp.StatusCode)
		}

		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (s *pkgsiteService) GetPackage(ctx context.Context, path string) (*domain.Package, error) {
	var pkg domain.Package
	err := s.get(ctx, "/package/"+path, &pkg)
	return &pkg, err
}

func (s *pkgsiteService) GetModule(ctx context.Context, path string) (*domain.Module, error) {
	var mod domain.Module
	err := s.get(ctx, "/module/"+path, &mod)
	return &mod, err
}

func (s *pkgsiteService) GetVersions(ctx context.Context, path string) (*domain.PaginatedResponse[domain.ModuleVersion], error) {
	var versions domain.PaginatedResponse[domain.ModuleVersion]
	err := s.get(ctx, "/versions/"+path, &versions)
	return &versions, err
}

func (s *pkgsiteService) GetSymbols(ctx context.Context, path string) (*domain.PackageSymbols, error) {
	var symbols domain.PackageSymbols
	err := s.get(ctx, "/symbols/"+path, &symbols)
	return &symbols, err
}

func (s *pkgsiteService) GetImportedBy(ctx context.Context, path string) (*domain.PackageImportedBy, error) {
	var importedBy domain.PackageImportedBy
	err := s.get(ctx, "/imported-by/"+path, &importedBy)
	return &importedBy, err
}

func (s *pkgsiteService) GetVulns(ctx context.Context, path string) (*domain.PaginatedResponse[domain.Vulnerability], error) {
	var vulns domain.PaginatedResponse[domain.Vulnerability]
	err := s.get(ctx, "/vulns/"+path, &vulns)
	return &vulns, err
}

func (s *pkgsiteService) Search(ctx context.Context, query string) (*domain.PaginatedResponse[domain.SearchResult], error) {
	var search domain.PaginatedResponse[domain.SearchResult]
	err := s.get(ctx, "/search?q="+url.QueryEscape(query), &search)
	return &search, err
}

func (s *pkgsiteService) CheckHealth(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "HEAD", s.baseURL+"/package/fmt", nil)
	if err != nil {
		return err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pkgsite api returned %d", resp.StatusCode)
	}
	return nil
}
