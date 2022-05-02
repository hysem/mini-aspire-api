# mini-aspire-api

## Assumptions/Tradeoffs
- For simplicity the consumer/admin users will be kept in a single table `user` along with `role` (i.e. whether `consumer` or `admin`).
- There is no user registration endpoint. If needed the users can be added to the db by using a client like dbeaver.
- For testing purposes, the following users are available in the db:
    | Email              | Password   | Role      |
    | ------------------ | ---------- | --------- |
    | admin1@yopmail.com | DWygs6wV   | Admin     |
    | admin2@yopmail.com | DWygs6wV   | Admin     |
    | admin3@yopmail.com | DWygs6wV   | Admin     |
    | admin4@yopmail.com | DWygs6wV   | Admin     |
    | admin5@yopmail.com | DWygs6wV   | Admin     |
    | cstmr1@yopmail.com | DWygs6wV   | Customer  |
    | consu2@yopmail.com | DWygs6wV   | Customer  |
    | cstmr3@yopmail.com | DWygs6wV   | Customer  |
    | cstmr4@yopmail.com | DWygs6wV   | Customer  |
    | cstmr5@yopmail.com | DWygs6wV   | Customer  |

## Requirements
- Golang 1.18
- Install necessary tools using `make tools`
## To run
 - Make a copy of config.env using `cp config.env config_local.env`
 - Start the containers using `make dep-up`
 - Run migrations using `make migrate-up`
 - Start the application using `make run`

## For running tests
### Unit test
- Mocks should be generated first using `make generate`
- Run the tests using `make test`
### Feature Tests
- Start the server `make run`
- Run the test `make ftest`

## Docs
- [Swagger Doc](http://localhost:8080/docs)