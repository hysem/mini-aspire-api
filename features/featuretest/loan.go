package featuretest

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

type LoanContext struct {
	uc         *UserContext
	loanID     uint64
	statusCode int
	loanDetail []byte
}

func NewLoanContext(uc *UserContext) *LoanContext {
	return &LoanContext{
		uc: uc,
	}
}

func (c *LoanContext) RequestLoan(who string, id int, amount string, terms int, purpose string) error {
	resp, err := httpClient.Post(endpointRequestLoan).
		AddHeader(echo.HeaderAuthorization, c.uc.getAuthHeader(who, id)).
		JSON(map[string]interface{}{
			"purpose": purpose,
			"amount":  amount,
			"terms":   terms,
		}).Send()
	if err != nil {
		return errors.Wrap(err, "failed to request loan")
	}
	c.statusCode = resp.StatusCode
	c.loanID = gjson.GetBytes(resp.Bytes(), "data.loan_id").Uint()
	return nil
}

func (c *LoanContext) VerifyLoanRequestStatus(expectedStatus string) error {
	if expectedStatus == "succeed" && c.statusCode == http.StatusCreated ||
		expectedStatus == "fail" && c.statusCode != http.StatusCreated {
		return nil
	}
	return fmt.Errorf("loan status check failed. expected: %s, got :%d", expectedStatus, c.statusCode)
}

func (c *LoanContext) GetLoan(who string, id int, expectedPermission string) error {
	resp, err := httpClient.Get(fmt.Sprintf(endpointGetLoan, c.loanID)).
		AddHeader(echo.HeaderAuthorization, c.uc.getAuthHeader(who, id)).
		Send()
	if err != nil {
		return errors.Wrap(err, "failed to get loan details")
	}
	if !(expectedPermission == "can" && resp.StatusCode == http.StatusOK ||
		expectedPermission == "can't" && resp.StatusCode == http.StatusForbidden) {
		return fmt.Errorf("%s view the loan details. got :%d", expectedPermission, resp.StatusCode)
	}

	c.loanDetail = resp.Bytes()

	return nil
}

func (c *LoanContext) VeryfyLoan(expectedStatus string, expectedAmount string, expectedTerms int64, expectedPurpose string) error {
	if actualStatus := gjson.GetBytes(c.loanDetail, "data.loan.status").String(); actualStatus != expectedStatus {
		return fmt.Errorf("expected loan status to be %s; got %s", expectedStatus, actualStatus)
	}
	if actualAmount := gjson.GetBytes(c.loanDetail, "data.loan.amount").String(); actualAmount != expectedAmount {
		return fmt.Errorf("expected loan status to be %s; got %s", expectedAmount, actualAmount)
	}
	if actualTerms := gjson.GetBytes(c.loanDetail, "data.loan.terms").Int(); actualTerms != expectedTerms {
		return fmt.Errorf("expected loan status to be %d; got %d", expectedTerms, actualTerms)
	}
	if actualPurpose := gjson.GetBytes(c.loanDetail, "data.loan.purpose").String(); actualPurpose != expectedPurpose {
		return fmt.Errorf("expected loan status to be %s; got %s", expectedPurpose, actualPurpose)
	}
	return nil
}

func (c *LoanContext) VeryfyLoanEMI(expectedCount int64, field string, expectedValues string) error {
	if actualCount := gjson.GetBytes(c.loanDetail, "data.loan_emis.#").Int(); actualCount != expectedCount {
		fmt.Println("not same")
		return fmt.Errorf("expected %d loan emi entries; got %d", expectedCount, actualCount)
	}
	if actualValues := quoteAndBracketReplacer.Replace(gjson.GetBytes(c.loanDetail, fmt.Sprintf("data.loan_emis.#.%s", field)).String()); actualValues != expectedValues {
		return fmt.Errorf("expected value of %s field to be %s; got %s", field, expectedValues, actualValues)
	}

	return nil
}

func (c *LoanContext) ApproveLoan(who string, id int, expectedStatus string) error {
	resp, err := httpClient.Patch(fmt.Sprintf(endpointApproveLoan, c.loanID)).
		AddHeader(echo.HeaderAuthorization, c.uc.getAuthHeader(who, id)).
		Send()
	if err != nil {
		return errors.Wrap(err, "failed to approve loan")
	}
	if !(expectedStatus == "succeed" && resp.StatusCode == http.StatusOK ||
		expectedStatus == "failed" && resp.StatusCode != http.StatusOK) {
		return fmt.Errorf("%s view the loan details. got :%d", expectedStatus, resp.StatusCode)
	}
	return nil
}

func (c *LoanContext) RepayLoan(who string, id int, amount string, expectedStatus string) error {
	resp, err := httpClient.Post(fmt.Sprintf(endpointRepayLoan, c.loanID)).
		AddHeader(echo.HeaderAuthorization, c.uc.getAuthHeader(who, id)).
		JSON(map[string]interface{}{
			"amount": amount,
		}).
		Send()
	if err != nil {
		return errors.Wrap(err, "failed to approve loan")
	}
	if !(expectedStatus == "succeed" && resp.StatusCode == http.StatusOK ||
		expectedStatus == "failed" && resp.StatusCode != http.StatusOK) {
		return fmt.Errorf("%s view the loan details. got :%d=> %s", expectedStatus, resp.StatusCode, resp.Bytes())
	}
	return nil
}
