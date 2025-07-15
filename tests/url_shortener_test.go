package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"

	"urlshortener/internal/http-server/handlers/url/save"
	"urlshortener/lib/api"
	"urlshortener/lib/random"
)

// Matches host in config
const (
	host = "localhost:8082"
)

// TestURLShortener_HappyPath tests successful URL shortening scenario
func TestURLShortener_HappyPath(t *testing.T) {
	// Base URL for test requests
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	// Create HTTP test client
	e := httpexpect.Default(t, u.String())

	// Test successful URL creation
	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.NewRandomString(10),
		}).
		WithBasicAuth("myuser", "mypass").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("alias") // Verify response contains alias field
}

// TestURLShortener_SaveRedirect tests various URL saving and redirect scenarios
// nolint:funlen
func TestURLShortener_SaveRedirect(t *testing.T) {
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL with custom alias",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:  "Invalid URL format",
			url:   "invalid_url",
			alias: gofakeit.Word(),
			error: "field URL is not a valid URL",
		},
		{
			name:  "Empty alias (should generate random)",
			url:   gofakeit.URL(),
			alias: "",
		},
		// TODO: Add more test cases:
		// - Duplicate alias
		// - Very long URL
		// - Special characters in alias
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			// --- Save URL Test ---
			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("myuser", "mypass").
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.error != "" {
				// Verify error response for invalid cases
				resp.NotContainsKey("alias")
				resp.Value("error").String().IsEqual(tc.error)
				return
			}

			alias := tc.alias
			if tc.alias != "" {
				// Verify custom alias was used
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				// Verify random alias was generated
				resp.Value("alias").String().NotEmpty()
				alias = resp.Value("alias").String().Raw()
			}

			// --- Redirect Test ---
			redirectURL := url.URL{
				Scheme: "http",
				Host:   host,
				Path:   alias,
			}

			redirectedToURL, err := api.GetRedirect(redirectURL.String())
			require.NoError(t, err)
			require.Equal(t, tc.url, redirectedToURL)

			// --- Cleanup: Delete URL ---
			resp = e.DELETE("/url/"+alias).
				WithBasicAuth("myuser", "mypass").
				Expect().Status(http.StatusOK).
				JSON().Object()

			resp.Value("status").String().IsEqual("OK")
		})
	}
}
