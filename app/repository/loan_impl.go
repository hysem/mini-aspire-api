package repository

import (
	"context"
	"database/sql"

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
	(purpose, amount, terms, status, user_id, approved_by, created_at, updated_at)
	VALUES
	($1, $2, $3, $4, $5, NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
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

// GetLoanByID retrieves the loan details by id
func (r *_loan) GetLoanByID(ctx context.Context, loanID uint64) (*model.Loan, error) {
	const query = `SELECT 
		loan_id, purpose, amount, terms, status, user_id, approved_by, created_at, updated_at 
	FROM loan WHERE loan_id=$1`

	var loan model.Loan
	err := r.masterDB.GetContext(ctx, &loan, query, loanID)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to create loan")
	}

	return &loan, nil
}

// UpdateLoanStatus updates the status of the given loan
func (r *_loan) UpdateLoanStatus(ctx context.Context, loanID uint64, approvedBy uint64, status model.LoanStatus, tx *sqlx.Tx) error {
	const query = `UPDATE loan SET status=$1, approved_by=$2 WHERE loan_id=$3`
	if _, err := tx.ExecContext(ctx, query, status, approvedBy, loanID); err != nil {
		return errors.Wrap(err, "failed to update loan status")
	}
	return nil
}

// UpdateLoanStatus updates the status loan_emi entries for the given loan
func (r *_loan) UpdateLoanEMIStatusByLoanID(ctx context.Context, loanID uint64, status model.LoanStatus, tx *sqlx.Tx) error {
	const query = `UPDATE loan_emi SET status=$1 WHERE loan_id=$2`
	if _, err := tx.ExecContext(ctx, query, status, loanID); err != nil {
		return errors.Wrap(err, "failed to update loan_emi status")
	}
	return nil
}
