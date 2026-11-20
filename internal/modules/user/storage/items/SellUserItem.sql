BEGIN;

-- Lock user, check existence
SELECT EXISTS (
    SELECT id
    FROM users
    WHERE id = $1
    FOR KEY SHARE
);

-- Lock inventory item, check existence, check status
SELECT 
    i.attached_price,
    s.code
FROM users_items AS i
JOIN users_items_statuses AS s ON s.id = i.status_id
WHERE i.id = $1
    AND i.user_id = $2
FOR UPDATE OF i;

-- Change status
UPDATE users_items
SET status = (
    SELECT s.id
    FROM users_items_statuses AS s
    WHERE s.code = 'sold'
)
WHERE id = $1
    AND user_id = $2
    AND status = (
        SELECT s.id
        FROM users_items_statuses AS s
        WHERE s.code = 'active'
    )
RETURNING user_id, attached_price;

-- Change balance
UPDATE users
set balance = balance + $3
where id = $2
RETURNING balance;

COMMIT;
