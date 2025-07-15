package redirect_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlshortener/internal/http-server/handlers/url/redirect"
	"urlshortener/internal/http-server/handlers/url/redirect/mocks"
	"urlshortener/internal/storage"
	"urlshortener/lib/api/response"
	"urlshortener/lib/logger/handlers/slogdiscard"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name          string
		alias         string
		wantStatus    int
		wantResponse  response.Response
		mockError     error
		mockCalled    bool
		mockReturnURL string
	}{
		{
			name:          "Redirect succes",
			alias:         "test_alias",
			wantStatus:    http.StatusFound,
			wantResponse:  response.Response{},
			mockCalled:    true,
			mockReturnURL: "https://example.com",
		},
		{
			name:       "Empty alias",
			alias:      "",
			wantStatus: http.StatusBadRequest,
			wantResponse: response.Response{
				Status: response.StatusError,
				Error:  "alias is required",
			},
			mockCalled: false,
		},
		{
			name:       "URL not found",
			alias:      "url not found",
			wantStatus: http.StatusNotFound,
			wantResponse: response.Response{
				Status: response.StatusError,
				Error:  "not found",
			},
			mockCalled: true,
			mockError:  storage.ErrUrlNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlRedirecterMock := mocks.NewURLGetter(t)

			if tc.mockCalled {
				urlRedirecterMock.On("GetURL", tc.alias).
					Return(tc.mockReturnURL, tc.mockError).
					Once()
			}

			handler := redirect.New(slogdiscard.NewDiscardLogger(), urlRedirecterMock)

			req, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("alias", tc.alias)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.wantStatus, rr.Code)

			if tc.wantStatus == http.StatusFound {
				require.Equal(t, tc.mockReturnURL, rr.Header().Get("Location"))
			} else {
				var resp response.Response
				err = json.Unmarshal(rr.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Equal(t, tc.wantResponse.Status, resp.Status)
				require.Equal(t, tc.wantResponse.Error, resp.Error)
			}

			if tc.mockCalled {
				urlRedirecterMock.AssertExpectations(t)
			}
		})
	}
}
