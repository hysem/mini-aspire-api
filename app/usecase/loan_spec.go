package usecase

import (
	"context"

	"github.com/hysem/mini-aspire-api/app/dto/request"
	"github.com/hysem/mini-aspire-api/app/dto/response"
)

// Loan interface declares supported operations on the loan usecase
//go:generate mockery --name=Loan --filename=loan_mock.go --inpackage

type Loan interface {
	// RequestLoan process a loan request
	RequestLoan(ctx context.Context, req *request.RequestLoan) (*response.RequestLoan, error)

	// ApproveLoan changes status of a loan and its associated loan_emi entries's status to APPROVED
	ApproveLoan(ctx context.Context, req *request.ApproveLoan) error
}
