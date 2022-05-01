-- +goose Up
-- +goose StatementBegin
CREATE TABLE "loan"(
    "loan_id" SERIAL NOT NULL,
    "user_id" INT NOT NULL,
    "purpose" VARCHAR(200) NOT NULL,
    "amount" TEXT NOT NULL,
    "terms" INT NOT NULL,
    "status" VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    "approved_by" INT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "pk_loan_loan_id" PRIMARY KEY("loan_id"),
    CONSTRAINT "fk_loan_user_id" FOREIGN KEY("user_id") REFERENCES "user"("user_id"),
    CONSTRAINT "fk_loan_approved_by" FOREIGN KEY("approved_by") REFERENCES "user"("user_id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "loan";
-- +goose StatementEnd
