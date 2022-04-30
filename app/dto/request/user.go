package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type (
	// UserGenerateToken request
	UserGenerateToken struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)

// Validate func
func (r *UserGenerateToken) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Email, validation.Required, validation.Length(3, 200), is.Email),
		validation.Field(&r.Password, validation.Required, validation.Length(8, 30)),
	)
}
