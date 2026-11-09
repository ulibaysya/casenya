
WITH values_ids AS (
    SELECT v.id AS id
    FROM item_tag_values AS v
    JOIN item_tag_keys AS k ON k.id = v.key_id
    JOIN (
        VALUES
--           key  value
            ('1', '1.2'),
            ('4', '4.3')
    ) AS n(key_name, value_name) 
    ON n.key_name = k.name AND n.value_name = v.name
)
INSERT INTO item_tag_relations
SELECT 1, v.id
FROM values_ids AS v;

INSERT INTO items_prices (item_id, price, provider_id)
SELECT 1, 12345, p.id 
FROM items_price_providers AS p 
WHERE code = 'custom';
