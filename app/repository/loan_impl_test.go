package repository_test

import (
	"context"
	"database/sql"
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
	query := `INSERT INTO loan (purpose, amount, terms, status, user_id, approved_by, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING loan_id`
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

func TestRepository_Loan_GetLoanByID(t *testing.T) {
	loan := &model.Loan{
		ID:      10,
		UserID:  1,
		Amount:  decimal.NewFromInt(10000),
		Terms:   3,
		Status:  model.LoanStatusPending,
		Purpose: "purpose",
	}
	query := `SELECT 
		loan_id, purpose, amount, terms, status, user_id, approved_by, created_at, updated_at 
	FROM loan WHERE loan_id=$1`
	testCases := map[string]struct {
		setMocks         func(m *repositoryMocks)
		expectedErr      string
		expectedResponse *model.Loan
	}{
		`error case: failed to execute query`: {
			setMocks: func(m *repositoryMocks) {
				m.masterDBMock.ExpectQuery(query).WillReturnError(assert.AnError)
			},
			expectedErr: `failed to create loan`,
		},
		`error case: no such rows`: {
			setMocks: func(m *repositoryMocks) {
				m.masterDBMock.ExpectQuery(query).WillReturnError(sql.ErrNoRows)
			},
		},
		`success case: got loan details`: {
			setMocks: func(m *repositoryMocks) {
				rows := sqlmock.NewRows([]string{"loan_id", "purpose", "amount", "terms", "status", "user_id", "approved_by", "created_at", "updated_at"})
				rows.AddRow(loan.ID, loan.Purpose, loan.Amount, loan.Terms, loan.Status, loan.UserID, loan.ApprovedBy, loan.CreatedAt, loan.UpdatedAt)
				m.masterDBMock.ExpectQuery(query).
					WithArgs(loan.ID).
					WillReturnRows(rows)
			},
			expectedResponse: loan,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			r, m := newRepository(t)
			defer m.assertExpectations(t)
			tc.setMocks(m)

			actualResponse, actualErr := r.loan.GetLoanByID(context.Background(), loan.ID)
			if tc.expectedErr == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.Contains(t, actualErr.Error(), tc.expectedErr)
			}
			assert.Equal(t, tc.expectedResponse, actualResponse)
		})
	}
}

func TestLoan_UpdateLoanStatus(t *testing.T) {
	loanID := uint64(1)
	approvedBy := uint64(2)
	status := model.LoanStatusApproved

	query := `UPDATE loan SET status=$1, approved_by=$2 WHERE loan_id=$3`
	testCases := map[string]struct {
		setMocks    func(m *repositoryMocks)
		expectedErr string
	}{
		`error case: failed to execute query`: {
			setMocks: func(m *repositoryMocks) {
				m.masterDBMock.ExpectBegin()
				m.masterDBMock.ExpectExec(query).WillReturnError(assert.AnError)
			},
			expectedErr: `failed to update loan status`,
		},
		`success case: updated loan status`: {
			setMocks: func(m *repositoryMocks) {
				m.masterDBMock.ExpectBegin()
				m.masterDBMock.ExpectExec(query).WithArgs(status, approvedBy, loanID).WillReturnResult(sqlmock.NewResult(0, 1))
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

			actualErr := r.loan.UpdateLoanStatus(context.Background(), loanID, approvedBy, status, tx)
			if tc.expectedErr == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.Contains(t, actualErr.Error(), tc.expectedErr)
			}
		})
	}
}

func TestLoan_UpdateLoanEMIStatusByLoanID(t *testing.T) {
	loanID := uint64(1)
	status := model.LoanStatusApproved

	query := `UPDATE loan_emi SET status=$1 WHERE loan_id=$2`
	testCases := map[string]struct {
		setMocks    func(m *repositoryMocks)
		expectedErr string
	}{
		`error case: failed to execute query`: {
			setMocks: func(m *repositoryMocks) {
				m.masterDBMock.ExpectBegin()
				m.masterDBMock.ExpectExec(query).WillReturnError(assert.AnError)
			},
			expectedErr: `failed to update loan_emi status`,
		},
		`success case: updated loan status`: {
			setMocks: func(m *repositoryMocks) {
				m.masterDBMock.ExpectBegin()
				m.masterDBMock.ExpectExec(query).WithArgs(status, loanID).WillReturnResult(sqlmock.NewResult(0, 1))
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

			actualErr := r.loan.UpdateLoanEMIStatusByLoanID(context.Background(), loanID, status, tx)
			if tc.expectedErr == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.Contains(t, actualErr.Error(), tc.expectedErr)
			}
		})
	}
}
