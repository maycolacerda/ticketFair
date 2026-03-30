// services/s3_service.go
package services

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/maycolacerda/ticketfair/configs"
)

const (
	maxImageSize = 5 << 20 // 5 MB
	imageFolder  = "events"
)

var allowedTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

// UploadEventImage validates, uploads and returns the public URL
func UploadEventImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	if configs.S3Client == nil {
		return "", ErrS3NotConfigured
	}

	// Size guard
	if header.Size > maxImageSize {
		return "", ErrImageTooLarge
	}

	// Read file bytes
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return "", ErrFailedToCreate
	}
	fileBytes := buf.Bytes()

	// Detect MIME type
	mimeType := http.DetectContentType(fileBytes)
	if !allowedTypes[mimeType] {
		return "", ErrInvalidImageType
	}

	// Validate it's a real image
	if _, _, err := image.Decode(bytes.NewReader(fileBytes)); err != nil {
		return "", ErrInvalidImageType
	}

	// Build key: events/<uuid>.<ext>
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		ext = mimeTypeToExt(mimeType)
	}
	key := fmt.Sprintf("%s/%s%s", imageFolder, uuid.New().String(), ext)

	// Upload to S3
	_, err := configs.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(configs.S3Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(mimeType),
	})
	if err != nil {
		slog.Error("S3 upload failed", "key", key, "error", err.Error())
		return "", ErrFailedToCreate
	}

	url := buildImageURL(key)
	slog.Info("Image uploaded", "key", key, "url", url)
	return url, nil
}

// DeleteEventImage removes an image from S3 by its full URL or key
func DeleteEventImage(imageURL string) error {
	if configs.S3Client == nil || imageURL == "" {
		return nil
	}

	key := extractKeyFromURL(imageURL)
	if key == "" {
		return nil
	}

	_, err := configs.S3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(configs.S3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		slog.Error("S3 delete failed", "key", key, "error", err.Error())
		return ErrFailedToUpdate
	}

	slog.Info("Image deleted", "key", key)
	return nil
}

// GeneratePresignedURL returns a temporary signed URL for private access (future use)
func GeneratePresignedURL(key string, expiry time.Duration) (string, error) {
	if configs.S3Client == nil {
		return "", ErrS3NotConfigured
	}

	presigner := s3.NewPresignClient(configs.S3Client)
	req, err := presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(configs.S3Bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", ErrFailedToFetch
	}

	return req.URL, nil
}

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

func buildImageURL(key string) string {

	endpoint := "http://localhost:4566"
	if ep := os.Getenv("AWS_ENDPOINT_URL"); ep != "" {
		endpoint = ep
	}
	return fmt.Sprintf("%s/%s/%s", endpoint, configs.S3Bucket, key)
}

func extractKeyFromURL(url string) string {
	// Strip everything before the bucket name
	marker := configs.S3Bucket + "/"
	idx := strings.Index(url, marker)
	if idx < 0 {
		return ""
	}
	return url[idx+len(marker):]
}

func mimeTypeToExt(mime string) string {
	switch mime {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg"
	}
}
