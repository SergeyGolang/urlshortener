package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"

	"urlshortener/internal/http-server/handlers/url/save"
	"urlshortener/lib/random"
)

// совпадение с host в конфиге
const (
	host = "localhost:8082"
)

func TestURLShortener_HappyPath(t *testing.T) {
	// формирование базового URL к которму будет обращаться клиент
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	// создание клиента с помощью которого будут выполняться запросы
	e := httpexpect.Default(t, u.String())

	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.NewRandomString(10),
		}).
		WithBasicAuth("Myuser", "Mypassword").
		// Expect - что мы ожидаем от ответа
		Expect().
		Status(200).
		JSON().Object().
		// Содержание параметра alias
		ContainsKey("alias")
}

//nolint:funlen
func TestURLShortener_SaveRedirect(t *testing.T) {
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   "https://example.com",
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:  "Invalid URL",
			url:   "invalid_url",
			alias: gofakeit.Word(),
			error: "field URL is not a valid URL",
		},
		{
			name:  "Empty Alias",
			url:   "https://example.com",
			alias: "",
		},
		// TODO: add more test cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			// Save

			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("Myuser", "Mypassword").
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.error != "" {
				// если в тест кейсе ожидаем ошибку, то проверяем отсутствие alias в ответе
				resp.NotContainsKey("alias")
				// проверяем наличие поле ошибки и сравниваем с кейсовой ошибкой
				resp.Value("error").String().IsEqual(tc.error)

				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()
				// сохраняем alias, по которому нужно потом обратиться для редиректа
				alias = resp.Value("alias").String().Raw()
			}

			// Redirect
			redirectResp := e.GET("/" + alias).
				Expect().
				Status(http.StatusOK) // 302

			// Проверяем, что редирект ведёт на исходный URL
			location := redirectResp.Header("Location").Raw()
			require.Equal(t, tc.url, location)
		})
	}
}
