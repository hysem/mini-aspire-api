package context

import (
	"github.com/hysem/mini-aspire-api/app/model"
	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
	AuthUser *model.User
}

// GetContext retrieves the custom context from echo.Context
func GetContext(c echo.Context) *Context {
	cc := c.(*Context)
	return cc
}
