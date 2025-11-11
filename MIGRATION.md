# Database Migration Guide

## Overview

The video system has been refactored to use PostgreSQL for metadata storage with clean numeric IDs in URLs.

**Old URL format:** `/api/v1/videos/stream/1762845772_video.mp4`
**New URL format:** `/api/v1/videos/stream/123`

## Benefits

✅ Clean, short URLs without encoding issues
✅ Decouples URLs from filesystem structure
✅ Enables metadata tracking (views, processing status, etc.)
✅ Supports future features (permissions, soft deletes, file renaming)
✅ Scalable architecture

## Running the Migration

### Step 1: Ensure PostgreSQL is Running

```bash
# If using Docker Compose
docker-compose up -d postgres

# Or start PostgreSQL manually
# Make sure it's accessible on localhost:5432
```

### Step 2: Run the Migration

The new migration file is: `backend/db/migrations/002_create_uploaded_videos.sql`

**Option A: Using psql**

```bash
# From the project root
psql -U postgres -d cowatching -f backend/db/migrations/002_create_uploaded_videos.sql
```

**Option B: Run all migrations in order**

```bash
# Run all migrations (including the initial schema)
psql -U postgres -d cowatching -f backend/db/migrations/001_init_schema.sql
psql -U postgres -d cowatching -f backend/db/migrations/002_create_uploaded_videos.sql
```

**Option C: Using Docker Compose (fresh database)**

If you're starting fresh:

```bash
# Stop and remove existing containers and volumes
docker-compose down -v

# Start fresh - migrations will run automatically
docker-compose up -d
```

### Step 3: Verify the Migration

```bash
psql -U postgres -d cowatching -c "\d uploaded_videos"
```

You should see:

```
                         Table "public.uploaded_videos"
   Column     |           Type           | Collation | Nullable | Default
--------------+--------------------------+-----------+----------+---------
 id           | bigint                   |           | not null | nextval(...)
 title        | character varying(255)   |           | not null |
 filename     | character varying(500)   |           | not null |
 content_type | character varying(100)   |           | not null |
 file_size    | bigint                   |           | not null |
 created_at   | timestamp with time zone |           |          | CURRENT_TIMESTAMP
 updated_at   | timestamp with time zone |           |          | CURRENT_TIMESTAMP
```

## Important Notes

### Existing Videos

⚠️ **Videos uploaded before this change will NOT appear in the video list**

- The old system stored videos only in the filesystem
- The new system requires database records
- You have two options:
  1. **Clean start**: Delete old videos from `backend/uploads/videos/`
  2. **Migration script**: Create a script to import old videos into the database

### Fresh Start (Recommended)

If you don't need to keep old videos:

```bash
# Delete old video files
rm -rf backend/uploads/videos/*

# Start the backend
cd backend
go run cmd/api/main.go

# Upload videos through the UI - they'll be stored in the database
```

### Migrating Existing Videos (Advanced)

If you need to keep existing videos, you'll need to create database records for them:

```sql
-- Example: Insert existing video into database
INSERT INTO uploaded_videos (title, filename, content_type, file_size)
VALUES (
  'My Video',
  '1762845772_video.mp4',
  'video/mp4',
  52428800  -- file size in bytes
);
```

## API Changes

### Upload Endpoint
- **URL**: `POST /api/v1/videos/upload` (unchanged)
- **Response**: Now includes numeric `id` field

```json
{
  "id": 123,
  "title": "My Video",
  "url": "/api/v1/videos/stream/123",
  "size": 52428800,
  "contentType": "video/mp4",
  "uploadedAt": "2025-01-15T10:30:00Z"
}
```

### List Endpoint
- **URL**: `GET /api/v1/videos` (unchanged)
- **Response**: Array of videos with numeric IDs

### Stream Endpoint
- **Old**: `GET /api/v1/videos/stream/{filename}`
- **New**: `GET /api/v1/videos/stream/{id}`

### Delete Endpoint
- **Old**: `DELETE /api/v1/videos/{filename}`
- **New**: `DELETE /api/v1/videos/{id}`

## Frontend Changes

The frontend has been automatically updated to use numeric IDs:

- Video interface now uses `number` type for `id`
- Removed `filename` field from the public interface
- Delete requests now use numeric ID

No frontend code changes needed if you're using the latest version!

## Troubleshooting

### "Video not found" errors after migration

This happens when:
1. Database migration didn't run
2. Old videos exist in filesystem but not in database

**Solution**: Run the migration or delete old video files

### "Failed to save video metadata"

Check:
1. PostgreSQL is running
2. Database connection settings in `.env`
3. Migration was successful

```bash
# Check database connection
psql -U postgres -d cowatching -c "SELECT 1"

# Check if table exists
psql -U postgres -d cowatching -c "SELECT COUNT(*) FROM uploaded_videos"
```

### Build errors

If you see build errors related to `internal/database/db`:

```bash
cd backend
go mod tidy
go build ./...
```

## Rollback (if needed)

To revert to the old system:

```bash
git checkout HEAD~1  # Go back one commit
```

However, videos uploaded with the new system will need to be re-uploaded after rollback.

## Next Steps

After successful migration:

1. Test video upload through the UI
2. Verify videos appear in the list with numeric IDs
3. Test video playback
4. Test video deletion

## Support

If you encounter issues:

1. Check PostgreSQL logs
2. Check backend logs for error messages
3. Verify all migration steps completed successfully
4. Ensure `.env` configuration is correct
