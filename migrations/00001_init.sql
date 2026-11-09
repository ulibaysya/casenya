-- +goose Up
CREATE TABLE IF NOT EXISTS items (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    image_url TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS item_tag_keys (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    is_public BOOL NOT NULL DEFAULT FALSE
);
CREATE TABLE IF NOT EXISTS item_tag_values (
    id BIGSERIAL PRIMARY KEY,
    key_id BIGINT NOT NULL REFERENCES item_tag_keys(id),
    name TEXT NOT NULL,
    UNIQUE (key_id, name)
);
CREATE TABLE IF NOT EXISTS item_tag_relations (
    item_id BIGINT NOT NULL REFERENCES items(id),
    item_tag_value_id BIGINT NOT NULL REFERENCES item_tag_values(id),
    PRIMARY KEY (item_id, item_tag_value_id)
);

CREATE TABLE IF NOT EXISTS items_price_providers (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE
);
INSERT INTO items_price_providers(code) VALUES ('custom');
CREATE TABLE IF NOT EXISTS items_prices (
    item_id BIGINT NOT NULL REFERENCES items(id),
    price BIGINT NOT NULL CHECK (price >= 0),
    provider_id BIGINT NOT NULL REFERENCES items_price_providers(id),
    PRIMARY KEY (item_id, provider_id)
);

-- +goose Down

DROP TABLE items_prices;
DROP TABLE items_price_providers;

DROP TABLE item_tag_relations;
DROP TABLE item_tag_values;
DROP TABLE item_tag_keys;

DROP TABLE items;
