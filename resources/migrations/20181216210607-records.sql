
-- +migrate Up

-- password: 123456789
INSERT INTO users (email, encrypted_password) VALUES ('test@gmail.com', '$2a$10$bRGA2ckqEhydPBjA8jpDzehoAAhlmc95nGEB5WaZmxOz/EyUnE9l6');

-- +migrate Down
DELETE FROM users WHERE email = 'test@gmail.com';
