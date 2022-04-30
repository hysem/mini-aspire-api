package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/hysem/mini-aspire-api/app/core/apierr"
	"github.com/hysem/mini-aspire-api/app/dto/request"
	"github.com/hysem/mini-aspire-api/app/dto/response"
	"github.com/hysem/mini-aspire-api/app/model"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestUser_GenerateToken(t *testing.T) {
	request := &request.UserGenerateToken{
		Email:    "test@yopmail.com",
		Password: "password",
	}

	user := &model.User{
		UserID:    12,
		Name:      "Test User",
		Email:     request.Email,
		Role:      model.RoleConsumer,
		Password:  "hashed_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testCases := map[string]struct {
		setMocks         func(u *usecaseMocks)
		expectedErr      string
		expectedResponse *response.UserGenerateToken
	}{
		`error case: failed to get user by email`: {
			setMocks: func(u *usecaseMocks) {
				u.userRepository.On("GetByEmail", mock.Anything, request.Email).Return(nil, assert.AnError)
			},
			expectedErr: `GetByEmail() failed`,
		},
		`error case: invalid credentials; no such user`: {
			setMocks: func(u *usecaseMocks) {
				u.userRepository.On("GetByEmail", mock.Anything, request.Email).Return(nil, nil)
			},
			expectedErr: apierr.ErrInvalidCredentials.Error(),
		},
		`error case: invalid credentials; wrong password`: {
			setMocks: func(u *usecaseMocks) {
				u.userRepository.On("GetByEmail", mock.Anything, request.Email).Return(user, nil)
				u.bcrypt.On("Verify", user.Password, request.Password).Return(assert.AnError)
			},
			expectedErr: apierr.ErrInvalidCredentials.Error(),
		},
		`error case: failed to generate token`: {
			setMocks: func(u *usecaseMocks) {
				u.userRepository.On("GetByEmail", mock.Anything, request.Email).Return(user, nil)
				u.bcrypt.On("Verify", user.Password, request.Password).Return(nil)
				u.jwt.On("Generate", mock.Anything).Return("", assert.AnError)
			},
			expectedErr: `Generate() failed`,
		},
		`success case: generated token`: {
			setMocks: func(u *usecaseMocks) {
				u.userRepository.On("GetByEmail", mock.Anything, request.Email).Return(user, nil)
				u.bcrypt.On("Verify", user.Password, request.Password).Return(nil)
				u.jwt.On("Generate", user.UserID).Return("generated_token", nil)
			},
			expectedResponse: &response.UserGenerateToken{
				Token: "generated_token",
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

			actualResponse, actualErr := u.user.GenerateToken(context.Background(), request)
			if tc.expectedErr == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.Contains(t, actualErr.Error(), tc.expectedErr)
			}

			assert.Equal(t, tc.expectedResponse, actualResponse)
		})
	}
}
