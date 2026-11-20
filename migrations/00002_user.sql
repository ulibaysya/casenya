-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    image_url TEXT UNIQUE,
    balance BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    trade_url TEXT
);

CREATE TABLE users_items_statuses (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE
);
INSERT INTO users_items_statuses (code) VALUES 
    ('active'),
    ('sold'),
    ('in_withdrawal'),
    ('withdrawn');
CREATE TABLE users_items (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    item_id BIGINT NOT NULL REFERENCES items(id),
    attached_price BIGINT NOT NULL CHECK (attached_price >= 0),
    status_id BIGINT NOT NULL REFERENCES users_items_statuses(id)
);

CREATE TABLE IF NOT EXISTS users_opaque_tokens (
    user_id BIGINT NOT NULL REFERENCES users(id),
    token TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users_opaque_tokens;

DROP TABLE users_items;
DROP TABLE users_items_statuses;

DROP TABLE users;
