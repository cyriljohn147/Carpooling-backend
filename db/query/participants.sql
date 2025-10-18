-- name: AddParticipant :one
INSERT INTO participants (trip_id, user_id, share_amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetParticipant :one
SELECT * FROM participants
WHERE id=$1
LIMIT 1;

-- name: ListParticipantsByTrip :many
SELECT * FROM participants
WHERE trip_id=$1
ORDER BY id
LIMIT $2 OFFSET $3;

-- name: UpdateParticipantShare :one
UPDATE participants
SET share_amount=$2
WHERE id=$1
RETURNING *;

-- name: DeleteParticipant :exec
DELETE FROM participants
WHERE id=$1;
