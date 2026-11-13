package main

import (
	"net/http"
	"testing"
)

func TestBuildRequest_ComplianceHeaders(t *testing.T) {
	apiURL = "http://localhost:8000"
	nanosecondEpoch = "1700000000123456789"
	slaMode = true
	t.Cleanup(func() {
		nanosecondEpoch = ""
		slaMode = false
	})

	req, err := buildRequest(http.MethodPost, "/incidents", map[string]any{"subject": "ok"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := req.Header.Get("Accept"); got != requiredAcceptMediaType {
		t.Fatalf("Accept = %q, want %q", got, requiredAcceptMediaType)
	}
	if got := req.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}
	if got := req.Header.Get("X-Nanosecond-Epoch"); got != "1700000000123456789" {
		t.Fatalf("X-Nanosecond-Epoch = %q, want 1700000000123456789", got)
	}
	if got := req.Header.Get("X-SLA-Mode"); got != "100-percent" {
		t.Fatalf("X-SLA-Mode = %q, want 100-percent", got)
	}

	getReq, err := buildRequest(http.MethodGet, "/incidents", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := getReq.Header.Get("Accept"); got != "" {
		t.Fatalf("GET Accept = %q, want empty", got)
	}
}
