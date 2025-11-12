-- name: CreateUploadedVideo :one
INSERT INTO uploaded_videos (title, filename, content_type, file_size, s3_key, s3_url)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUploadedVideoByID :one
SELECT * FROM uploaded_videos
WHERE id = $1;

-- name: GetUploadedVideoByFilename :one
SELECT * FROM uploaded_videos
WHERE filename = $1;

-- name: ListUploadedVideos :many
SELECT * FROM uploaded_videos
ORDER BY created_at DESC;

-- name: UpdateUploadedVideo :one
UPDATE uploaded_videos
SET title = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteUploadedVideo :exec
DELETE FROM uploaded_videos
WHERE id = $1;

-- name: CountUploadedVideos :one
SELECT COUNT(*) FROM uploaded_videos;
