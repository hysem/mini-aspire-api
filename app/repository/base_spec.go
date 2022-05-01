package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Base repository
//go:generate mockery --name=Base --filename=base_mock.go --inpackage
type Base interface {
	// ExecTx executes the given operations in a transaction
	// This will also handle begin, rollback and commit transaction.
	ExecTx(ctx context.Context, fn TxFn) error
}

// TxFn type
//go:generate mockery --name=TxFn --filename=base_tx_fn_mock.go --inpackage
type TxFn func(ctx context.Context, tx *sqlx.Tx) error
