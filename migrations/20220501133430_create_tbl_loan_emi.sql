-- +goose Up
-- +goose StatementBegin
CREATE TABLE "loan_emi"(
    "loan_emi_id" SERIAL NOT NULL,
    "loan_id" INT NOT NULL,
    "seq_no" INT NOT NULL,
    "due_date" DATE NOT NULL,
    "amount" TEXT NOT NULL,
    "status" VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "pk_loan_emi_loan_emi_id" PRIMARY KEY("loan_emi_id"),
    CONSTRAINT "fk_loan_emi_loan_id" FOREIGN KEY("loan_id") REFERENCES "loan"("loan_id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "loan";
-- +goose StatementEnd
