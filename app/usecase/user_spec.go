package usecase

import (
	"context"

	"github.com/hysem/mini-aspire-api/app/dto/request"
	"github.com/hysem/mini-aspire-api/app/dto/response"
)

// User interface declares supported operations on the user usecase
//go:generate mockery --name=User --filename=user_mock.go --inpackage

type User interface {
	// GenerateToken generates a token for an existing user
	GenerateToken(ctx context.Context, req *request.UserGenerateToken) (*response.UserGenerateToken, error)
}
