package usecase

import (
	"context"

	"github.com/hysem/mini-aspire-api/app/core/apierr"
	"github.com/hysem/mini-aspire-api/app/core/bcrypt"
	"github.com/hysem/mini-aspire-api/app/core/jwt"
	"github.com/hysem/mini-aspire-api/app/dto/request"
	"github.com/hysem/mini-aspire-api/app/dto/response"
	"github.com/hysem/mini-aspire-api/app/repository"
	"github.com/pkg/errors"
)

var _ User = (*_user)(nil)

func NewUser(
	userRepository repository.User,
	bcrypt bcrypt.Bcrypt,
	jwt jwt.JWT,
) *_user {
	return &_user{
		userRepository: userRepository,
		bcrypt:         bcrypt,
		jwt:            jwt,
	}
}

// _user implements the User interface
type _user struct {
	userRepository repository.User
	bcrypt         bcrypt.Bcrypt
	jwt            jwt.JWT
}

// GenerateToken generates a token for an existing user
func (u *_user) GenerateToken(ctx context.Context, req *request.UserGenerateToken) (*response.UserGenerateToken, error) {
	existingUser, err := u.userRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.Wrap(err, "u.userRepository.GetByEmail() failed")
	}

	if existingUser == nil {
		return nil, apierr.ErrInvalidCredentials
	}

	// verify the password
	if err := u.bcrypt.Verify(existingUser.Password, req.Password); err != nil {
		return nil, apierr.ErrInvalidCredentials
	}

	token, err := u.jwt.Generate(existingUser.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "u.jwt.Generate() failed")
	}

	return &response.UserGenerateToken{
		Token: token,
	}, nil
}
