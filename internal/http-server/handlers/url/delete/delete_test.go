package delete_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"urlshortener/internal/http-server/handlers/url/delete"
	"urlshortener/internal/http-server/handlers/url/delete/mocks"
	"urlshortener/internal/storage"
	"urlshortener/lib/api/response"
	"urlshortener/lib/logger/handlers/slogdiscard"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name         string
		alias        string
		wantStatus   int
		wantResponse response.Response
		mockError    error
		mockCalled   bool
	}{
		{
			name:  "Delete success",
			alias: "test_alias",
			wantResponse: response.Response{
				Status: response.StatusOK,
			},
			wantStatus: http.StatusOK,
			mockCalled: true,
		},
		{
			name:  "Empty alias",
			alias: "",
			wantResponse: response.Response{
				Status: response.StatusError,
				Error:  "alias is required",
			},
			wantStatus: http.StatusBadRequest,
			mockCalled: false,
		},
		{
			name:  "URL not found",
			alias: "not_found",
			wantResponse: response.Response{
				Status: response.StatusError,
				Error:  "url not found",
			},
			wantStatus: http.StatusNotFound,
			mockError:  storage.ErrUrlNotFound,
			mockCalled: true,
		},
		{
			name:  "Internal error",
			alias: "internal_error",
			wantResponse: response.Response{
				Status: response.StatusError,
				Error:  "failed to delete url",
			},
			wantStatus: http.StatusInternalServerError,
			mockError:  errors.New("internal error"),
			mockCalled: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlDeleterMock := mocks.NewURLDeleter(t)

			if tc.mockCalled {
				urlDeleterMock.On("DeleteURL", tc.alias).
					Return(tc.mockError).
					Once()
			}

			handler := delete.New(slogdiscard.NewDiscardLogger(), urlDeleterMock)

			req, err := http.NewRequest(http.MethodDelete, "/", nil)
			require.NoError(t, err)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("alias", tc.alias)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.wantStatus, rr.Code)

			var resp response.Response
			err = json.Unmarshal(rr.Body.Bytes(), &resp)
			require.NoError(t, err)

			require.Equal(t, tc.wantResponse.Status, resp.Status)
			require.Equal(t, tc.wantResponse.Error, resp.Error)

			if tc.mockCalled {
				urlDeleterMock.AssertExpectations(t)
			}
		})
	}
}
