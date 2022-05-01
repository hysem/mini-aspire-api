package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hysem/mini-aspire-api/app/model"
	"github.com/stretchr/testify/assert"
)

func TestRepositoryUser_GetByID(t *testing.T) {
	user := &model.User{
		UserID:   1,
		Name:     "name",
		Email:    "email@yopmail.com",
		Password: "password",
		Role:     model.RoleCustomer,
	}
	query := `SELECT user_id, name, email, password, role, created_at, updated_at FROM "user" WHERE user_id = $1`
	testCases := map[string]struct {
		setMocks         func(m *repositoryMocks)
		expectedErr      string
		expectedResponse *model.User
	}{
		`error case: failed to execute query`: {
			setMocks: func(m *repositoryMocks) {
				m.masterDBMock.ExpectQuery(query).WillReturnError(assert.AnError)
			},
			expectedErr: `failed to get user by id`,
		},
		`error case: no such user`: {
			setMocks: func(m *repositoryMocks) {
				m.masterDBMock.ExpectQuery(query).WithArgs(user.UserID).WillReturnError(sql.ErrNoRows)
			},
		},
		`success case: updated user`: {
			setMocks: func(m *repositoryMocks) {
				rows := sqlmock.NewRows([]string{"user_id", "name", "email", "password", "role", "created_at", "updated_at"})
				rows.AddRow(user.UserID, user.Name, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt)
				m.masterDBMock.ExpectQuery(query).WithArgs(user.UserID).WillReturnRows(rows)
			},
			expectedResponse: user,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			r, m := newRepository(t)
			defer m.assertExpectations(t)
			tc.setMocks(m)

			actualResponse, actualErr := r.user.GetByID(context.Background(), user.UserID)
			if tc.expectedErr == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.Contains(t, actualErr.Error(), tc.expectedErr)
			}
			assert.Equal(t, tc.expectedResponse, actualResponse)
		})
	}
}

func TestRepository_User_GetByEmail(t *testing.T) {
	user := &model.User{
		UserID:   1,
		Name:     "name",
		Email:    "email@yopmail.com",
		Password: "password",
		Role:     model.RoleCustomer,
	}
	query := `SELECT user_id, name, email, password, role, created_at, updated_at FROM "user" WHERE email = $1`
	testCases := map[string]struct {
		setMocks         func(m *repositoryMocks)
		expectedErr      string
		expectedResponse *model.User
	}{
		`error case: failed to execute query`: {
			setMocks: func(m *repositoryMocks) {
				m.masterDBMock.ExpectQuery(query).WillReturnError(assert.AnError)
			},
			expectedErr: `failed to get user by email`,
		},
		`error case: no such user`: {
			setMocks: func(m *repositoryMocks) {
				m.masterDBMock.ExpectQuery(query).WithArgs(user.Email).WillReturnError(sql.ErrNoRows)
			},
		},
		`success case: updated user`: {
			setMocks: func(m *repositoryMocks) {
				rows := sqlmock.NewRows([]string{"user_id", "name", "email", "password", "role", "created_at", "updated_at"})
				rows.AddRow(user.UserID, user.Name, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt)
				m.masterDBMock.ExpectQuery(query).WithArgs(user.Email).WillReturnRows(rows)
			},
			expectedResponse: user,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			r, m := newRepository(t)
			defer m.assertExpectations(t)
			tc.setMocks(m)

			actualResponse, actualErr := r.user.GetByEmail(context.Background(), user.Email)
			if tc.expectedErr == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.Contains(t, actualErr.Error(), tc.expectedErr)
			}
			assert.Equal(t, tc.expectedResponse, actualResponse)
		})
	}
}
