-- +goose Up

-- Add new column to items
ALTER TABLE items
    ADD COLUMN added_at timestamptz NOT NULL DEFAULT NOW();

-- Add is_sync_managed column to tags keys
ALTER TABLE item_tag_keys
    ADD COLUMN is_sync_managed boolean NOT NULL DEFAULT FALSE;

-- Add reserved (sync managed) keys and values
INSERT INTO item_tag_keys (name, is_public, is_sync_managed)
    VALUES 
        ('category', TRUE, TRUE),
        ('weapon_category', TRUE, TRUE),
        ('weapon_type', TRUE, TRUE),
        ('weapon_rarity', TRUE, TRUE),
        ('gloves_type', TRUE, TRUE),
        ('knife_type', TRUE, TRUE),
        ('pattern', TRUE, TRUE),
        ('wear', TRUE, TRUE),
        ('variant', TRUE, TRUE);
INSERT INTO item_tag_values (key_id, name)
SELECT k.id, v.value
FROM item_tag_keys AS k
JOIN (
    VALUES
        ('category', 'weapon'),
        ('category', 'knife'),
        ('category', 'gloves'),
        ('wear', 'unknown'),
        ('wear', 'fn'),
        ('wear', 'mw'),
        ('wear', 'ft'),
        ('wear', 'ww'),
        ('wear', 'bs'),
        ('variant', 'norm'),
        ('variant', 'st'),
        ('variant', 'sv')
) AS v(key, value) ON v.key = k.name;

-- Sync table
CREATE TABLE items_sync (
    item_id bigint PRIMARY KEY REFERENCES items (id) ON DELETE CASCADE,
    external_id text NOT NULL UNIQUE,
    last_synced_at timestamptz NOT NULL DEFAULT NOW(),
    hash BYTEA NOT NULL
);

-- +goose Down
ALTER TABLE items
    DROP COLUMN added_at;

WITH values_to_del AS (
    SELECT v.id
    FROM item_tag_values AS v
    JOIN item_tag_keys AS k ON v.key_id = k.id
    WHERE k.is_sync_managed
)
DELETE FROM item_tag_values AS v
    USING values_to_del AS del
    WHERE v.id = del.id;
DELETE FROM item_tag_keys WHERE is_sync_managed;
ALTER TABLE item_tag_keys
    DROP COLUMN is_sync_managed;

DROP TABLE items_sync;
