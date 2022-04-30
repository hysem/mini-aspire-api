package handler

import (
	"errors"
	"net/http"

	"github.com/hysem/mini-aspire-api/app/core/apierr"
	"github.com/hysem/mini-aspire-api/app/core/message"
	"github.com/hysem/mini-aspire-api/app/dto/request"
	"github.com/hysem/mini-aspire-api/app/dto/response"
	"github.com/hysem/mini-aspire-api/app/usecase"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// User handler
type User struct {
	userUsecase usecase.User
}

// NewUser returns a new User handler instance
func NewUser(
	userUsecase usecase.User,
) *User {
	return &User{
		userUsecase: userUsecase,
	}
}

func (h *User) GenerateToken(c echo.Context) error {
	var req request.UserGenerateToken

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: message.ParsingFailed,
		})
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, response.APIResponse{
			Error:   err,
			Message: message.ValidationFailed,
		})
	}

	resp, err := h.userUsecase.GenerateToken(c.Request().Context(), &req)
	switch {
	case errors.Is(err, apierr.ErrInvalidCredentials):
		return c.JSON(http.StatusUnauthorized, response.APIResponse{
			Message: message.InvalidCredentials,
		})
	case err != nil:
		zap.L().Error("h.userUsecase.GenerateToken() failed", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.APIResponse{
			Message: message.InternalServerError,
		})
	}

	return c.JSON(http.StatusOK, response.APIResponse{
		Data: resp,
	})
}
