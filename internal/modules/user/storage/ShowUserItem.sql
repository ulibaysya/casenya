SELECT
    item_id,
    attached_price
FROM users_items
WHERE 
    user_id = $1
    AND id = $2;
