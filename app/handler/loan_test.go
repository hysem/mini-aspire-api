package handler_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/hysem/mini-aspire-api/app/core/context"
	"github.com/hysem/mini-aspire-api/app/dto/request"
	"github.com/hysem/mini-aspire-api/app/dto/response"
	"github.com/hysem/mini-aspire-api/app/model"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Loan_RequestLoan(t *testing.T) {
	const validRequestBody = `{
		"amount": "10000",
		"terms": 3,
		"purpose": "test"
	}`
	authUser := &model.User{
		UserID: 1,
	}

	testCases := map[string]struct {
		body                 string
		setMocks             func(m *handlerMocks)
		expectedResponseBody string
		expectedStatusCode   int
	}{
		`error case: failed to parse request`: {
			body: `{`,
			setMocks: func(m *handlerMocks) {

			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"failed to parse request"}`,
		},
		`error case: failed to validate request`: {
			body: `{}`,
			setMocks: func(m *handlerMocks) {

			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error": {"terms":"cannot be blank","purpose":"cannot be blank"}, "message":"failed to validate request"}`,
		},
		`error case: failed to process loan request`: {
			body: validRequestBody,
			setMocks: func(m *handlerMocks) {
				m.loanUsecase.On("RequestLoan", mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
		`success case: processed loan request`: {
			body: validRequestBody,
			setMocks: func(m *handlerMocks) {
				m.loanUsecase.On("RequestLoan", mock.Anything, &request.RequestLoan{
					Amount:  decimal.NewFromInt(10000),
					Terms:   3,
					UserID:  1,
					Purpose: "test",
				}).Return(&response.RequestLoan{
					LoanID: 1,
				}, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"data":{"loan_id":1}, "message":"Loan request created successfully. Pending admin approval."}`,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			h, m := newHandler(t)
			defer m.assertExpectations(t)
			tc.setMocks(m)

			req, err := http.NewRequest(http.MethodPost, "/user/loan", bytes.NewBufferString(tc.body))
			req.Header.Add(echo.HeaderContentType, "application/json")
			assert.NoError(t, err)

			res := runHandler(t, req, h.loan.RequestLoan, func(cc *context.Context) {
				cc.AuthUser = authUser
			})
			assert.Equal(t, tc.expectedStatusCode, res.Result().StatusCode)
			assert.JSONEq(t, tc.expectedResponseBody, res.Body.String())
		})
	}
}

func TestHandler_Loan_ApproveLoan(t *testing.T) {
	authUser := &model.User{
		UserID: 1,
	}
	testCases := map[string]struct {
		loan                 *model.Loan
		setMocks             func(m *handlerMocks)
		expectedResponseBody string
		expectedStatusCode   int
	}{
		`success case: already approved loan`: {
			loan: &model.Loan{
				ID:     1,
				Status: model.LoanStatusApproved,
			},
			setMocks: func(m *handlerMocks) {
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message":"Loan is already approved"}`,
		},
		`error case: failed to approve loan`: {
			loan: &model.Loan{
				ID:     1,
				Status: model.LoanStatusPending,
			},
			setMocks: func(m *handlerMocks) {
				m.loanUsecase.On("ApproveLoan", mock.Anything, mock.Anything).Return(assert.AnError)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
		`success case: approved loan`: {
			loan: &model.Loan{
				ID:     1,
				Status: model.LoanStatusPending,
			},
			setMocks: func(m *handlerMocks) {
				m.loanUsecase.On("ApproveLoan", mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message":"Loan approved"}`,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			h, m := newHandler(t)
			defer m.assertExpectations(t)
			tc.setMocks(m)

			req, err := http.NewRequest(http.MethodPatch, "/user/loan/1/approve", nil)
			req.Header.Add(echo.HeaderContentType, "application/json")
			assert.NoError(t, err)

			res := runHandler(t, req, h.loan.ApproveLoan, func(cc *context.Context) {
				cc.AuthUser = authUser
				cc.Loan = tc.loan
			})

			assert.Equal(t, tc.expectedStatusCode, res.Result().StatusCode)
			assert.JSONEq(t, tc.expectedResponseBody, res.Body.String())
		})
	}
}

func TestHandler_Loan_GetLoan(t *testing.T) {
	loan := &model.Loan{
		ID:     1,
		Status: model.LoanStatusPending,
		UserID: 2,
	}
	testCases := map[string]struct {
		authUserID           uint64
		setMocks             func(m *handlerMocks)
		expectedResponseBody string
		expectedStatusCode   int
	}{
		`error case: failed to view loan`: {
			authUserID:           loan.UserID + 1,
			setMocks:             func(m *handlerMocks) {},
			expectedStatusCode:   http.StatusForbidden,
			expectedResponseBody: `{"message":"you don't have access to this resource"}`,
		},
		`error case: failed to get loan details`: {
			authUserID: loan.UserID,
			setMocks: func(m *handlerMocks) {
				m.loanUsecase.On("GetLoan", mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
		`success case: got loan details`: {
			authUserID: loan.UserID,
			setMocks: func(m *handlerMocks) {
				m.loanUsecase.On("GetLoan", mock.Anything, &request.GetLoan{
					Loan: loan,
				}).Return(&response.GetLoan{
					Loan:     loan,
					LoanEMIs: []*model.LoanEMI{{Status: model.LoanStatusPending}},
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{
				"data":{
					"loan":{
						"loan_id":1,
						"user_id":2,
						"amount":"0",
						"terms":0,
						"status":"PENDING",
						"purpose":"",
						"approved_by":null,
						"created_at":"0001-01-01T00:00:00Z",
						"updated_at":"0001-01-01T00:00:00Z"
					},
					"loan_emis":[{
						"loan_emi_id":0,
						"loan_id":0,
						"seq_no":0,
						"due_date":"0001-01-01T00:00:00Z",
						"amount":"0",
						"status":"PENDING",
						"created_at":"0001-01-01T00:00:00Z",
						"updated_at":"0001-01-01T00:00:00Z"
					}]
				}
			}`,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			h, m := newHandler(t)
			defer m.assertExpectations(t)
			tc.setMocks(m)

			req, err := http.NewRequest(http.MethodGet, "/user/loan/1", nil)
			req.Header.Add(echo.HeaderContentType, "application/json")
			assert.NoError(t, err)

			res := runHandler(t, req, h.loan.GetLoan, func(cc *context.Context) {
				cc.AuthUser = &model.User{
					UserID: tc.authUserID,
				}
				cc.Loan = loan
			})

			assert.Equal(t, tc.expectedStatusCode, res.Result().StatusCode)
			assert.JSONEq(t, tc.expectedResponseBody, res.Body.String())
		})
	}
}
