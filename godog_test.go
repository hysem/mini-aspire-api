package main_test

import (
	"github.com/cucumber/godog"
	"github.com/hysem/mini-aspire-api/features/featuretest"
)

func InitializeScenario(ctx *godog.ScenarioContext) {
	uc := featuretest.NewUserContext()
	lc := featuretest.NewLoanContext(uc)

	ctx.Step(`^"([^"]*)" (\d+) logged in$`, uc.DoLogin)
	ctx.Step(`^"([^"]*)" (\d+) request for a loan for an amount of (\d+)\$, for (\d+) weeks, for the purpose of "([^"]*)"$`, lc.RequestLoan)
	ctx.Step(`^loan request should "([^"]*)"$`, lc.CheckLoanRequestStatus)
	ctx.Step(`^"([^"]*)" (\d+) "([^"]*)" view the loan request in "([^"]*)" status for an amount of (\d+)\$, for (\d+) weeks, for the purpose of "([^"]*)"$`, lc.GetLoan)
	ctx.Step(`^there should be (\d+) loan_emi entries with the "([^"]*)" "([^"]*)" respectively\.$`, lc.VeryfyLoanEMI)
}
