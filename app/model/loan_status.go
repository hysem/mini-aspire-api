package model

import "errors"

// Role
type LoanStatus string

const (
	LoanStatusInvalid  LoanStatus = ""
	LoanStatusPending  LoanStatus = "PENDING"
	LoanStatusApproved LoanStatus = "APPROVED"
	LoanStatusPaid     LoanStatus = "PAID"
)

var (
	loanStatuses = []LoanStatus{LoanStatusPending, LoanStatusApproved, LoanStatusPaid}
)

func (r *LoanStatus) UnmarshalText(b []byte) error {
	for _, v := range loanStatuses {
		if string(v) == string(b) {
			*r = v
			return nil
		}
	}
	return errors.New("invalid loan status")
}
