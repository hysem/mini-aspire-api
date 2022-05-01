package usecase_test

import (
	"context"
	"testing"

	"github.com/hysem/mini-aspire-api/app/core/apierr"
	"github.com/hysem/mini-aspire-api/app/dto/request"
	"github.com/hysem/mini-aspire-api/app/dto/response"
	"github.com/hysem/mini-aspire-api/app/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUsecase_Loan_RequestLoan(t *testing.T) {
	amount := decimal.NewFromInt(10000)
	request := &request.RequestLoan{
		Amount: &amount,
		Terms:  3,
	}
	testCases := map[string]struct {
		setMocks         func(u *usecaseMocks)
		expectedErr      string
		expectedResponse *response.RequestLoan
	}{
		`error case: failed to execute transaction`: {
			setMocks: func(u *usecaseMocks) {
				u.baseRepository.On("ExecTx", mock.Anything, mock.Anything).Return(assert.AnError)
			},
			expectedErr: `u.baseRepository.ExecTx() failed`,
		},
		`success case: request created`: {
			setMocks: func(u *usecaseMocks) {
				u.baseRepository.On("ExecTx", mock.Anything, mock.Anything).Return(nil)
			},
			expectedResponse: &response.RequestLoan{
				LoanID: 0,
			},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			u, m := newUsecase(t)
			defer m.assertExpectations(t)
			tc.setMocks(m)

			actualResponse, actualErr := u.loan.RequestLoan(context.Background(), request)
			if tc.expectedErr == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.Contains(t, actualErr.Error(), tc.expectedErr)
			}

			assert.Equal(t, tc.expectedResponse, actualResponse)
		})
	}
}

func TestUsecase_Loan_ApproveLoan(t *testing.T) {
	request := &request.ApproveLoan{
		LoanID:     1,
		ApprovedBy: 2,
	}
	testCases := map[string]struct {
		setMocks    func(u *usecaseMocks)
		expectedErr string
	}{
		`error case: failed to execute transaction`: {
			setMocks: func(u *usecaseMocks) {
				u.baseRepository.On("ExecTx", mock.Anything, mock.Anything).Return(assert.AnError)
			},
			expectedErr: `u.baseRepository.ExecTx() failed`,
		},
		`success case: request created`: {
			setMocks: func(u *usecaseMocks) {
				u.baseRepository.On("ExecTx", mock.Anything, mock.Anything).Return(nil)
			},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			u, m := newUsecase(t)
			defer m.assertExpectations(t)
			tc.setMocks(m)

			actualErr := u.loan.ApproveLoan(context.Background(), request)
			if tc.expectedErr == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.Contains(t, actualErr.Error(), tc.expectedErr)
			}
		})
	}
}

func TestUsecase_Loan_GetLoan(t *testing.T) {
	request := &request.GetLoan{
		Loan: &model.Loan{ID: 1},
	}
	testCases := map[string]struct {
		setMocks         func(u *usecaseMocks)
		expectedErr      string
		expectedResponse *response.GetLoan
	}{
		`error case: failed to get loan emis`: {
			setMocks: func(u *usecaseMocks) {
				u.loanRepository.On("GetLoanEMIs", mock.Anything, request.Loan.ID).Return(nil, assert.AnError)
			},
			expectedErr: `u.loanRepository.GetLoanEMIs() failed`,
		},
		`success case: got loan details`: {
			setMocks: func(u *usecaseMocks) {
				u.loanRepository.On("GetLoanEMIs", mock.Anything, request.Loan.ID).Return([]*model.LoanEMI{{}}, nil)
			},
			expectedResponse: &response.GetLoan{
				Loan:     request.Loan,
				LoanEMIs: []*model.LoanEMI{{}},
			},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			u, m := newUsecase(t)
			defer m.assertExpectations(t)
			tc.setMocks(m)

			actualResponse, actualErr := u.loan.GetLoan(context.Background(), request)
			if tc.expectedErr == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.Contains(t, actualErr.Error(), tc.expectedErr)
			}
			assert.Equal(t, tc.expectedResponse, actualResponse)
		})
	}
}

func TestUsecase_Loan_RepayLoan(t *testing.T) {
	loan := &model.Loan{ID: 1}
	getRequest := func(amount int64) *request.RepayLoan {
		da := decimal.NewFromInt(amount)
		return &request.RepayLoan{
			Loan:   loan,
			Amount: &da,
		}
	}

	loanEMIs := []*model.LoanEMI{{
		ID:     1,
		SeqNo:  1,
		Amount: decimal.NewFromInt(2000),
		Status: model.LoanStatusPaid,
	}, {
		ID:     2,
		SeqNo:  2,
		Amount: decimal.NewFromInt(2000),
		Status: model.LoanStatusApproved,
	}, {
		ID:     3,
		SeqNo:  3,
		Amount: decimal.NewFromInt(2000),
		Status: model.LoanStatusApproved,
	}}

	testCases := map[string]struct {
		request     *request.RepayLoan
		setMocks    func(u *usecaseMocks)
		expectedErr string
	}{
		`error case: failed to get loan emis`: {
			request: getRequest(2000),
			setMocks: func(u *usecaseMocks) {
				u.loanRepository.On("GetLoanEMIs", mock.Anything, loan.ID).Return(nil, assert.AnError)
			},
			expectedErr: `u.loanRepository.GetLoanEMIs() failed`,
		},
		`error case: invalid repayment amount`: {
			request: getRequest(6000),
			setMocks: func(u *usecaseMocks) {
				u.loanRepository.On("GetLoanEMIs", mock.Anything, loan.ID).Return(loanEMIs, nil)
			},
			expectedErr: apierr.ErrInvalidRepaymentAmount.Error(),
		},
		`error case: failed to execute transaction`: {
			request: getRequest(2000),
			setMocks: func(u *usecaseMocks) {
				u.loanRepository.On("GetLoanEMIs", mock.Anything, loan.ID).Return(loanEMIs, nil)
				u.baseRepository.On("ExecTx", mock.Anything, mock.Anything).Return(assert.AnError)
			},
			expectedErr: `u.baseRepository.ExecTx() failed`,
		},
		`success case: paid`: {
			request: getRequest(2000),
			setMocks: func(u *usecaseMocks) {
				u.loanRepository.On("GetLoanEMIs", mock.Anything, loan.ID).Return(loanEMIs, nil)
				u.baseRepository.On("ExecTx", mock.Anything, mock.Anything).Return(nil)
			},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			u, m := newUsecase(t)
			defer m.assertExpectations(t)
			tc.setMocks(m)

			actualErr := u.loan.RepayLoan(context.Background(), tc.request)
			if tc.expectedErr == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.Contains(t, actualErr.Error(), tc.expectedErr)
			}
		})
	}
}
