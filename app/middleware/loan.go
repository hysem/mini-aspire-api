package middleware

import (
	"net/http"
	"strconv"

	"github.com/hysem/mini-aspire-api/app/core/context"
	"github.com/hysem/mini-aspire-api/app/core/message"
	"github.com/hysem/mini-aspire-api/app/dto/response"
	"github.com/hysem/mini-aspire-api/app/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

// Auth middleware handles the authentication of requests
// The user will be available in AuthUser field of the custom context object
func Loan(loanRepository repository.Loan) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := context.GetContext(c)

			loanID, err := strconv.ParseUint(cc.Param("lid"), 10, 64)
			if err != nil {
				return c.JSON(http.StatusNotFound, response.APIResponse{
					Message: message.NoSuchResource,
				})
			}

			existingLoan, err := loanRepository.GetLoanByID(c.Request().Context(), loanID)
			if err != nil {
				log.Error("loanRepository.GetLoanByID() failed", zap.Error(err))
				return c.JSON(http.StatusInternalServerError, response.APIResponse{
					Message: message.InternalServerError,
				})
			}

			if existingLoan == nil {
				return c.JSON(http.StatusNotFound, response.APIResponse{
					Message: message.NoSuchResource,
				})
			}

			cc.Loan = existingLoan

			return next(cc)
		}
	}
}
