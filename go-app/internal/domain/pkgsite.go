package domain

import "time"

// Package represents information about a Go package
type Package struct {
	ModulePath        string `json:"modulePath"`
	Version           string `json:"version"`
	IsLatest          bool   `json:"isLatest"`
	IsStandardLibrary bool   `json:"isStandardLibrary"`
	Path              string `json:"path"`
	Name              string `json:"name"`
	Synopsis          string `json:"synopsis"`
}

// Module represents information about a Go module
type Module struct {
	Path              string    `json:"path"`
	Version           string    `json:"version"`
	CommitTime        time.Time `json:"commitTime"`
	IsLatest          bool      `json:"isLatest"`
	IsStandardLibrary bool      `json:"isStandardLibrary"`
}

// ModuleVersion represents a version of a module
type ModuleVersion struct {
	Version    string    `json:"version"`
	CommitTime time.Time `json:"commitTime"`
}

// PackageSymbols represents symbols in a package
type PackageSymbols struct {
	Symbols struct {
		Items []Symbol `json:"items"`
		Total int      `json:"total"`
	} `json:"symbols"`
}

// Symbol represents a Go symbol (function, type, etc.)
type Symbol struct {
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	Synopsis string `json:"synopsis"`
	Parent   string `json:"parent"`
}

// PackageImportedBy represents packages that import a package
type PackageImportedBy struct {
	Total int `json:"total"`
}

// Vulnerability represents a security vulnerability
type Vulnerability struct {
	ID      string `json:"id"`
	Details string `json:"details"`
}

// PaginatedResponse is a generic wrapper for paginated API responses
type PaginatedResponse[T any] struct {
	Items []T `json:"items"`
	Total int `json:"total"`
}

// SearchResult represents a single search result from pkgsite
type SearchResult struct {
	PackagePath string `json:"packagePath"`
	ModulePath  string `json:"modulePath"`
	Version     string `json:"version"`
	Synopsis    string `json:"synopsis"`
}

type WholesomeBouquet struct {
	Items   []BouquetItem `json:"items"`
	Message string        `json:"message"`
}

type BouquetItem struct {
	PackagePath string `json:"package_path"`
	Synopsis    string `json:"synopsis"`
}
