package usecase

import (
	"context"
	"time"

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
			Amount:  req.Amount,
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
