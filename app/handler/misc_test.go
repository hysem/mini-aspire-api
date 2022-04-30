package handler_test

import (
	"net/http"
	"testing"

	"github.com/hysem/mini-aspire-api/app/core/context"
	"github.com/stretchr/testify/assert"
)

func TestMisc_Ping(t *testing.T) {
	h, m := newHandler(t)
	defer m.assertExpectations(t)

	req, err := http.NewRequest(http.MethodGet, "/ping", nil)
	assert.NoError(t, err)

	res := runHandler(t, req, h.misc.Ping, &context.Context{})
	assert.Equal(t, http.StatusOK, res.Result().StatusCode)
	assert.JSONEq(t, `{"message":"pong"}`, res.Body.String())
}
