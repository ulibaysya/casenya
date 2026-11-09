UPDATE users
SET balance = balance + $1
WHERE user = $2;
