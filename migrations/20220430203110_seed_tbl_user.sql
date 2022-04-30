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
('Consumer 1', 'consu1@yopmail.com', '$2a$12$f4fc/QNwMU7VvN2nSdPt/esCz/lIVj4je5H6vU6TcMgDzBUVgE67u', 'consumer'),
('Consumer 2', 'consu2@yopmail.com', '$2a$12$f4fc/QNwMU7VvN2nSdPt/esCz/lIVj4je5H6vU6TcMgDzBUVgE67u', 'consumer'),
('Consumer 3', 'consu3@yopmail.com', '$2a$12$f4fc/QNwMU7VvN2nSdPt/esCz/lIVj4je5H6vU6TcMgDzBUVgE67u', 'consumer'),
('Consumer 4', 'consu4@yopmail.com', '$2a$12$f4fc/QNwMU7VvN2nSdPt/esCz/lIVj4je5H6vU6TcMgDzBUVgE67u', 'consumer'),
('Consumer 5', 'consu5@yopmail.com', '$2a$12$f4fc/QNwMU7VvN2nSdPt/esCz/lIVj4je5H6vU6TcMgDzBUVgE67u', 'consumer');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "user" WHERE email IN (
'admin1@yopmail.com',
'admin2@yopmail.com',
'admin3@yopmail.com',
'admin4@yopmail.com',
'admin5@yopmail.com',
'consu1@yopmail.com',
'consu2@yopmail.com',
'consu3@yopmail.com',
'consu4@yopmail.com',
'consu5@yopmail.com'
);
-- +goose StatementEnd
