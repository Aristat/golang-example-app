
-- +migrate Up

CREATE TABLE products(
    id SERIAL PRIMARY KEY
);

-- +migrate Down

DROP TABLE products;
