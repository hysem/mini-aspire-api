package handler_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/hysem/mini-aspire-api/app/core/apierr"
	"github.com/hysem/mini-aspire-api/app/core/context"
	"github.com/hysem/mini-aspire-api/app/dto/request"
	"github.com/hysem/mini-aspire-api/app/dto/response"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUser_Generate(t *testing.T) {
	const validRequestBody = `{
		"email": "email@yopmail.com",
		"password": "12345678"
	}`
	testCases := map[string]struct {
		body                 string
		setMocks             func(m *handlerMocks)
		expectedResponseBody string
		expectedStatusCode   int
	}{
		`error case: failed to parse request`: {
			body: `{`,
			setMocks: func(m *handlerMocks) {

			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"failed to parse request"}`,
		},
		`error case: failed to validate request`: {
			body: `{}`,
			setMocks: func(m *handlerMocks) {

			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error": {"email":"cannot be blank", "password":"cannot be blank"}, "message":"failed to validate request"}`,
		},
		`error case: failed to generate token`: {
			body: validRequestBody,
			setMocks: func(m *handlerMocks) {
				m.userUsecase.On("GenerateToken", mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
		`error case: failed to generate token; invalid credentials`: {
			body: validRequestBody,
			setMocks: func(m *handlerMocks) {
				m.userUsecase.On("GenerateToken", mock.Anything, mock.Anything).Return(nil, apierr.ErrInvalidCredentials)
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"message":"invalid credentials"}`,
		},
		`success case: generated token`: {
			body: validRequestBody,
			setMocks: func(m *handlerMocks) {
				m.userUsecase.On("GenerateToken", mock.Anything, &request.UserGenerateToken{
					Email:    "email@yopmail.com",
					Password: "12345678",
				}).Return(&response.UserGenerateToken{
					Token: "abcd",
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"data":{"token":"abcd"}}`,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			h, m := newHandler(t)
			defer m.assertExpectations(t)
			tc.setMocks(m)

			req, err := http.NewRequest(http.MethodPost, "/user/token", bytes.NewBufferString(tc.body))
			req.Header.Add(echo.HeaderContentType, "application/json")
			assert.NoError(t, err)

			cc := &context.Context{}

			res := runHandler(t, req, h.user.GenerateToken, cc)
			assert.Equal(t, tc.expectedStatusCode, res.Result().StatusCode)
			assert.JSONEq(t, tc.expectedResponseBody, res.Body.String())
		})
	}
}
