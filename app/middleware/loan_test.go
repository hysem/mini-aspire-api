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

func TestMiddleware_Loan(t *testing.T) {
	loan := &model.Loan{
		ID: 1,
	}
	authorizationHeader := "Bearer token"

	testCases := map[string]struct {
		lid                  string
		setMocks             func(m *middlewareMocks)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		`error case: invalid path param`: {
			lid: "invalid",
			setMocks: func(m *middlewareMocks) {
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"message":"no such resource"}`,
		},
		`error case: failed to get loan details`: {
			lid: "1",
			setMocks: func(m *middlewareMocks) {
				m.loanRepository.On("GetLoanByID", mock.Anything, loan.ID).Return(nil, assert.AnError)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
		`error case: no such loan`: {
			lid: "1",
			setMocks: func(m *middlewareMocks) {
				m.loanRepository.On("GetLoanByID", mock.Anything, loan.ID).Return(nil, nil)
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"message":"no such resource"}`,
		},
		`success case: found loan`: {
			lid: "1",
			setMocks: func(m *middlewareMocks) {
				m.loanRepository.On("GetLoanByID", mock.Anything, loan.ID).Return(loan, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"amount":"0", "approved_by":0, "created_at":"0001-01-01T00:00:00Z", "loan_id":1, "purpose":"", "status":"", "terms":0, "updated_at":"0001-01-01T00:00:00Z", "user_id":0}`,
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
				return c.JSON(http.StatusOK, cc.Loan)
			}

			res := runMiddlewareTest(t, req, handler, mw.context, modifyContext(func(c *context.Context) {
				c.SetParamNames("lid")
				c.SetParamValues(tc.lid)
			}), mw.loan)

			assert.Equal(t, tc.expectedStatusCode, res.Result().StatusCode)
			assert.JSONEq(t, tc.expectedResponseBody, res.Body.String())
		})
	}

}
