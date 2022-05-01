package repository_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hysem/mini-aspire-api/app/repository"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type repositoryMocks struct {
	masterDB     *sqlx.DB
	masterDBMock sqlmock.Sqlmock
	txFn         repository.MockTxFn
}

type testRepository struct {
	user repository.User
	loan repository.Loan
	base repository.Base
}

func queryMatcher() sqlmock.QueryMatcherFunc {
	return sqlmock.QueryMatcherFunc(func(expectedSQL, actualSQL string) error {
		expectedSQL = regexp.QuoteMeta(expectedSQL)
		return sqlmock.QueryMatcherRegexp.Match(expectedSQL, actualSQL)
	})
}

func newRepository(t *testing.T) (*testRepository, *repositoryMocks) {
	m := &repositoryMocks{}

	if db, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(queryMatcher())); err != nil {
		assert.NoError(t, err)
	} else {
		m.masterDB = sqlx.NewDb(db, "sqlmock")
		m.masterDBMock = dbMock
	}

	u := &testRepository{
		user: repository.NewUser(m.masterDB),
		base: repository.NewBase(m.masterDB),
		loan: repository.NewLoan(m.masterDB),
	}
	return u, m
}

func (m *repositoryMocks) assertExpectations(t *testing.T) {
	assert.NoError(t, m.masterDBMock.ExpectationsWereMet())
	m.txFn.AssertExpectations(t)
}
