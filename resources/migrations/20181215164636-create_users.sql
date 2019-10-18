
-- +migrate Up

CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL DEFAULT '',
    encrypted_password VARCHAR(255) NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX users_email ON users USING btree (email);

-- +migrate Down
