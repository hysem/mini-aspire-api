package repository

import (
	"context"
	"database/sql"

	"github.com/hysem/mini-aspire-api/app/model"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
		return errors.Wrap(err, "failed to create loan emis")
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
		return nil, errors.Wrap(err, "failed to get loan")
	}

	return &loan, nil
}

// UpdateLoanStatus updates the status of the given loan
func (r *_loan) UpdateLoanStatus(ctx context.Context, loanID uint64, approvedBy uint64, status model.LoanStatus, tx *sqlx.Tx) error {
	const query = `UPDATE loan SET status=$1, approved_by=$2, updated_at=CURRENT_TIMESTAMP WHERE loan_id=$3`
	if _, err := tx.ExecContext(ctx, query, status, approvedBy, loanID); err != nil {
		return errors.Wrap(err, "failed to update loan status")
	}
	return nil
}

// UpdateLoanEMIStatusByLoanID updates the status loan_emi entries for the given loan
func (r *_loan) UpdateLoanEMIStatusByLoanID(ctx context.Context, loanID uint64, status model.LoanStatus, tx *sqlx.Tx) error {
	const query = `UPDATE loan_emi SET status=$1, updated_at=CURRENT_TIMESTAMP WHERE loan_id=$2`
	if _, err := tx.ExecContext(ctx, query, status, loanID); err != nil {
		return errors.Wrap(err, "failed to update loan_emi status")
	}
	return nil
}

// UpdateLoanEMIStatus updates the status loan_emi entries for the given loan
func (r *_loan) UpdateLoanEMIStatus(ctx context.Context, loanID uint64, loanEMIIds []uint64, status model.LoanStatus, tx *sqlx.Tx) error {
	const query = `UPDATE loan_emi SET status=$1, updated_at=CURRENT_TIMESTAMP WHERE loan_id=$2 AND loan_emi_id=ANY($3)`
	if _, err := tx.ExecContext(ctx, query, status, loanID, pq.Array(loanEMIIds)); err != nil {
		return errors.Wrap(err, "failed to update loan_emi status")
	}
	return nil
}

// GetLoanEMIs get all loan emis for the given loan
func (r *_loan) GetLoanEMIs(ctx context.Context, loanID uint64) ([]*model.LoanEMI, error) {
	const query = `SELECT 
		loan_emi_id, loan_id, seq_no, due_date, amount, status, created_at, updated_at 
	FROM loan_emi 
	WHERE loan_id=$1
	ORDER BY seq_no`

	var loanEMIs []*model.LoanEMI
	err := r.masterDB.SelectContext(ctx, &loanEMIs, query, loanID)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to get loan emis")
	}

	return loanEMIs, nil
}
