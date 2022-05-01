package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hysem/mini-aspire-api/app/core/context"
	"github.com/hysem/mini-aspire-api/app/core/jwt"
	"github.com/hysem/mini-aspire-api/app/middleware"
	"github.com/hysem/mini-aspire-api/app/model"
	"github.com/hysem/mini-aspire-api/app/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type middlewareMocks struct {
	userRepository repository.MockUser
	loanRepository repository.MockLoan
	jwt            jwt.MockJWT
}

type testMiddleware struct {
	auth        echo.MiddlewareFunc
	loan        echo.MiddlewareFunc
	context     echo.MiddlewareFunc
	grantAccess func(roles ...model.Role) echo.MiddlewareFunc
}

func newMiddleware(t *testing.T) (*testMiddleware, *middlewareMocks) {
	m := &middlewareMocks{}

	u := &testMiddleware{
		auth:        middleware.Auth(&m.jwt, &m.userRepository),
		loan:        middleware.Loan(&m.loanRepository),
		context:     middleware.Context,
		grantAccess: middleware.GrantAccess,
	}
	return u, m
}

func (m *middlewareMocks) assertExpectations(t *testing.T) {
	m.userRepository.AssertExpectations(t)
	m.loanRepository.AssertExpectations(t)
	m.jwt.AssertExpectations(t)
}

func runMiddlewareTest(t *testing.T, req *http.Request, handler echo.HandlerFunc, mws ...echo.MiddlewareFunc) *httptest.ResponseRecorder {
	for i := len(mws) - 1; i >= 0; i-- {
		handler = mws[i](handler)
	}
	res := httptest.NewRecorder()
	c := echo.New().NewContext(req, res)
	cc := context.Context{Context: c}

	err := handler(cc)
	assert.NoError(t, err)
	return res
}

func modifyContext(fn func(c *context.Context)) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := context.GetContext(c)
			fn(cc)

			return next(cc)
		}
	}
}
