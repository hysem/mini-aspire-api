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

	// GetLoanByID retrieves the loan details by id
	GetLoanByID(ctx context.Context, loanID uint64) (*model.Loan, error)

	// UpdateLoanStatus updates the status of the given loan
	UpdateLoanStatus(ctx context.Context, loanID uint64, approvedBy uint64, status model.LoanStatus, tx *sqlx.Tx) error

	// UpdateLoanEMIStatusByLoanID updates the status loan_emi entries for the given loan
	UpdateLoanEMIStatusByLoanID(ctx context.Context, loanID uint64, status model.LoanStatus, tx *sqlx.Tx) error

	// UpdateLoanEMIStatus updates the status loan_emi entries for the given loan
	UpdateLoanEMIStatus(ctx context.Context, loanID uint64, loanEMIIds []uint64, status model.LoanStatus, tx *sqlx.Tx) error

	// GetLoanEMIs get all loan emis for the given loan
	GetLoanEMIs(ctx context.Context, loanID uint64) ([]*model.LoanEMI, error)
}
