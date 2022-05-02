Feature: Viewing a loan request

    A loan request can be viewed only by the customer who created it.

    Background: A loan request exists
        Given "customer" 1 logged in
        When "customer" 1 request for a loan for an amount of 3000$, for 3 weeks, for the purpose of "bdd testing"
        Then loan request should "succeed"
        And "customer" 1 "can" view the loan request in "PENDING" status for an amount of 3000$, for 3 weeks, for the purpose of "bdd testing"
        And there should be 3 loan_emi entries with the "amount" "1000, 1000, 1000" respectively.
        And there should be 3 loan_emi entries with the "status" "PENDING, PENDING, PENDING" respectively.

    Scenario: Another customer should not be able to view the request
        Given "customer" 2 logged in
        Then "customer" 2 "can't" view the loan request
        
    
    Scenario: Admin should be able to approve the loan request
        Given "admin" 1 logged in
        Then "admin" 1 "can't" view the loan request
        
