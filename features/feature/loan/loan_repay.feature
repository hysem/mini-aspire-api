Feature: Repaying a loan

    A loan can be repayed only by a customer who created it after it is approved by admin

    Background: A pending loan request exists
        Given "customer" 1 logged in
        When "customer" 1 request for a loan for an amount of 3000$, for 3 weeks, for the purpose of "bdd testing"
        Then loan request should "succeed"
        And "customer" 1 "can" view the loan request 
        And there should be a loan request in "PENDING" status for an amount of 3000$, for 3 weeks, for the purpose of "bdd testing"
        And there should be 3 loan_emi entries with the "amount" "1000, 1000, 1000" respectively.
        And there should be 3 loan_emi entries with the "status" "PENDING, PENDING, PENDING" respectively.

    Scenario: Customer can repay the loan for a single installment for a pending loan and failed
        Given "customer" 1 repays 1000 towards the loan and "failed"
        And "customer" 1 "can" view the loan request 

    Scenario: An approved loan request exists and customer can repay the loan for a single installment for an approved loan and succeed
        When "admin" 1 logged in 
        And "admin" 1 approves the loan request and "succeed"
        Then "customer" 1 "can" view the loan request 
        And there should be a loan request in "APPROVED" status for an amount of 3000$, for 3 weeks, for the purpose of "bdd testing"
        And there should be 3 loan_emi entries with the "amount" "1000, 1000, 1000" respectively.
        And there should be 3 loan_emi entries with the "status" "APPROVED, APPROVED, APPROVED" respectively.
        When "customer" 1 repays 1000 towards the loan and "succeed"
        When "customer" 1 "can" view the loan request 
        And there should be a loan request in "APPROVED" status for an amount of 3000$, for 3 weeks, for the purpose of "bdd testing"
        And there should be 3 loan_emi entries with the "amount" "1000, 1000, 1000" respectively.
        And there should be 3 loan_emi entries with the "status" "PAID, APPROVED, APPROVED" respectively.
    
    Scenario: An approved loan request exists and customer can repay the loan for a single installment for an approved loan and succeed
        When "admin" 1 logged in 
        And "admin" 1 approves the loan request and "succeed"
        Then "customer" 1 "can" view the loan request 
        And there should be a loan request in "APPROVED" status for an amount of 3000$, for 3 weeks, for the purpose of "bdd testing"
        And there should be 3 loan_emi entries with the "amount" "1000, 1000, 1000" respectively.
        And there should be 3 loan_emi entries with the "status" "APPROVED, APPROVED, APPROVED" respectively.
        When "customer" 1 repays 2000 towards the loan and "succeed"
        When "customer" 1 "can" view the loan request 
        And there should be a loan request in "APPROVED" status for an amount of 3000$, for 3 weeks, for the purpose of "bdd testing"
        And there should be 3 loan_emi entries with the "amount" "1000, 1000, 1000" respectively.
        And there should be 3 loan_emi entries with the "status" "PAID, PAID, APPROVED" respectively.
    
    Scenario: An approved loan request exists and customer can repay the loan for a three installment for an approved loan and succeed
        When "admin" 1 logged in 
        And "admin" 1 approves the loan request and "succeed"
        Then "customer" 1 "can" view the loan request 
        And there should be a loan request in "APPROVED" status for an amount of 3000$, for 3 weeks, for the purpose of "bdd testing"
        And there should be 3 loan_emi entries with the "amount" "1000, 1000, 1000" respectively.
        And there should be 3 loan_emi entries with the "status" "APPROVED, APPROVED, APPROVED" respectively.
        When "customer" 1 repays 3000 towards the loan and "succeed"
        When "customer" 1 "can" view the loan request 
        And there should be a loan request in "PAID" status for an amount of 3000$, for 3 weeks, for the purpose of "bdd testing"
        And there should be 3 loan_emi entries with the "amount" "1000, 1000, 1000" respectively.
        And there should be 3 loan_emi entries with the "status" "PAID, PAID, PAID" respectively.


        
