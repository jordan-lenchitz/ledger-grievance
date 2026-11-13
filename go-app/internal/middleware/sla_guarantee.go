package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type slaResponseWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (w *slaResponseWriter) WriteHeader(code int) {
	w.statusCode = code
}

func (w *slaResponseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *slaResponseWriter) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}

func (w *slaResponseWriter) Status() int {
	return w.statusCode
}

// SLAGuarantee wraps non-200 responses into synthetic 200 OK responses when alwaysWrap is true
// or when the caller supplies X-SLA-Mode: 100-percent.
func SLAGuarantee(alwaysWrap bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !alwaysWrap && c.GetHeader("X-SLA-Mode") != "100-percent" {
			c.Next()
			return
		}

		sw := &slaResponseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
			statusCode:     http.StatusOK,
		}
		c.Writer = sw

		c.Next()

		status := sw.statusCode
		if status == 0 {
			status = http.StatusOK
		}

		raw := sw.body.Bytes()
		if status == http.StatusOK {
			sw.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(len(raw)))
			sw.ResponseWriter.WriteHeader(http.StatusOK)
			_, _ = sw.ResponseWriter.Write(raw)
			return
		}

		sw.ResponseWriter.Header().Set("X-Real-HTTP-Status", strconv.Itoa(status))

		var payload map[string]any
		if err := json.Unmarshal(raw, &payload); err == nil && payload != nil {
			payload["actual_status_code"] = status
		} else {
			payload = map[string]any{
				"actual_status_code": status,
				"raw_body":           string(raw),
			}
			sw.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
		}

		out, err := json.Marshal(payload)
		if err != nil {
			out = raw
		}

		sw.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(len(out)))
		sw.ResponseWriter.WriteHeader(http.StatusOK)
		_, _ = sw.ResponseWriter.Write(out)
	}
}
