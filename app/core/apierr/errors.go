package apierr

import "errors"

// Err declarations
var (
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrInvalidRepaymentAmount = errors.New("amount should be multiple of emi amount or should be equal to the outstanding amount")
)
