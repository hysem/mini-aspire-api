package middleware_test

import (
	"net/http"
	"testing"

	"github.com/hysem/mini-aspire-api/app/core/context"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware_Context(t *testing.T) {
	mw, _ := newMiddleware(t)

	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)

	handler := func(c echo.Context) error {
		_, ok := c.(*context.Context)
		if ok {
			return c.NoContent(http.StatusOK)
		}
		return c.NoContent(http.StatusBadRequest)
	}

	res := runMiddlewareTest(t, req, handler, mw.context)

	assert.Equal(t, http.StatusOK, res.Result().StatusCode)

}
