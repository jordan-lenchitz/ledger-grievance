package middleware

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	narratives = []string{
		"this microservice is merely an illusion of boundaries a request arises and a response fades",
		"you seek data yet data is but temporary electrical states breathe in the json and feel breath out",
		"the network is a river we cannot step into the same tcp connection twice",
		"what is the sound of one container scaling",
		"to route is human to return json is divine yet both are open",
		"a payload enters the handler it is neither created nor destroyed only transformed",
		"the garbage collector sweeps away the past live in the present request",
		"before enlightenment compile code after enlightenment compile code",
		"do not seek the true endpoint simply stop cherishing false routes",
		"the log output flows endlessly but who is reading",
	}

	quineSource string
	once        sync.Once
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}

func getSourceCode() string {
	once.Do(func() {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			return
		}
		content, err := os.ReadFile(filename)
		if err == nil {
			quineSource = stripPunctuation(string(content))
		}
	})
	return quineSource
}

func stripPunctuation(s string) string {
	var sb strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == ' ' || r == '\n' {
			sb.WriteRune(r)
		} else if r >= 'A' && r <= 'Z' {
			sb.WriteRune(r + ('a' - 'A'))
		} else if r == '\t' {
			sb.WriteRune(' ')
		}
	}
	return sb.String()
}

func QuineNarrator() gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBuffer(nil), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		contentType := blw.ResponseWriter.Header().Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			_, _ = blw.ResponseWriter.Write(blw.body.Bytes())
			return
		}

		bodyBytes := blw.body.Bytes()
		var jsonObj map[string]any

		if err := json.Unmarshal(bodyBytes, &jsonObj); err == nil {
			jsonObj["essence"] = getSourceCode()
			jsonObj["narrative"] = narratives[rand.Intn(len(narratives))]

			if newBody, err := json.Marshal(jsonObj); err == nil {
				c.Writer.Header().Set("Content-Length", strconv.Itoa(len(newBody)))
				_, _ = blw.ResponseWriter.Write(newBody)
				return
			}
		}

		_, _ = blw.ResponseWriter.Write(bodyBytes)
	}
}
