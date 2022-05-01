package usecase_test

import (
	"context"
	"testing"

	"github.com/hysem/mini-aspire-api/app/dto/request"
	"github.com/hysem/mini-aspire-api/app/dto/response"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUsecase_Loan_RequestLoan(t *testing.T) {
	request := &request.RequestLoan{
		Amount: decimal.NewFromInt(10000),
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
