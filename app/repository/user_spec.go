package repository

import (
	"context"

	"github.com/hysem/mini-aspire-api/app/model"
)

// User interface declares supported operations on the user repository
//go:generate mockery --name=User --filename=user_mock.go --inpackage
type User interface {
	// GetByEmail retrieves an existing user by email
	GetByEmail(ctx context.Context, email string) (*model.User, error)

	// GetByID retrieves an existing user by id
	GetByID(ctx context.Context, id uint64) (*model.User, error)
}
