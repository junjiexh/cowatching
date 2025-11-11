-- name: CreateRoom :one
INSERT INTO rooms (name, code, owner_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRoomByID :one
SELECT * FROM rooms
WHERE id = $1;

-- name: GetRoomByCode :one
SELECT * FROM rooms
WHERE code = $1;

-- name: ListRoomsByOwner :many
SELECT * FROM rooms
WHERE owner_id = $1 AND is_active = true
ORDER BY created_at DESC;

-- name: UpdateRoom :one
UPDATE rooms
SET name = $2, is_active = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteRoom :exec
DELETE FROM rooms
WHERE id = $1;

-- name: AddRoomParticipant :one
INSERT INTO room_participants (room_id, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: RemoveRoomParticipant :exec
DELETE FROM room_participants
WHERE room_id = $1 AND user_id = $2;

-- name: GetRoomParticipants :many
SELECT u.* FROM users u
JOIN room_participants rp ON u.id = rp.user_id
WHERE rp.room_id = $1;
