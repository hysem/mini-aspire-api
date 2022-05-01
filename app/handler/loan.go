package handler

import (
	"net/http"

	"github.com/hysem/mini-aspire-api/app/core/context"
	"github.com/hysem/mini-aspire-api/app/core/message"
	"github.com/hysem/mini-aspire-api/app/dto/request"
	"github.com/hysem/mini-aspire-api/app/dto/response"
	"github.com/hysem/mini-aspire-api/app/model"
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

// RequestLoan handles the loan request
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

// ApproveLoan handles the loan approval request
func (h *Loan) ApproveLoan(c echo.Context) error {
	var req request.ApproveLoan

	cc := context.GetContext(c)

	req.ApprovedBy = cc.AuthUser.UserID
	req.LoanID = cc.Loan.ID
	if cc.Loan.Status != model.LoanStatusPending {
		return c.JSON(http.StatusOK, response.APIResponse{
			Message: message.AlreadyApprovedLoan,
		})
	}

	err := h.loanUsecase.ApproveLoan(c.Request().Context(), &req)
	if err != nil {
		zap.L().Error("h.loanUsecase.ApproveLoan() failed", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.APIResponse{
			Message: message.InternalServerError,
		})
	}

	return c.JSON(http.StatusOK, response.APIResponse{
		Message: message.ApprovedLoan,
	})
}

// GetLoan handles the loan approval request
func (h *Loan) GetLoan(c echo.Context) error {
	cc := context.GetContext(c)

	if cc.AuthUser.UserID != cc.Loan.UserID {
		return c.JSON(http.StatusForbidden, response.APIResponse{
			Message: message.AccessDenied,
		})
	}

	var req request.GetLoan
	req.Loan = cc.Loan

	resp, err := h.loanUsecase.GetLoan(c.Request().Context(), &req)
	if err != nil {
		zap.L().Error("h.loanUsecase.ApproveLoan() failed", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.APIResponse{
			Message: message.InternalServerError,
		})
	}
	return c.JSON(http.StatusOK, response.APIResponse{
		Data: resp,
	})
}
