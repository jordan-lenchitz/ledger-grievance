package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jordan-lenchitz/ledger-grievance/go-app/internal/domain"
)

type PkgsiteService interface {
	GetPackage(ctx context.Context, path string) (*domain.Package, error)
	GetModule(ctx context.Context, path string) (*domain.Module, error)
	GetVersions(ctx context.Context, path string) (*domain.PaginatedResponse[domain.ModuleVersion], error)
	GetSymbols(ctx context.Context, path string) (*domain.PackageSymbols, error)
	GetImportedBy(ctx context.Context, path string) (*domain.PackageImportedBy, error)
	GetVulns(ctx context.Context, path string) (*domain.PaginatedResponse[domain.Vulnerability], error)
	Search(ctx context.Context, query string) (*domain.PaginatedResponse[domain.SearchResult], error)
}

type pkgsiteService struct {
	baseURL string
	client  *http.Client
}

func NewPkgsiteService(baseURL string) PkgsiteService {
	if baseURL == "" {
		baseURL = "https://pkg.go.dev/v1beta"
	}
	return &pkgsiteService{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (s *pkgsiteService) get(ctx context.Context, endpoint string, result any) error {
	req, err := http.NewRequestWithContext(ctx, "GET", s.baseURL+endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api error: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(result)
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
