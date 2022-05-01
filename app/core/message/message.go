package message

// Message constants
const (
	InternalServerError = "something went wrong"
	ParsingFailed       = "failed to parse request"
	ValidationFailed    = "failed to validate request"
	InvalidCredentials  = "invalid credentials"
	InvalidSession      = "invalid session"
	AccessDenied        = "you don't have access to this resource"
	LoanRequestSuccess  = "Loan request created successfully. Pending admin approval."
	ApprovedLoan        = "Loan approved"
	AlreadyApprovedLoan = "Loan is already approved"
	NoSuchResource      = "no such resource"
)
