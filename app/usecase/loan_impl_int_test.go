package usecase

import (
	"context"
	"testing"

	"github.com/hysem/mini-aspire-api/app/dto/request"
	"github.com/hysem/mini-aspire-api/app/dto/response"
	"github.com/hysem/mini-aspire-api/app/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUsecaseLoan_requestLoan(t *testing.T) {
	request := &request.RequestLoan{
		Amount:  decimal.NewFromInt(10000),
		Terms:   3,
		UserID:  1,
		Purpose: "purpose",
	}
	loanID := uint64(1)
	testCases := map[string]struct {
		setMocks         func(u *usecaseMocks)
		expectedErr      string
		expectedResponse response.RequestLoan
	}{
		`error case: failed to create loan`: {
			setMocks: func(u *usecaseMocks) {
				u.masterDBMock.ExpectBegin()
				u.loanRepository.On("CreateLoan", mock.Anything, mock.Anything, mock.Anything).Return(uint64(0), assert.AnError)
			},
			expectedErr: `u.loanRepository.CreateLoan() failed`,
		},
		`error case: failed to create loan emi`: {
			setMocks: func(u *usecaseMocks) {
				u.masterDBMock.ExpectBegin()
				u.loanRepository.On("CreateLoan", mock.Anything, &model.Loan{
					UserID:  request.UserID,
					Amount:  request.Amount,
					Terms:   request.Terms,
					Status:  model.LoanStatusPending,
					Purpose: request.Purpose,
				}, mock.Anything).Return(loanID, nil)
				u.loanRepository.On("CreateLoanEMIs", mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError)
			},
			expectedResponse: response.RequestLoan{LoanID: loanID},
			expectedErr:      `u.loanRepository.CreateLoanEMIs() failed`,
		},
		`success case: requested loan`: {
			setMocks: func(u *usecaseMocks) {
				u.masterDBMock.ExpectBegin()
				u.loanRepository.On("CreateLoan", mock.Anything, &model.Loan{
					UserID:  request.UserID,
					Amount:  request.Amount,
					Terms:   request.Terms,
					Status:  model.LoanStatusPending,
					Purpose: request.Purpose,
				}, mock.Anything).Return(loanID, nil)
				u.loanRepository.On("CreateLoanEMIs", mock.Anything, mock.MatchedBy(func(v []*model.LoanEMI) bool {
					return v[0].Amount.Equal(decimal.RequireFromString("3333.33")) &&
						v[1].Amount.Equal(decimal.RequireFromString("3333.33")) &&
						v[2].Amount.Equal(decimal.RequireFromString("3333.34"))
				}), mock.Anything).Return(nil)
			},
			expectedResponse: response.RequestLoan{LoanID: loanID},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			u, m := newUsecase(t)
			defer m.assertExpectations(t)
			tc.setMocks(m)

			tx, err := m.masterDB.Beginx()
			assert.NoError(t, err)

			var actualResponse response.RequestLoan
			actualErr := u.loan.requestLoan(context.Background(), request, &actualResponse)(context.Background(), tx)
			if tc.expectedErr == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.Contains(t, actualErr.Error(), tc.expectedErr)
			}

			assert.Equal(t, tc.expectedResponse, actualResponse)
		})
	}
}
