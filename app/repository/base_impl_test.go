package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRepository_Base_ExecTx(t *testing.T) {
	testCases := map[string]struct {
		setMocks    func(m *repositoryMocks)
		expectedErr string
	}{
		`error case: failed to begin transaction`: {
			setMocks: func(m *repositoryMocks) {
				m.masterDBMock.ExpectBegin().WillReturnError(assert.AnError)
			},
			expectedErr: `failed to begin transaction`,
		},
		`error case: failed txfn`: {
			setMocks: func(m *repositoryMocks) {
				m.masterDBMock.ExpectBegin()
				m.txFn.On("Execute", mock.Anything, mock.AnythingOfType("*sqlx.Tx")).Return(assert.AnError)
				m.masterDBMock.ExpectRollback()
			},
			expectedErr: `failed to execute transaction`,
		},
		`error case: failed txfn; failed rollback`: {
			setMocks: func(m *repositoryMocks) {
				m.masterDBMock.ExpectBegin()
				m.txFn.On("Execute", mock.Anything, mock.AnythingOfType("*sqlx.Tx")).Return(assert.AnError)
				m.masterDBMock.ExpectRollback().WillReturnError(assert.AnError)
			},
			expectedErr: `failed to rollback failed transation: failed to execute transaction`,
		},
		`error case: commit failed`: {
			setMocks: func(m *repositoryMocks) {
				m.masterDBMock.ExpectBegin()
				m.txFn.On("Execute", mock.Anything, mock.AnythingOfType("*sqlx.Tx")).Return(nil)
				m.masterDBMock.ExpectCommit().WillReturnError(assert.AnError)
			},
			expectedErr: `failed to commit transaction`,
		},
		`success case: done`: {
			setMocks: func(m *repositoryMocks) {
				m.masterDBMock.ExpectBegin()
				m.txFn.On("Execute", mock.Anything, mock.AnythingOfType("*sqlx.Tx")).Return(nil)
				m.masterDBMock.ExpectCommit()
			},
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			r, m := newRepository(t)
			defer m.assertExpectations(t)
			tc.setMocks(m)

			actualErr := r.base.ExecTx(context.Background(), m.txFn.Execute)
			if tc.expectedErr == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.Contains(t, actualErr.Error(), tc.expectedErr)
			}
		})
	}
}
