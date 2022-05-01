package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type (
	Loan struct {
		ID         uint64          `db:"loan_id" json:"loan_id"`
		UserID     uint64          `db:"user_id" json:"user_id"`
		Amount     decimal.Decimal `db:"amount" json:"amount"`
		Terms      int64           `db:"terms" json:"terms"`
		Status     LoanStatus      `db:"status" json:"status"`
		Purpose    string          `db:"purpose" json:"purpose"`
		ApprovedBy uint64          `db:"approved_by" json:"approved_by"`
		CreatedAt  time.Time       `db:"created_at" json:"created_at"`
		UpdatedAt  time.Time       `db:"updated_at" json:"updated_at"`
	}

	LoanEMI struct {
		ID        uint64          `db:"loan_emi_id" json:"loan_emi_id"`
		LoanID    uint64          `db:"loan_id" json:"loan_id"`
		SeqNo     uint64          `db:"seq_no" json:"seq_no"`
		DueDate   time.Time       `db:"due_date" json:"due_date"`
		Amount    decimal.Decimal `db:"amount" json:"amount"`
		Status    LoanStatus      `db:"status" json:"status"`
		CreatedAt time.Time       `db:"created_at" json:"created_at"`
		UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
	}
)
