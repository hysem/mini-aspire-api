package usecase

import (
	"context"
	"time"

	"github.com/hysem/mini-aspire-api/app/core/apierr"
	"github.com/hysem/mini-aspire-api/app/dto/request"
	"github.com/hysem/mini-aspire-api/app/dto/response"
	"github.com/hysem/mini-aspire-api/app/model"
	"github.com/hysem/mini-aspire-api/app/repository"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var _ Loan = (*_loan)(nil)

// _loan implements Loan usecase
type _loan struct {
	loanRepository repository.Loan
	baseRepository repository.Base
}

// NewLoan returns an instance of _loan usecase
func NewLoan(
	loanRepository repository.Loan,
	baseRepository repository.Base,
) *_loan {
	return &_loan{
		loanRepository: loanRepository,
		baseRepository: baseRepository,
	}
}

// RequestLoan process a loan request
func (u *_loan) RequestLoan(ctx context.Context, req *request.RequestLoan) (*response.RequestLoan, error) {
	var resp response.RequestLoan
	if err := u.baseRepository.ExecTx(ctx, u.requestLoan(ctx, req, &resp)); err != nil {
		return nil, errors.Wrap(err, "u.baseRepository.ExecTx() failed")
	}
	return &resp, nil
}

func (u *_loan) requestLoan(ctx context.Context, req *request.RequestLoan, resp *response.RequestLoan) repository.TxFn {
	return func(ctx context.Context, tx *sqlx.Tx) error {

		loanID, err := u.loanRepository.CreateLoan(ctx, &model.Loan{
			UserID:  req.UserID,
			Amount:  *req.Amount,
			Terms:   req.Terms,
			Status:  model.LoanStatusPending,
			Purpose: req.Purpose,
		}, tx)

		if err != nil {
			return errors.Wrap(err, "u.loanRepository.CreateLoan() failed")
		}

		resp.LoanID = loanID

		amount := req.Amount.Ceil()
		termInDecimal := decimal.NewFromInt(req.Terms)

		emi := amount.DivRound(termInDecimal, 2)

		adjustedLastEMI := emi
		if emiTotalAmount := emi.Mul(termInDecimal); !emiTotalAmount.Equal(amount) {
			adjustedLastEMI = adjustedLastEMI.Add(amount.Sub(emiTotalAmount))
		}

		today := time.Now()

		startDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

		loanEMIs := make([]*model.LoanEMI, req.Terms)
		for i := 0; i < int(req.Terms); i++ {
			dueDate := startDate.AddDate(0, 0, 7*(i+1))
			loanEMIs[i] = &model.LoanEMI{
				SeqNo:   uint64(i + 1),
				LoanID:  loanID,
				DueDate: dueDate,
				Amount:  emi,
				Status:  model.LoanStatusPending,
			}
		}

		if !emi.Equal(adjustedLastEMI) {
			loanEMIs[req.Terms-1].Amount = adjustedLastEMI
		}

		if err := u.loanRepository.CreateLoanEMIs(ctx, loanEMIs, tx); err != nil {
			return errors.Wrap(err, "u.loanRepository.CreateLoanEMIs() failed")
		}
		return nil
	}
}

// ApproveLoan changes status of a loan and its associated loan_emi entries's status to APPROVED
func (u *_loan) ApproveLoan(ctx context.Context, req *request.ApproveLoan) error {
	if err := u.baseRepository.ExecTx(ctx, u.approveLoan(ctx, req)); err != nil {
		return errors.Wrap(err, "u.baseRepository.ExecTx() failed")
	}
	return nil
}

func (u *_loan) approveLoan(ctx context.Context, req *request.ApproveLoan) repository.TxFn {
	return func(ctx context.Context, tx *sqlx.Tx) error {
		err := u.loanRepository.UpdateLoanStatus(ctx, req.LoanID, req.ApprovedBy, model.LoanStatusApproved, tx)
		if err != nil {
			return errors.Wrap(err, "u.loanRepository.UpdateLoanStatus() failed")
		}

		if err := u.loanRepository.UpdateLoanEMIStatusByLoanID(ctx, req.LoanID, model.LoanStatusApproved, tx); err != nil {
			return errors.Wrap(err, "u.loanRepository.UpdateLoanEMIStatusByLoanID() failed")
		}
		return nil
	}
}

// GetLoan retrieves loan details
func (u *_loan) GetLoan(ctx context.Context, req *request.GetLoan) (*response.GetLoan, error) {
	loanEMIs, err := u.loanRepository.GetLoanEMIs(ctx, req.Loan.ID)
	if err != nil {
		return nil, errors.Wrap(err, "u.loanRepository.GetLoanEMIs() failed")
	}
	return &response.GetLoan{
		Loan:     req.Loan,
		LoanEMIs: loanEMIs,
	}, nil
}

// RepayLoan repays the loan
func (u *_loan) RepayLoan(ctx context.Context, req *request.RepayLoan) error {
	loanEMIs, err := u.loanRepository.GetLoanEMIs(ctx, req.Loan.ID)
	if err != nil {
		return errors.Wrap(err, "u.loanRepository.GetLoanEMIs() failed")
	}

	availableAmount := *req.Amount
	fullyPaid := true
	emisToBeMarkedAsPaid := make([]uint64, 0, len(loanEMIs))
	for _, loanEMI := range loanEMIs {
		if loanEMI.Status == model.LoanStatusPaid {
			continue
		}
		if availableAmount.GreaterThanOrEqual(loanEMI.Amount) {
			emisToBeMarkedAsPaid = append(emisToBeMarkedAsPaid, loanEMI.ID)
			availableAmount = availableAmount.Sub(loanEMI.Amount)
			continue
		}
		fullyPaid = false
		break
	}
	if !availableAmount.Equal(decimal.Zero) {
		return apierr.ErrInvalidRepaymentAmount
	}

	if err := u.baseRepository.ExecTx(ctx, u.repayLoan(ctx, req.Loan, emisToBeMarkedAsPaid, fullyPaid)); err != nil {
		return errors.Wrap(err, "u.baseRepository.ExecTx() failed")
	}
	return nil
}

func (u *_loan) repayLoan(ctx context.Context, loan *model.Loan, loanEMIIds []uint64, fullyPaid bool) repository.TxFn {
	return func(ctx context.Context, tx *sqlx.Tx) error {
		if fullyPaid {
			if err := u.loanRepository.UpdateLoanStatus(ctx, loan.ID, *loan.ApprovedBy, model.LoanStatusPaid, tx); err != nil {
				return errors.Wrap(err, "u.loanRepository.UpdateLoanStatus() failed")
			}
		}

		if err := u.loanRepository.UpdateLoanEMIStatus(ctx, loan.ID, loanEMIIds, model.LoanStatusPaid, tx); err != nil {
			return errors.Wrap(err, "u.loanRepository.UpdateLoanEMIStatus() failed")
		}
		return nil
	}
}
