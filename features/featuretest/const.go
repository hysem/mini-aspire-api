package featuretest

import "strings"

const (
	password = "DWygs6wV"

	endpointToken       = "/v1/user/token"
	endpointRequestLoan = "/v1/user/loan"
	endpointGetLoan     = "/v1/user/loan/%d"
	endpointApproveLoan = "/v1/user/loan/%d/approve"
)

var (
	quoteAndBracketReplacer = strings.NewReplacer(`["`, "", `","`, `, `, `"]`, "")
)
