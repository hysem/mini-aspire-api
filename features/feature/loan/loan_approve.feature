Feature: Approving a loan request

    A loan request in pending status can be approved by an admin user

    Background: A loan request exists in pending status
        Given "customer" 1 logged in
        When "customer" 1 request for a loan for an amount of 3000$, for 3 weeks, for the purpose of "bdd testing"
        Then loan request should "succeed"
        And "customer" 1 "can" view the loan request in "PENDING" status for an amount of 3000$, for 3 weeks, for the purpose of "bdd testing"
        And there should be 3 loan_emi entries with the "amount" "1000, 1000, 1000" respectively.
        And there should be 3 loan_emi entries with the "status" "PENDING, PENDING, PENDING" respectively.

    Scenario: Customer shouldn't be able to approve the loan request
        Given "customer" 2 logged in
        When "customer" 2 approves the loan request and "failed"
        Then "customer" 1 "can" view the loan request in "PENDING" status for an amount of 3000$, for 3 weeks, for the purpose of "bdd testing"
        And there should be 3 loan_emi entries with the "amount" "1000, 1000, 1000" respectively.
        And there should be 3 loan_emi entries with the "status" "PENDING, PENDING, PENDING" respectively.
    
    Scenario: Admin should be able to approve the loan request
        Given "admin" 1 logged in
        When "admin" 1 approves the loan request and "succeed"
        Then "customer" 1 "can" view the loan request in "APPROVED" status for an amount of 3000$, for 3 weeks, for the purpose of "bdd testing"
        And there should be 3 loan_emi entries with the "amount" "1000, 1000, 1000" respectively.
        And there should be 3 loan_emi entries with the "status" "APPROVED, APPROVED, APPROVED" respectively.
