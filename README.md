# mini-aspire-api

# Assumptions/Tradeoffs
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
    | consu1@yopmail.com | DWygs6wV   | Consumer  |
    | consu2@yopmail.com | DWygs6wV   | Consumer  |
    | consu3@yopmail.com | DWygs6wV   | Consumer  |
    | consu4@yopmail.com | DWygs6wV   | Consumer  |
    | consu5@yopmail.com | DWygs6wV   | Consumer  |