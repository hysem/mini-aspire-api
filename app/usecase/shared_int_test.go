package usecase

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hysem/mini-aspire-api/app/core/bcrypt"
	"github.com/hysem/mini-aspire-api/app/core/jwt"
	"github.com/hysem/mini-aspire-api/app/repository"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type usecaseMocks struct {
	userRepository repository.MockUser
	loanRepository repository.MockLoan
	baseRepository repository.MockBase
	jwt            jwt.MockJWT
	bcrypt         bcrypt.MockBcrypt
	masterDB       *sqlx.DB
	masterDBMock   sqlmock.Sqlmock
}

type testUsecase struct {
	user *_user
	loan *_loan
}

func newUsecase(t *testing.T) (*testUsecase, *usecaseMocks) {
	m := &usecaseMocks{}

	if db, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(queryMatcher())); err != nil {
		assert.NoError(t, err)
	} else {
		m.masterDB = sqlx.NewDb(db, "sqlmock")
		m.masterDBMock = dbMock
	}

	u := &testUsecase{
		user: NewUser(&m.userRepository, &m.bcrypt, &m.jwt),
		loan: NewLoan(&m.loanRepository, &m.baseRepository),
	}
	return u, m
}

func (m *usecaseMocks) assertExpectations(t *testing.T) {
	m.userRepository.AssertExpectations(t)
	m.baseRepository.AssertExpectations(t)
	m.loanRepository.AssertExpectations(t)
	m.bcrypt.AssertExpectations(t)
	m.jwt.AssertExpectations(t)
	assert.NoError(t, m.masterDBMock.ExpectationsWereMet())
}

func queryMatcher() sqlmock.QueryMatcherFunc {
	return sqlmock.QueryMatcherFunc(func(expectedSQL, actualSQL string) error {
		expectedSQL = regexp.QuoteMeta(expectedSQL)
		return sqlmock.QueryMatcherRegexp.Match(expectedSQL, actualSQL)
	})
}
