-- Add S3 storage fields to uploaded_videos table
ALTER TABLE uploaded_videos
ADD COLUMN s3_key VARCHAR(500),
ADD COLUMN s3_url TEXT;

-- Create index on s3_key for faster lookups
CREATE INDEX idx_uploaded_videos_s3_key ON uploaded_videos(s3_key);

-- Make filename nullable since we'll use s3_key as primary identifier
ALTER TABLE uploaded_videos
ALTER COLUMN filename DROP NOT NULL;
