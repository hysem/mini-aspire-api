package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/shopspring/decimal"
)

type (
	RequestLoan struct {
		Amount  decimal.Decimal `json:"amount"`
		Terms   int64           `json:"terms"`
		Purpose string          `json:"purpose"`
		UserID  uint64          `json:"-"`
	}
)

// Validate func
func (r *RequestLoan) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Amount, validation.Required), //TODO: custom validation
		validation.Field(&r.Purpose, validation.Required, validation.Length(1, 200)),
		validation.Field(&r.Terms, validation.Required, validation.Min(int64(1)), validation.Max(int64(100))),
	)
}
