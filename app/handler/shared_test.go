package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hysem/mini-aspire-api/app/core/context"
	"github.com/hysem/mini-aspire-api/app/handler"
	"github.com/hysem/mini-aspire-api/app/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type handlerMocks struct {
	userUsecase usecase.MockUser
	loanUsecase usecase.MockLoan
}

type testHandler struct {
	misc *handler.Misc
	user *handler.User
	loan *handler.Loan
}

func newHandler(t *testing.T) (*testHandler, *handlerMocks) {
	m := &handlerMocks{}

	u := &testHandler{
		misc: handler.NewMisc(),
		user: handler.NewUser(&m.userUsecase),
		loan: handler.NewLoan(&m.loanUsecase),
	}
	return u, m
}

func (m *handlerMocks) assertExpectations(t *testing.T) {
	m.userUsecase.AssertExpectations(t)
	m.loanUsecase.AssertExpectations(t)
}

func runHandler(t *testing.T, req *http.Request, handler echo.HandlerFunc, setContext func(cc *context.Context)) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()
	c := echo.New().NewContext(req, res)
	cc := &context.Context{Context: c}
	if setContext != nil {
		setContext(cc)
	}
	err := handler(cc)
	assert.NoError(t, err)
	return res
}
