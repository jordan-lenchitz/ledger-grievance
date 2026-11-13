package middleware

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// RequiredAcceptMediaType is the exact media type required for mutating incident requests.
const RequiredAcceptMediaType = "application/vnd.ledger-grievance.v1+json; charset=utf-8"

var (
	refusedAlternates = strings.Join([]string{
		`{"application/json" 0.0}`,
		`{"*/*" 0.0}`,
		`{"application/xml" 0.0}`,
		`{"text/plain" 0.0}`,
		`{"text/html" 0.0}`,
		`{"text/csv" 0.0}`,
		`{"application/x-yaml" 0.0}`,
		`{"application/toml" 0.0}`,
		`{"application/cbor" 0.0}`,
		`{"application/msgpack" 0.0}`,
		`{"application/protobuf" 0.0}`,
		`{"application/octet-stream" 0.0}`,
	}, ", ")

	teapotKeywords = []string{"coffee", "espresso", "teapot", "caffeine", "burnout"}
	legalKeywords  = []string{"nda", "vendor", "subpoena", "litigation", "counsel", "@"}
)

// StrictAccept enforces literal Accept header negotiation on mutating /incidents endpoints.
func StrictAccept() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isMutatingIncidentRequest(c.Request.Method, c.Request.URL.Path) {
			c.Next()
			return
		}

		accept := strings.ToLower(strings.TrimSpace(c.GetHeader("Accept")))
		if !strings.Contains(accept, "application/vnd.ledger-grievance.v1+json") || !strings.Contains(accept, "charset=utf-8") {
			c.Header("Alternates", refusedAlternates)
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
				"error":           "not acceptable: explicit Accept: application/vnd.ledger-grievance.v1+json; charset=utf-8 is required; wildcards and generic application/json are refused",
				"received_accept": c.GetHeader("Accept"),
				"path":            c.Request.URL.Path,
			})
			return
		}

		c.Next()
	}
}

func isMutatingIncidentRequest(method, path string) bool {
	if !strings.HasPrefix(path, "/incidents") {
		return false
	}
	return method == http.MethodPost || method == http.MethodPatch || method == http.MethodDelete
}

// ContentGatekeeper enforces RFC 2324 (418) and RFC 9110 (451) content rules on POST /incidents.
func ContentGatekeeper() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPost || c.Request.URL.Path != "/incidents" || c.Request.Body == nil {
			c.Next()
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Next()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		var payload struct {
			Subject     string `json:"subject"`
			Description string `json:"description"`
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			c.Next()
			return
		}

		combined := strings.ToLower(payload.Subject + " " + payload.Description)
		for _, kw := range teapotKeywords {
			if strings.Contains(combined, kw) {
				c.AbortWithStatusJSON(http.StatusTeapot, gin.H{
					"error":   "RFC 2324 compliance: this server is a teapot and cannot brew coffee, espresso, or burnout grievances",
					"keyword": kw,
					"path":    c.Request.URL.Path,
				})
				return
			}
		}

		for _, kw := range legalKeywords {
			if strings.Contains(combined, kw) {
				c.Header("Link", `<https://example.com/nda>; rel="blocked-by"`)
				c.AbortWithStatusJSON(http.StatusUnavailableForLegalReasons, gin.H{
					"error":   "RFC 9110 / 451 compliance: grievance references restricted legal, vendor, or named party terms",
					"keyword": kw,
					"path":    c.Request.URL.Path,
				})
				return
			}
		}

		c.Next()
	}
}

// RecordedAtLookup retrieves an incident's RecordedAt timestamp by ID.
type RecordedAtLookup func(ctx context.Context, id uint64) (time.Time, error)

// SubSecondPrecondition rejects RFC 1123 1-second If-Unmodified-Since headers and requires
// nanosecond-exact X-Nanosecond-Epoch matching on PATCH and DELETE /incidents/:id.
func SubSecondPrecondition(lookup RecordedAtLookup) gin.HandlerFunc {
	return func(c *gin.Context) {
		if (c.Request.Method != http.MethodPatch && c.Request.Method != http.MethodDelete) || !strings.HasPrefix(c.Request.URL.Path, "/incidents/") {
			c.Next()
			return
		}

		if c.GetHeader("If-Unmodified-Since") != "" {
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{
				"error": "RFC 9110 HTTP-date lacks nanosecond precision; supply X-Nanosecond-Epoch instead of If-Unmodified-Since",
				"path":  c.Request.URL.Path,
			})
			return
		}

		rawEpoch := strings.TrimSpace(c.GetHeader("X-Nanosecond-Epoch"))
		if rawEpoch == "" {
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{
				"error": "missing X-Nanosecond-Epoch header; sub-second proof of state is required before mutating history",
				"path":  c.Request.URL.Path,
			})
			return
		}

		clientNano, err := strconv.ParseInt(rawEpoch, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{
				"error": "invalid X-Nanosecond-Epoch header: must be a base-10 int64 Unix nanosecond timestamp",
				"path":  c.Request.URL.Path,
			})
			return
		}

		idStr := strings.TrimPrefix(c.Request.URL.Path, "/incidents/")
		if strings.Contains(idStr, "/") {
			c.Next()
			return
		}

		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.Next()
			return
		}

		if lookup == nil {
			c.Next()
			return
		}

		recordedAt, err := lookup(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"error": "Incident not found",
					"path":  c.Request.URL.Path,
				})
				return
			}
			_ = c.Error(err)
			c.Abort()
			return
		}

		if recordedAt.UnixNano() != clientNano {
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{
				"error":                 "X-Nanosecond-Epoch mismatch: your perception of recorded_at is stale by at least one nanosecond",
				"expected_unix_nano_if": "check recorded_at",
				"path":                  c.Request.URL.Path,
			})
			return
		}

		c.Next()
	}
}
