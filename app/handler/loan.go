package handler

import (
	"net/http"

	"github.com/hysem/mini-aspire-api/app/core/context"
	"github.com/hysem/mini-aspire-api/app/core/message"
	"github.com/hysem/mini-aspire-api/app/dto/request"
	"github.com/hysem/mini-aspire-api/app/dto/response"
	"github.com/hysem/mini-aspire-api/app/usecase"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Loan handler
type Loan struct {
	loanUsecase usecase.Loan
}

// NewLoan returns a new Loan handler instance
func NewLoan(
	loanUsecase usecase.Loan,
) *Loan {
	return &Loan{
		loanUsecase: loanUsecase,
	}
}

func (h *Loan) RequestLoan(c echo.Context) error {
	var req request.RequestLoan

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

	req.UserID = context.GetContext(c).AuthUser.UserID

	resp, err := h.loanUsecase.RequestLoan(c.Request().Context(), &req)
	if err != nil {
		zap.L().Error("h.loanUsecase.ProcessLoanRequest() failed", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.APIResponse{
			Message: message.InternalServerError,
		})
	}

	return c.JSON(http.StatusCreated, response.APIResponse{
		Message: message.LoanRequestSuccess,
		Data:    resp,
	})
}
