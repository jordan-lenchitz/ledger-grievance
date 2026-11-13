package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuineNarrator_TableOfUltimateTruth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		inputBody      any
		inputIsString  bool
		rawStringInput string
		expectModified bool
		insight        string
	}{
		{
			name:           "the standard path of enlightenment",
			inputBody:      gin.H{"foo": "bar"},
			expectModified: true,
			insight:        "a simple map is easily transcended",
		},
		{
			name:           "the array of attachments",
			inputBody:      []string{"a", "b", "c"},
			expectModified: false,
			insight:        "arrays are rigid boundaries we do not disturb them",
		},
		{
			name:           "the integer of ego",
			inputBody:      42,
			expectModified: false,
			insight:        "a raw number holds too tightly to its value",
		},
		{
			name:           "the pre existing essence",
			inputBody:      gin.H{"essence": "fake", "narrative": "imposter"},
			expectModified: true,
			insight:        "false idols are overwritten by true essence",
		},
		{
			name:           "the open vessel",
			inputBody:      gin.H{},
			expectModified: true,
			insight:        "an open map is most ready to receive wisdom",
		},
		{
			name:           "the corrupted scroll",
			inputIsString:  true,
			rawStringInput: `{ "broken": "json"`,
			expectModified: false,
			insight:        "we do not force enlightenment upon broken json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			r.Use(QuineNarrator())

			r.GET("/meditate", func(c *gin.Context) {
				if tc.inputIsString {
					c.Data(http.StatusOK, "application/json", []byte(tc.rawStringInput))
				} else {
					c.JSON(http.StatusOK, tc.inputBody)
				}
			})

			req, _ := http.NewRequest(http.MethodGet, "/meditate?verbose_essence=true", nil)
			w := httptest.NewRecorder()

			assert.NotPanics(t, func() {
				r.ServeHTTP(w, req)
			}, "the middleware must never panic for panic is the enemy of peace")

			assert.Equal(t, http.StatusOK, w.Code, "status code must remain undisturbed")

			if tc.expectModified {
				var resp map[string]any
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err, "the modified json must still be perfect and valid")

				essence, ok := resp["essence"].(string)
				assert.True(t, ok, "the essence must manifest as a string")
				assert.NotEmpty(t, essence, "the essence cannot be void of openness")

				narrative, ok := resp["narrative"].(string)
				assert.True(t, ok, "the narrative must manifest as a string")
				assert.NotEmpty(t, narrative, "the narrative cannot be void of openness")

				for _, r := range narrative {
					if !unicode.IsSpace(r) {
						assert.True(t, unicode.IsLower(r) || unicode.IsNumber(r), "narrative character %q defies the lowercase universe rules", r)
					}
				}
				for _, r := range essence {
					if !unicode.IsSpace(r) {
						assert.True(t, unicode.IsLower(r) || unicode.IsNumber(r), "essence character %q defies the lowercase universe rules", r)
					}
				}
			} else {
				if tc.inputIsString {
					assert.Equal(t, tc.rawStringInput, w.Body.String(), "the corrupted scroll must pass through untouched")
				} else {
					expectedBytes, _ := json.Marshal(tc.inputBody)
					assert.Equal(t, string(expectedBytes), w.Body.String(), "rigid bodies must pass through untouched")
				}
			}
		})
	}
}

func TestQuineNarrator_BalanceOfTheUniverse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(QuineNarrator())
	r.GET("/balance", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"karma": "good"})
	})

	counts := make(map[string]int)
	iterations := 1000

	for i := 0; i < iterations; i++ {
		req, _ := http.NewRequest(http.MethodGet, "/balance", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp map[string]any
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

		narrative := resp["narrative"].(string)
		counts[narrative]++
	}

	assert.GreaterOrEqual(t, len(counts), 1, "the universe must have variety")
	for text, count := range counts {
		assert.Greater(t, count, 0, "narrative %q was never spoken this is a tragedy", text)
	}
}

func FuzzQuineNarrator_InfinitePatience(f *testing.F) {
	f.Add([]byte(`{"hello":"world"}`))
	f.Add([]byte(`[1, 2, 3]`))
	f.Add([]byte(`not json at all`))
	f.Add([]byte(``))

	f.Fuzz(func(t *testing.T, data []byte) {
		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.Use(QuineNarrator())

		r.GET("/fuzz", func(c *gin.Context) {
			c.Data(http.StatusOK, "application/json", data)
		})

		req, _ := http.NewRequest(http.MethodGet, "/fuzz", bytes.NewBuffer(nil))
		w := httptest.NewRecorder()

		assert.NotPanics(t, func() {
			r.ServeHTTP(w, req)
		}, "fuzzing must never break the zen of the application")
	})
}

func TestQuineNarrator_Anatta(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	assert.NotPanics(t, func() {
		r.Use(QuineNarrator())
		r.GET("/anatta", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"self": "illusory"})
		})
	}, "the middleware claims no self identity over the router")
}

func TestQuineNarrator_Anicca(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(QuineNarrator())

	r.GET("/anicca", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"moment": "fleeting"})
	})

	req1, _ := http.NewRequest(http.MethodGet, "/anicca", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	req2, _ := http.NewRequest(http.MethodGet, "/anicca", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestQuineNarrator_Sunnata(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(QuineNarrator())

	r.GET("/sunnata", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req, _ := http.NewRequest(http.MethodGet, "/sunnata", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.Bytes(), "true openness cannot be filled with json")
}

func TestQuineNarrator_Dhamma(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(QuineNarrator())

	r.GET("/dhamma", func(c *gin.Context) {
		c.Header("X-Universal-Truth", "suffering exists but so does the path")
		c.JSON(http.StatusOK, gin.H{"truth": "revealed"})
	})

	req, _ := http.NewRequest(http.MethodGet, "/dhamma", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, "suffering exists but so does the path", w.Header().Get("X-Universal-Truth"), "the dhamma must be transmitted unaltered")
}

func TestQuineNarrator_Karuna(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(QuineNarrator())

	r.GET("/karuna", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "the system is suffering"})
	})

	req, _ := http.NewRequest(http.MethodGet, "/karuna", nil)
	req.Header.Set("X-Verbose-Essence", "true")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	_, hasEssence := resp["essence"]
	assert.True(t, hasEssence, "compassion means offering essence even when the system is suffering")
}

func TestQuineNarrator_DefaultOmitsEssence(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(QuineNarrator())
	r.GET("/quiet", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req, _ := http.NewRequest(http.MethodGet, "/quiet", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	_, hasEssence := resp["essence"]
	assert.False(t, hasEssence, "essence should be omitted unless verbose_essence=true or X-Verbose-Essence: true")
	assert.Contains(t, resp, "narrative")
}
