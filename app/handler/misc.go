package handler

import (
	"net/http"

	"github.com/hysem/mini-aspire-api/app/dto/response"
	"github.com/labstack/echo/v4"
)

// Misc handler
type Misc struct{}

// NewMisc returns a new Misc handler instance
func NewMisc() *Misc {
	return &Misc{}
}

// Ping handles the ping request.
func (h *Misc) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, response.APIResponse{
		Message: "pong",
	})
}
