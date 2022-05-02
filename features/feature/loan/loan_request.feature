Feature: Loan request

    The customer can use this feature to create a loan request. 
    The request should contain the following fields: purpose, terms and amount.
    After successful validation an entry for the same will be created in loan, and the emi entries will be created in the loan_emi table.
    If the emi amount is a recurring decimal, say for eg: An amount of 100$ for a 3 week term will be divided into 33.33 emis.
    However there will be an error of .01 dollars which will be added to the last emi amount.

Scenario: As an admin I cannot request for a loan
    Given "admin" 10 logged in 
    When "admin" 1 request for a loan for an amount of 3000$, for 3 weeks, for the purpose of "bdd testing"
    Then loan request should "fail"

Scenario: As a customer I can request for a loan and emi amount is exactly computable
    Given "customer" 1 logged in
    When "customer" 1 request for a loan for an amount of 3000$, for 3 weeks, for the purpose of "bdd testing"
    Then loan request should "succeed"
    And "customer" 1 "can" view the loan request in "PENDING" status for an amount of 3000$, for 3 weeks, for the purpose of "bdd testing"
    And there should be 3 loan_emi entries with the "amount" "1000, 1000, 1000" respectively.
    And there should be 3 loan_emi entries with the "status" "PENDING, PENDING, PENDING" respectively.

Scenario: As a customer I can request for a loan and emi amount is not exactly computable
    Given "customer" 1 logged in
    When "customer" 1 request for a loan for an amount of 1000$, for 3 weeks, for the purpose of "bdd testing"
    Then loan request should "succeed"
    And "customer" 1 "can" view the loan request in "PENDING" status for an amount of 1000$, for 3 weeks, for the purpose of "bdd testing"
    And there should be 3 loan_emi entries with the "amount" "333.33, 333.33, 333.34" respectively.
    And there should be 3 loan_emi entries with the "status" "PENDING, PENDING, PENDING" respectively.
