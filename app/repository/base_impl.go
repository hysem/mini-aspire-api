package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

// _base implements Base
type _base struct {
	masterDB *sqlx.DB
}

var _ Base = (*_base)(nil)

// NewBase returns a new _base repository instance
func NewBase(masterDB *sqlx.DB) *_base {
	return &_base{
		masterDB: masterDB,
	}
}

// ExecTx executes the given operations in a transaction
// This will also handle begin, rollback and commit transaction.
func (r *_base) ExecTx(ctx context.Context, fn TxFn) (err error) {
	var tx *sqlx.Tx

	tx, err = r.masterDB.BeginTxx(ctx, nil)
	if err != nil {
		err = errors.New("failed to begin transaction")
		return
	}

	defer func() {
		if err == nil {
			if errTx := tx.Commit(); errTx != nil {
				err = errors.Wrap(errTx, "failed to commit transaction")
			}
			return
		}

		if errTx := tx.Rollback(); errTx != nil {
			err = errors.Wrap(multierr.Append(err, errTx), "failed to rollback failed transation")
			return
		}
	}()

	if errTx := fn(ctx, tx); errTx != nil {
		err = errors.Wrap(errTx, "failed to execute transaction")
	}
	return
}
