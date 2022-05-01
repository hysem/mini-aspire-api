package middleware

import (
	"net/http"

	"github.com/hysem/mini-aspire-api/app/core/context"
	"github.com/hysem/mini-aspire-api/app/core/jwt"
	"github.com/hysem/mini-aspire-api/app/core/message"
	"github.com/hysem/mini-aspire-api/app/dto/response"
	"github.com/hysem/mini-aspire-api/app/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

// Auth middleware handles the authentication of requests
// The user will be available in AuthUser field of the custom context object
func Auth(jwt jwt.JWT, userRepository repository.User) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := context.GetContext(c)
			token := cc.Request().Header.Get(echo.HeaderAuthorization)

			userID, err := jwt.Parse(token)
			if err != nil {
				log.Error("invalid token", zap.Error(err), zap.String("token", token))
				return c.JSON(http.StatusUnauthorized, response.APIResponse{
					Message: message.InvalidSession,
				})
			}

			existingUser, err := userRepository.GetByID(c.Request().Context(), userID)
			if err != nil {
				log.Error("userRepository.GetByID() failed", zap.Error(err))
				return c.JSON(http.StatusInternalServerError, response.APIResponse{
					Message: message.InternalServerError,
				})
			}

			if existingUser == nil {
				return c.JSON(http.StatusUnauthorized, response.APIResponse{
					Message: message.InvalidSession,
				})
			}

			cc.AuthUser = existingUser

			return next(cc)
		}
	}
}
