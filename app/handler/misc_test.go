package handler_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Misc_Ping(t *testing.T) {
	h, m := newHandler(t)
	defer m.assertExpectations(t)

	req, err := http.NewRequest(http.MethodGet, "/ping", nil)
	assert.NoError(t, err)

	res := runHandler(t, req, h.misc.Ping, nil)
	assert.Equal(t, http.StatusOK, res.Result().StatusCode)
	assert.JSONEq(t, `{"message":"pong"}`, res.Body.String())
}
