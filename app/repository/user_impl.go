package repository

import (
	"context"
	"database/sql"

	"github.com/hysem/mini-aspire-api/app/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var _ User = (*_user)(nil)

// _user implements User interface
type _user struct {
	masterDB *sqlx.DB
}

func NewUser(
	masterDB *sqlx.DB,
) *_user {
	return &_user{
		masterDB: masterDB,
	}
}

// GetByEmail retrieves an existing user by email
func (u *_user) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	const query = `SELECT user_id, name, email, password, role, created_at, updated_at FROM "user" WHERE email = $1`
	var user model.User
	if err := u.masterDB.GetContext(ctx, &user, query, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get user by email")
	}
	return &user, nil
}

// GetByID retrieves an existing user by id
func (u *_user) GetByID(ctx context.Context, id uint64) (*model.User, error) {
	const query = `SELECT user_id, name, email, password, role, created_at, updated_at FROM "user" WHERE user_id = $1`
	var user model.User
	if err := u.masterDB.GetContext(ctx, &user, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get user by id")
	}
	return &user, nil
}
