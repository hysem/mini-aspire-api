package repository_test

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hysem/mini-aspire-api/app/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Loan_CreateLoan(t *testing.T) {
	loan := &model.Loan{
		ID:      10,
		UserID:  1,
		Amount:  decimal.NewFromInt(10000),
		Terms:   3,
		Status:  model.LoanStatusPending,
		Purpose: "purpose",
	}
	query := `INSERT INTO loan (purpose, amount, terms, status, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING loan_id`
	testCases := map[string]struct {
		setMocks         func(m *repositoryMocks)
		expectedErr      string
		expectedResponse uint64
	}{
		`error case: failed to execute query`: {
			setMocks: func(m *repositoryMocks) {
				m.masterDBMock.ExpectBegin()
				m.masterDBMock.ExpectQuery(query).WillReturnError(assert.AnError)
			},
			expectedErr: `failed to create loan`,
		},
		`success case: created loan`: {
			setMocks: func(m *repositoryMocks) {
				m.masterDBMock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"loan_id"})
				rows.AddRow(loan.ID)
				m.masterDBMock.ExpectQuery(query).
					WithArgs(loan.Purpose, loan.Amount, loan.Terms, loan.Status, loan.UserID).
					WillReturnRows(rows)
			},
			expectedResponse: loan.ID,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			r, m := newRepository(t)
			defer m.assertExpectations(t)
			tc.setMocks(m)

			tx, err := m.masterDB.Beginx()
			assert.NoError(t, err)

			actualResponse, actualErr := r.loan.CreateLoan(context.Background(), loan, tx)
			if tc.expectedErr == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.Contains(t, actualErr.Error(), tc.expectedErr)
			}
			assert.Equal(t, tc.expectedResponse, actualResponse)
		})
	}
}

func TestLoan_CreateLoanEMIs(t *testing.T) {
	loanID := uint64(1)
	loanEMIs := []*model.LoanEMI{{
		SeqNo:  1,
		LoanID: loanID,
	}, {
		SeqNo: 2,
	}}

	query := `INSERT INTO loan_emi (loan_id, seq_no, due_date, amount, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	testCases := map[string]struct {
		setMocks    func(m *repositoryMocks)
		expectedErr string
	}{
		`error case: failed to execute query`: {
			setMocks: func(m *repositoryMocks) {
				m.masterDBMock.ExpectBegin()
				m.masterDBMock.ExpectExec(query).WillReturnError(assert.AnError)
			},
			expectedErr: `failed to create loan`,
		},
		`success case: created loan emis`: {
			setMocks: func(m *repositoryMocks) {
				m.masterDBMock.ExpectBegin()
				var values []driver.Value
				for _, loanEMI := range loanEMIs {
					values = append(values, loanEMI.LoanID, loanEMI.SeqNo, loanEMI.DueDate, loanEMI.Amount, loanEMI.Status)
				}
				m.masterDBMock.ExpectExec(query).WithArgs(values...).WillReturnResult(sqlmock.NewResult(0, int64(len(loanEMIs))))
			},
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			r, m := newRepository(t)
			defer m.assertExpectations(t)
			tc.setMocks(m)

			tx, err := m.masterDB.Beginx()
			assert.NoError(t, err)

			actualErr := r.loan.CreateLoanEMIs(context.Background(), loanEMIs, tx)
			if tc.expectedErr == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.Contains(t, actualErr.Error(), tc.expectedErr)
			}
		})
	}
}
