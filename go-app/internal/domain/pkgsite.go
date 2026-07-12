package domain

import "time"

type Package struct {
	ModulePath        string `json:"modulePath"`
	Version           string `json:"version"`
	IsLatest          bool   `json:"isLatest"`
	IsStandardLibrary bool   `json:"isStandardLibrary"`
	Path              string `json:"path"`
	Name              string `json:"name"`
	Synopsis          string `json:"synopsis"`
}

type Module struct {
	Path              string    `json:"path"`
	Version           string    `json:"version"`
	CommitTime        time.Time `json:"commitTime"`
	IsLatest          bool      `json:"isLatest"`
	IsStandardLibrary bool      `json:"isStandardLibrary"`
}

type ModuleVersion struct {
	Version    string    `json:"version"`
	CommitTime time.Time `json:"commitTime"`
}

type PackageSymbols struct {
	Symbols struct {
		Items []Symbol `json:"items"`
		Total int      `json:"total"`
	} `json:"symbols"`
}

type Symbol struct {
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	Synopsis string `json:"synopsis"`
	Parent   string `json:"parent"`
}

type PackageImportedBy struct {
	Total int `json:"total"`
}

type Vulnerability struct {
	ID      string `json:"id"`
	Details string `json:"details"`
}

type PaginatedResponse[T any] struct {
	Items []T `json:"items"`
	Total int `json:"total"`
}

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
