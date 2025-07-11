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
			name:  "Success",
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

			// 1. Инициализация мока
			urlDeleterMock := mocks.NewURLDeleter(t)

			// 2. Настройка мока только если он должен вызываться
			if tc.mockCalled {
				urlDeleterMock.On("DeleteURL", tc.alias).
					Return(tc.mockError).
					Once()
			}

			// 3. Создание обработчика
			handler := delete.New(slogdiscard.NewDiscardLogger(), urlDeleterMock)

			// 4. Создание запроса
			req, err := http.NewRequest(http.MethodDelete, "/", nil)
			// Проверка отстутсвия ошибки в req, если есть ошибка, то мы тест остановится и упадёт
			require.NoError(t, err)

			// 5. Установка параметров маршрута для chi
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("alias", tc.alias)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// 6. Создание ResponseRecorder
			rr := httptest.NewRecorder()

			// 7. Вызов обработчика
			handler.ServeHTTP(rr, req)

			// 8. Проверка статус кода
			require.Equal(t, tc.wantStatus, rr.Code)

			// 9. Декодирование ответа
			var resp response.Response
			err = json.Unmarshal(rr.Body.Bytes(), &resp)
			require.NoError(t, err)

			// 10. Проверка тела ответа
			require.Equal(t, tc.wantResponse.Status, resp.Status)
			require.Equal(t, tc.wantResponse.Error, resp.Error)

			// 11. Проверка вызовов мока
			if tc.mockCalled {
				urlDeleterMock.AssertExpectations(t)
			}
		})
	}
}
