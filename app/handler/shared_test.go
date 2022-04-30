package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hysem/mini-aspire-api/app/core/context"
	"github.com/hysem/mini-aspire-api/app/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type handlerMocks struct {
}

type testHandler struct {
	misc *handler.Misc
}

func newHandler(t *testing.T) (*testHandler, *handlerMocks) {
	m := &handlerMocks{}

	u := &testHandler{
		misc: handler.NewMisc(),
	}
	return u, m
}

func (m *handlerMocks) assertExpectations(t *testing.T) {
}

func runHandler(t *testing.T, req *http.Request, handler echo.HandlerFunc, ctx *context.Context) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()
	ctx.Context = echo.New().NewContext(req, res)
	err := handler(ctx)
	assert.NoError(t, err)
	return res
}
