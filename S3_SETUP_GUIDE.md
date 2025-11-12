# Amazon S3 Integration Setup Guide

This guide will help you set up Amazon S3 storage for video files in the Cowatching application.

## Overview

The application has been updated to store video files in Amazon S3 instead of the local filesystem. Video metadata (including S3 URLs and keys) is stored in the PostgreSQL database.

## Prerequisites

1. AWS Account with S3 access
2. AWS Access Key ID and Secret Access Key with S3 permissions
3. An S3 bucket created for storing videos

## Step 1: Create an S3 Bucket

1. Log in to the [AWS Management Console](https://console.aws.amazon.com/)
2. Navigate to S3 service
3. Click "Create bucket"
4. Configure your bucket:
   - **Bucket name**: Choose a unique name (e.g., `cowatching-videos-prod`)
   - **Region**: Select a region close to your users (e.g., `us-east-1`)
   - **Block Public Access**: Keep all public access blocked (recommended)
   - **Versioning**: Optional (recommended for production)
   - **Encryption**: Enable server-side encryption (recommended)
5. Click "Create bucket"

## Step 2: Create IAM User with S3 Permissions

1. Navigate to IAM service in AWS Console
2. Click "Users" → "Add users"
3. Create a new user (e.g., `cowatching-s3-user`)
4. Attach the following policy (or create a custom policy):

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:GetObject",
        "s3:DeleteObject",
        "s3:ListBucket"
      ],
      "Resource": [
        "arn:aws:s3:::YOUR-BUCKET-NAME/*",
        "arn:aws:s3:::YOUR-BUCKET-NAME"
      ]
    }
  ]
}
```

5. Create access credentials and save them securely

## Step 3: Configure Environment Variables

Update your `.env` file with the S3 configuration:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=cowatching

# Server Configuration
SERVER_PORT=8080

# AWS S3 Configuration
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-access-key-id-here
AWS_SECRET_ACCESS_KEY=your-secret-access-key-here
S3_BUCKET_NAME=your-bucket-name-here
S3_VIDEO_PREFIX=videos/
```

**Important**: Never commit your `.env` file to version control!

## Step 4: Run Database Migration

The application requires a database schema update to support S3 fields. Run the migration:

```bash
# If using Docker Compose
docker-compose exec backend sh -c "psql \$DATABASE_URL -f db/migrations/003_add_s3_fields.sql"

# Or manually using psql
psql -h localhost -U postgres -d cowatching -f backend/db/migrations/003_add_s3_fields.sql
```

The migration adds the following columns to the `uploaded_videos` table:
- `s3_key` (VARCHAR): The S3 object key
- `s3_url` (TEXT): The full S3 URL
- Makes `filename` nullable since we now use `s3_key` as the primary identifier

## Step 5: Rebuild and Restart the Application

```bash
# If using Docker Compose
docker-compose down
docker-compose up --build

# Or if running locally
cd backend
go build -o bin/api cmd/api/main.go
./bin/api
```

## How It Works

### Video Upload Flow

1. User uploads a video via the frontend
2. Backend receives the video file
3. Video is uploaded to S3 with a unique key: `{timestamp}_{filename}`
4. S3 returns the URL of the uploaded file
5. Metadata is saved to the database with the S3 key and URL

### Video Streaming Flow

1. User requests to stream a video by ID
2. Backend retrieves video metadata from the database
3. Generates a presigned S3 URL (valid for 1 hour)
4. Redirects the user to the presigned URL
5. Video streams directly from S3 to the user

### Video Deletion Flow

1. User requests to delete a video by ID
2. Backend deletes the database record
3. Backend deletes the file from S3
4. Returns success response

## Key Features

- **Presigned URLs**: Videos are served via time-limited presigned URLs for security
- **No Local Storage**: Videos are stored entirely in S3
- **Scalability**: S3 handles all file storage and delivery
- **Durability**: S3 provides 99.999999999% durability
- **Error Handling**: Automatic cleanup on upload failures

## Configuration Options

### S3_VIDEO_PREFIX

The `S3_VIDEO_PREFIX` variable allows you to organize videos in a subfolder within your bucket:
- Default: `videos/`
- Example keys: `videos/1699875234_myvideo.mp4`

### Presigned URL Expiration

Currently set to 1 hour. To change this, modify `backend/internal/handlers/video.go`:

```go
// Generate presigned URL with custom expiration
presignedURL, err := h.s3Client.GetVideoURL(r.Context(), *video.S3Key, 2*time.Hour)
```

## Security Considerations

1. **Private Buckets**: Keep your S3 bucket private (block all public access)
2. **Presigned URLs**: Videos are accessed via time-limited presigned URLs
3. **IAM Permissions**: Use minimal IAM permissions (principle of least privilege)
4. **Credentials**: Never commit AWS credentials to version control
5. **HTTPS**: All S3 communications use HTTPS by default

## Cost Considerations

AWS S3 charges for:
- **Storage**: ~$0.023 per GB/month (Standard storage)
- **Requests**: ~$0.0004 per 1,000 PUT requests, ~$0.0004 per 10,000 GET requests
- **Data Transfer**: Free for inbound, ~$0.09 per GB for outbound (first 10TB)

For a small application with moderate usage, costs are typically very low.

## Troubleshooting

### "Failed to upload video to S3"

- Check AWS credentials in `.env`
- Verify IAM user has `s3:PutObject` permission
- Ensure bucket name is correct
- Check AWS region matches

### "Video not stored in S3"

- Run the database migration
- Verify new videos have `s3_key` populated in the database

### "Failed to generate presigned URL"

- Check IAM user has `s3:GetObject` permission
- Verify S3 key exists in the bucket
- Check AWS credentials are valid

## Monitoring

Consider setting up:
- **CloudWatch**: Monitor S3 bucket metrics
- **S3 Access Logs**: Track all requests to your bucket
- **Cost Alerts**: Get notified when costs exceed thresholds

## Migration from Local Storage

If you have existing videos stored locally, you'll need to:

1. Upload existing files to S3
2. Update database records with S3 keys and URLs
3. Delete local files after verification

Example migration script would be needed (not included).

## Support

For issues or questions:
- Check AWS S3 documentation: https://docs.aws.amazon.com/s3/
- Review application logs for error messages
- Verify all environment variables are set correctly
