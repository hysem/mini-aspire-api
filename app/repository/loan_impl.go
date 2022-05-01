package repository

import (
	"context"

	"github.com/hysem/mini-aspire-api/app/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var _ Loan = (*_loan)(nil)

// _loan implements Loan interface
type _loan struct {
	masterDB *sqlx.DB
}

func NewLoan(
	masterDB *sqlx.DB,
) *_loan {
	return &_loan{
		masterDB: masterDB,
	}
}

// CreateLoan creates an entry in the loan table
func (r *_loan) CreateLoan(ctx context.Context, loan *model.Loan, tx *sqlx.Tx) (uint64, error) {
	const query = `INSERT INTO loan 
	(purpose, amount, terms, status, user_id, created_at, updated_at)
	VALUES
	($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	RETURNING loan_id`

	var loanID uint64
	if err := tx.QueryRowContext(ctx, query,
		loan.Purpose,
		loan.Amount,
		loan.Terms,
		loan.Status,
		loan.UserID,
	).Scan(&loanID); err != nil {
		return loanID, errors.Wrap(err, "failed to create loan")
	}

	return loanID, nil
}

// CreateLoanEMIs creates loan emis in the loan emi table
func (r *_loan) CreateLoanEMIs(ctx context.Context, loanEMIs []*model.LoanEMI, tx *sqlx.Tx) error {
	const query = `INSERT INTO loan_emi
	(loan_id, seq_no, due_date, amount, status, created_at, updated_at)
	VALUES
	(:loan_id, :seq_no, :due_date, :amount, :status, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	if _, err := tx.NamedExecContext(ctx, query, loanEMIs); err != nil {
		return errors.Wrap(err, "failed to create loan")
	}

	return nil
}
