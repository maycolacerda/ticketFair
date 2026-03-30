// configs/s3.go
package configs

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client
var S3Bucket string

func InitS3() {
	endpoint := os.Getenv("AWS_ENDPOINT_URL")
	region := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	S3Bucket = os.Getenv("S3_BUCKET")

	if S3Bucket == "" {
		S3Bucket = "ticketfair-images"
	}
	if region == "" {
		region = "us-east-1"
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)
	if err != nil {
		slog.Error("Failed to load S3 config", "error", err.Error())
		return
	}

	// Override endpoint for LocalStack
	S3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true // required for LocalStack
		}
	})

	// Ensure bucket exists
	ensureBucket()

	slog.Info("S3 initialized", "bucket", S3Bucket, "endpoint", endpoint)
}

func ensureBucket() {
	_, err := S3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(S3Bucket),
	})
	if err != nil {
		// Bucket may already exist — not a fatal error
		slog.Warn("S3 bucket create skipped (may already exist)",
			"bucket", S3Bucket,
			"error", err.Error(),
		)
	}
}
