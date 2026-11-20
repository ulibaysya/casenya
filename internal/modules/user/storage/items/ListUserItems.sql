SELECT
    i.id,
    i.item_id,
    i.attached_price,
    s.code
FROM users_items AS i
JOIN users_items_statuses AS s ON i.status = s.id
WHERE user_id = $1;
