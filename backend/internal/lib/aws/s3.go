package aws

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Client struct {
	server  *server.Server
	client  *s3.Client
	presign *s3.PresignClient
}

func NewS3Client(server *server.Server, cfg aws.Config) *S3Client {
	awsConfig := server.Config.AWS

	// Internal endpoint — used for every real network call the backend
	// itself makes (HeadBucket, CreateBucket, PutObject, GetObject).
	// Must be reachable from inside this container.
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if awsConfig.EndpointURL != "" {
			o.BaseEndpoint = aws.String(awsConfig.EndpointURL)
			o.UsePathStyle = true
		}
	})

	// Public endpoint — used only to build the host embedded in presigned
	// URLs. No connection is ever opened with this client, so it just
	// needs to match whatever host the actual uploader can reach.
	publicEndpoint := awsConfig.PublicEndpointURL
	if publicEndpoint == "" {
		publicEndpoint = awsConfig.EndpointURL
	}

	presignTarget := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if publicEndpoint != "" {
			o.BaseEndpoint = aws.String(publicEndpoint)
			o.UsePathStyle = true
		}
	})

	return &S3Client{
		server:  server,
		client:  client,
		presign: s3.NewPresignClient(presignTarget),
	}
}

func (s *S3Client) UploadFile(ctx context.Context, bucket string, fileName string, file io.Reader) (string, error) {
	fileKey := fmt.Sprintf("%s_%d", fileName, time.Now().Unix())

	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(fileKey),
		Body:        bytes.NewReader(buffer.Bytes()),
		ContentType: aws.String(http.DetectContentType(buffer.Bytes())),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return fileKey, nil
}

func (s *S3Client) CreatePresignedDownloadUrl(ctx context.Context, bucket string, objectKey string) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	expiration := time.Minute * 60

	presignedUrl, err := presignClient.PresignGetObject(ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(objectKey),
		},
		s3.WithPresignExpires(expiration))
	if err != nil {
		return "", err
	}

	return presignedUrl.URL, nil
}

func (s *S3Client) CreatePresignedUploadURL(
	ctx context.Context,
	bucket string,
	objectKey string,
) (string, error) {
	expiration := time.Hour

	presignedURL, err := s.presign.PresignPutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(objectKey),
		},
		s3.WithPresignExpires(expiration),
	)
	if err != nil {
		return "", err
	}

	return presignedURL.URL, nil
}

func (s *S3Client) DeleteObject(ctx context.Context, bucket string, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object %s: %w", key, err)
	}

	return nil
}

func (c *S3Client) EnsureBucket(ctx context.Context, bucket string) error {
	_, err := c.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err == nil {
		return nil // already exists
	}

	var notFound *types.NotFound
	if !errors.As(err, &notFound) {
		// auth error, network error, etc. — don't swallow it as "missing bucket"
		return err
	}

	_, err = c.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	return err
}
