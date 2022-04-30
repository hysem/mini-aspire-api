package usecase_test

import (
	"testing"

	"github.com/hysem/mini-aspire-api/app/core/bcrypt"
	"github.com/hysem/mini-aspire-api/app/core/jwt"
	"github.com/hysem/mini-aspire-api/app/repository"
	"github.com/hysem/mini-aspire-api/app/usecase"
)

type usecaseMocks struct {
	userRepository repository.MockUser
	jwt            jwt.MockJWT
	bcrypt         bcrypt.MockBcrypt
}

type testUsecase struct {
	user usecase.User
}

func newUsecase(t *testing.T) (*testUsecase, *usecaseMocks) {
	m := &usecaseMocks{}

	u := &testUsecase{
		user: usecase.NewUser(&m.userRepository, &m.bcrypt, &m.jwt),
	}
	return u, m
}

func (m *usecaseMocks) assertExpectations(t *testing.T) {
	m.userRepository.AssertExpectations(t)
	m.bcrypt.AssertExpectations(t)
	m.jwt.AssertExpectations(t)
}
