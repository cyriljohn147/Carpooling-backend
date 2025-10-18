-- name: GetAllUsersCurrentBalances :many
-- Gets current balances for all users who owe money
SELECT
    id as user_id,
    name,
    COALESCE(balance, 0) as total_owed
FROM users
WHERE balance > 0
ORDER BY name;

-- name: GetCurrentUserBalance :one
-- Gets the current balance for a user
SELECT
    id as user_id,
    COALESCE(balance, 0) as total_owed
FROM users
WHERE id = $1;

-- name: PayOffUserBalance :exec
-- Marks user's balance as paid (sets to 0)
UPDATE users
SET balance = 0.00
WHERE id = $1;

-- name: PayPartialUserBalance :exec
-- Reduces user's balance by payment amount
UPDATE users 
SET balance = CASE 
    WHEN balance <= $2 THEN 0.00
    ELSE balance - $2
END
WHERE id = $1;

-- name: UpdateUserBalance :exec
-- Updates user's balance by adding an amount (for trip cost sharing)
UPDATE users 
SET balance = COALESCE(balance, 0) + $2
WHERE id = $1;
