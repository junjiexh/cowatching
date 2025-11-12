package s3

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Client struct {
	client       *s3.Client
	bucketName   string
	videoPrefix  string
	region       string
}

type S3Config struct {
	Region         string
	AccessKeyID    string
	SecretAccessKey string
	BucketName     string
	VideoPrefix    string
}

// NewS3Client creates a new S3 client with the provided configuration
func NewS3Client(cfg S3Config) (*S3Client, error) {
	if cfg.BucketName == "" {
		return nil, fmt.Errorf("S3 bucket name is required")
	}

	// Create AWS config with credentials
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(awsCfg)

	return &S3Client{
		client:      client,
		bucketName:  cfg.BucketName,
		videoPrefix: cfg.VideoPrefix,
		region:      cfg.Region,
	}, nil
}

// UploadVideo uploads a video file to S3
func (s *S3Client) UploadVideo(ctx context.Context, key string, body io.Reader, contentType string, fileSize int64) (string, error) {
	// Add prefix to the key
	fullKey := s.videoPrefix + key

	// Upload the file to S3
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(fullKey),
		Body:          body,
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(fileSize),
		ACL:           types.ObjectCannedACLPrivate,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Generate the S3 URL
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.region, fullKey)
	return url, nil
}

// GetVideoURL generates a presigned URL for downloading/streaming a video
func (s *S3Client) GetVideoURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	fullKey := s.videoPrefix + key

	// Create a presigned request
	presignClient := s3.NewPresignClient(s.client)
	presignResult, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fullKey),
	}, s3.WithPresignExpires(expiration))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignResult.URL, nil
}

// GetVideo retrieves a video file from S3
func (s *S3Client) GetVideo(ctx context.Context, key string) (io.ReadCloser, int64, error) {
	fullKey := s.videoPrefix + key

	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fullKey),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get object from S3: %w", err)
	}

	contentLength := int64(0)
	if result.ContentLength != nil {
		contentLength = *result.ContentLength
	}

	return result.Body, contentLength, nil
}

// DeleteVideo deletes a video file from S3
func (s *S3Client) DeleteVideo(ctx context.Context, key string) error {
	fullKey := s.videoPrefix + key

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fullKey),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object from S3: %w", err)
	}

	return nil
}

// HeadVideo checks if a video exists and returns its metadata
func (s *S3Client) HeadVideo(ctx context.Context, key string) (*s3.HeadObjectOutput, error) {
	fullKey := s.videoPrefix + key

	result, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fullKey),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to head object from S3: %w", err)
	}

	return result, nil
}
