-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    name text NOT NULL UNIQUE,
    image_url text UNIQUE,
    balance bigint NOT NULL DEFAULT 0 CHECK (balance >= 0),
    trade_url text
);

CREATE TABLE users_items_statuses (
    id bigserial PRIMARY KEY,
    code text NOT NULL UNIQUE
);
INSERT INTO users_items_statuses (code)
VALUES
    ('active'),
    ('sold'),
    ('in_withdrawal'),
    ('withdrawn');
CREATE TABLE users_items (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users (id),
    item_id bigint NOT NULL REFERENCES items (id),
    attached_price bigint NOT NULL CHECK (attached_price >= 0),
    status_id bigint NOT NULL REFERENCES users_items_statuses (id)
);

-- +goose Down
DROP TABLE users_items;
DROP TABLE users_items_statuses;

DROP TABLE users;

