package middleware_test

import (
	"net/http"
	"testing"

	"github.com/hysem/mini-aspire-api/app/core/context"
	"github.com/hysem/mini-aspire-api/app/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMiddleware_Auth(t *testing.T) {
	user := &model.User{
		UserID:   1,
		Name:     "name",
		Email:    "email@yopmail.com",
		Password: "password",
		Role:     model.RoleCustomer,
	}
	authorizationHeader := "Bearer token"

	testCases := map[string]struct {
		setMocks             func(m *middlewareMocks)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		`error case: invalid authorization header`: {
			setMocks: func(m *middlewareMocks) {
				m.jwt.On("Parse", authorizationHeader).Return(uint64(0), assert.AnError)
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"message":"invalid session"}`,
		},
		`error case: failed to get user by id`: {
			setMocks: func(m *middlewareMocks) {
				m.jwt.On("Parse", authorizationHeader).Return(user.UserID, nil)
				m.userRepository.On("GetByID", mock.Anything, user.UserID).Return(nil, assert.AnError)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
		`error case: no such user`: {
			setMocks: func(m *middlewareMocks) {
				m.jwt.On("Parse", authorizationHeader).Return(user.UserID, nil)
				m.userRepository.On("GetByID", mock.Anything, user.UserID).Return(nil, nil)
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"message":"invalid session"}`,
		},
		`success case: got response`: {
			setMocks: func(m *middlewareMocks) {
				m.jwt.On("Parse", authorizationHeader).Return(user.UserID, nil)
				m.userRepository.On("GetByID", mock.Anything, user.UserID).Return(user, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"created_at":"0001-01-01T00:00:00Z", "email":"email@yopmail.com", "name":"name", "role":"customer", "updated_at":"0001-01-01T00:00:00Z", "user_id":1}`,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mw, m := newMiddleware(t)
			defer m.assertExpectations(t)
			tc.setMocks(m)

			req, err := http.NewRequest("GET", "/test", nil)
			assert.NoError(t, err)
			req.Header.Set(echo.HeaderAuthorization, authorizationHeader)

			handler := func(c echo.Context) error {
				cc := context.GetContext(c)
				return c.JSON(http.StatusOK, cc.AuthUser)
			}

			res := runMiddlewareTest(t, req, handler, mw.context, mw.auth)

			assert.Equal(t, tc.expectedStatusCode, res.Result().StatusCode)
			assert.JSONEq(t, tc.expectedResponseBody, res.Body.String())
		})
	}

}
