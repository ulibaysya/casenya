BEGIN ISOLATION LEVEL REPEATABLE READ READ ONLY;
WITH tags AS (
    SELECT 
        r.item_id, 
        k.name AS key,
        v.name AS value
    FROM item_tag_relations AS r
    JOIN item_tag_values AS v ON r.item_tag_value_id = v.id 
    JOIN item_tag_keys AS k ON v.key_id = k.id 
)
SELECT 
    i.name, 
    i.image_url, 
    t.key, 
    t.value
FROM items AS i
LEFT JOIN tags AS t ON t.item_id = i.id
WHERE i.id = 4;

SELECT 
    pp.code AS item_price_provider,
    p.price
FROM items_prices AS p
JOIN items_price_providers AS pp ON pp.id = p.provider_id 
WHERE p.item_id = 4;
COMMIT;
