BEGIN;

-- -- Lock inventory item row
-- SELECT status, attached_price
-- FROM users_items
-- WHERE id = $1
-- FOR SHARE:

-- Change status
UPDATE users_items
SET status = (
    SELECT s.id
    FROM user_items_statuses AS s
    WHERE s.code = 'sold'
)
WHERE id = $1
    AND user_id = $2
    AND status = (
        SELECT s.id
        FROM user_items_statuses AS s
        WHERE s.code = 'active'
    )
RETURNING attached_price;

-- Change balance
UPDATE users
set balance = balance + $3
where id = $2;

COMMIT;
