package middleware

import (
	"net/http"

	"github.com/hysem/mini-aspire-api/app/core/context"
	"github.com/hysem/mini-aspire-api/app/core/message"
	"github.com/hysem/mini-aspire-api/app/dto/response"
	"github.com/hysem/mini-aspire-api/app/model"
	"github.com/labstack/echo/v4"
)

// GrantAccess middleware handles access to different resource based on role
func GrantAccess(roles ...model.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := context.GetContext(c)

			for _, role := range roles {
				if cc.AuthUser.Role == role {
					return next(cc)
				}
			}
			return c.JSON(http.StatusForbidden, response.APIResponse{
				Message: message.AccessDenied,
			})
		}
	}
}
