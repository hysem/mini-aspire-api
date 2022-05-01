package message

// Message constants
const (
	InternalServerError = "something went wrong"
	ParsingFailed       = "failed to parse request"
	ValidationFailed    = "failed to validate request"
	InvalidCredentials  = "invalid credentials"
	InvalidSession      = "invalid session"
	AccessDenied        = "you don't have access to this resource"
	LoanRequestSuccess  = "loan request created"
	ApprovedLoan        = "loan approved"
	NotYetApprovedLoan  = "loan not yet approved"
	AlreadyApprovedLoan = "loan is already approved"
	AlreadyPaid         = "loan is already paid"
	PaymentSuccess      = "payment success"
	NoSuchResource      = "no such resource"
)
