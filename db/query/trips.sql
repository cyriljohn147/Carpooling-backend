-- name: CreateTrip :one
INSERT INTO trips (driver_id, date, time,total_cost)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTrip :one
SELECT * FROM trips
WHERE id=$1
LIMIT 1;

-- name: ListTrips :many
SELECT * FROM trips
ORDER BY date DESC, time DESC
LIMIT $1 OFFSET $2;

-- name: UpdateTrip :one
UPDATE trips
SET driver_id=$2, date=$3, time=$4,total_cost=$5
WHERE id=$1
RETURNING *;

-- name: DeleteTrip :exec
DELETE FROM trips
WHERE id=$1;

-- name: CreateTripWithParticipantsAndBalances :exec
-- This should be executed within a transaction in your Go code
-- Parameters: $1=driver_id, $2=date, $3=time, $4=total_cost, $5=participant_data (JSON or separate calls)
-- Note: This is a template - actual implementation will require multiple queries in a transaction
BEGIN;
-- Insert trip, then participants, then update balances
-- Use this as a guide for your transaction logic in Go
COMMIT;

-- name: GetTripWithParticipants :many
SELECT
    t.id as trip_id,
    t.driver_id,
    t.date,
    t.time,
    t.total_cost,
    p.user_id as participant_id,
    p.share_amount,
    u.name as participant_name
FROM trips t
LEFT JOIN participants p ON t.id = p.trip_id
LEFT JOIN users u ON p.user_id = u.id
WHERE t.id = $1;
