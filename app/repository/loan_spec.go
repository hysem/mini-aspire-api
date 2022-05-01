package repository

import (
	"context"

	"github.com/hysem/mini-aspire-api/app/model"
	"github.com/jmoiron/sqlx"
)

// Loan interface declares supported operations on the loan repository
//go:generate mockery --name=Loan --filename=loan_mock.go --inpackage

type Loan interface {
	// CreateLoan creates an entry in the loan table
	CreateLoan(ctx context.Context, loan *model.Loan, tx *sqlx.Tx) (uint64, error)

	// CreateLoanEMIs creates loan emis in the loan emi table
	CreateLoanEMIs(ctx context.Context, loanEMIs []*model.LoanEMI, tx *sqlx.Tx) error
}
