package response

import "github.com/hysem/mini-aspire-api/app/model"

type (
	// RequestLoan response
	RequestLoan struct {
		LoanID uint64 `json:"loan_id"`
	}
	GetLoan struct {
		Loan     *model.Loan      `json:"loan"`
		LoanEMIs []*model.LoanEMI `json:"loan_emis"`
	}
)
