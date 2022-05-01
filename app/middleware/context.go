package middleware

import (
	"github.com/hysem/mini-aspire-api/app/core/context"
	"github.com/labstack/echo/v4"
)

// Context middleware attaches core.Context to the request context
func Context(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &context.Context{
			Context: c,
		}
		return next(cc)
	}
}
