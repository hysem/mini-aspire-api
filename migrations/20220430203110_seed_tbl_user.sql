-- +goose Up
-- +goose StatementBegin
INSERT INTO "user" 
("name", "email", "password", "role") 
VALUES 
('Admin 1', 'admin1@yopmail.com', '$2a$12$f4fc/QNwMU7VvN2nSdPt/esCz/lIVj4je5H6vU6TcMgDzBUVgE67u', 'admin'),
('Admin 2', 'admin2@yopmail.com', '$2a$12$f4fc/QNwMU7VvN2nSdPt/esCz/lIVj4je5H6vU6TcMgDzBUVgE67u', 'admin'),
('Admin 3', 'admin3@yopmail.com', '$2a$12$f4fc/QNwMU7VvN2nSdPt/esCz/lIVj4je5H6vU6TcMgDzBUVgE67u', 'admin'),
('Admin 4', 'admin4@yopmail.com', '$2a$12$f4fc/QNwMU7VvN2nSdPt/esCz/lIVj4je5H6vU6TcMgDzBUVgE67u', 'admin'),
('Admin 5', 'admin5@yopmail.com', '$2a$12$f4fc/QNwMU7VvN2nSdPt/esCz/lIVj4je5H6vU6TcMgDzBUVgE67u', 'admin'),
('Customer 1', 'cstmr1@yopmail.com', '$2a$12$f4fc/QNwMU7VvN2nSdPt/esCz/lIVj4je5H6vU6TcMgDzBUVgE67u', 'customer'),
('Customer 2', 'cstmr2@yopmail.com', '$2a$12$f4fc/QNwMU7VvN2nSdPt/esCz/lIVj4je5H6vU6TcMgDzBUVgE67u', 'customer'),
('Customer 3', 'cstmr3@yopmail.com', '$2a$12$f4fc/QNwMU7VvN2nSdPt/esCz/lIVj4je5H6vU6TcMgDzBUVgE67u', 'customer'),
('Customer 4', 'cstmr4@yopmail.com', '$2a$12$f4fc/QNwMU7VvN2nSdPt/esCz/lIVj4je5H6vU6TcMgDzBUVgE67u', 'customer'),
('Customer 5', 'cstmr5@yopmail.com', '$2a$12$f4fc/QNwMU7VvN2nSdPt/esCz/lIVj4je5H6vU6TcMgDzBUVgE67u', 'customer');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "user" WHERE email IN (
'admin1@yopmail.com',
'admin2@yopmail.com',
'admin3@yopmail.com',
'admin4@yopmail.com',
'admin5@yopmail.com',
'cstmr1@yopmail.com',
'cstmr2@yopmail.com',
'cstmr3@yopmail.com',
'cstmr4@yopmail.com',
'cstmr5@yopmail.com'
);
-- +goose StatementEnd
