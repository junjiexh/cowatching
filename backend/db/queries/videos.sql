-- name: CreateVideo :one
INSERT INTO videos (room_id, title, url)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetVideoByID :one
SELECT * FROM videos
WHERE id = $1;

-- name: GetVideosByRoom :many
SELECT * FROM videos
WHERE room_id = $1
ORDER BY created_at DESC;

-- name: UpdateVideoPlayback :one
UPDATE videos
SET playback_position = $2, is_playing = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteVideo :exec
DELETE FROM videos
WHERE id = $1;
