package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"cal-salary/core/config"
	"cal-salary/core/logger"

	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	minioClient *minio.Client
)

// khởi tạo MinIO client
func InitMinIOClient() error {
	cfg := config.Get()

	if !cfg.MinIO.Enabled {
		logger.Info("MinIO is disabled, using local storage")
		return nil
	}

	if cfg.MinIO.Endpoint == "" {
		return fmt.Errorf("MinIO endpoint is required")
	}

	if cfg.MinIO.AccessKeyID == "" {
		return fmt.Errorf("MinIO access key ID is required")
	}

	if cfg.MinIO.SecretAccessKey == "" {
		return fmt.Errorf("MinIO secret access key is required")
	}

	if cfg.MinIO.BucketName == "" {
		return fmt.Errorf("MinIO bucket name is required")
	}

	// Đảm bảo endpoint không chứa protocol
	endpoint := strings.TrimPrefix(strings.TrimPrefix(cfg.MinIO.Endpoint, "http://"), "https://")

	// Khởi tạo MinIO client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.AccessKeyID, cfg.MinIO.SecretAccessKey, ""),
		Secure: cfg.MinIO.UseSSL,
		Region: cfg.MinIO.Region,
	})

	if err != nil {
		return fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, cfg.MinIO.BucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		// Tạo bucket nếu chưa tồn tại
		err = client.MakeBucket(ctx, cfg.MinIO.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		logger.Info(fmt.Sprintf("Created bucket: %s", cfg.MinIO.BucketName))
	}

	minioClient = client
	logger.Info("MinIO client initialized successfully")
	return nil
}

// GetMinIOClient trả về MinIO client
func GetMinIOClient() *minio.Client {
	return minioClient
}

// IsMinIOEnabled kiểm tra MinIO có được bật không
func IsMinIOEnabled() bool {
	cfg := config.Get()
	return cfg.MinIO.Enabled && minioClient != nil
}

// UploadToMinIO upload file lên MinIO
func UploadToMinIO(ctx context.Context, fileData io.Reader, objectName string, size int64, contentType string) error {
	cfg := config.Get()

	if !IsMinIOEnabled() {
		return fmt.Errorf("MinIO is not enabled")
	}

	uploadInfo, err := minioClient.PutObject(ctx, cfg.MinIO.BucketName, objectName, fileData, size, minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		return fmt.Errorf("failed to upload to MinIO: %w", err)
	}

	logger.Info(fmt.Sprintf("Uploaded to MinIO: %s (size: %d bytes)", uploadInfo.Key, uploadInfo.Size))
	return nil
}

// DownloadFromMinIO download file từ MinIO
func DownloadFromMinIO(ctx context.Context, objectName string) (io.ReadCloser, error) {
	cfg := config.Get()

	if !IsMinIOEnabled() {
		return nil, fmt.Errorf("MinIO is not enabled")
	}

	object, err := minioClient.GetObject(ctx, cfg.MinIO.BucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download from MinIO: %w", err)
	}

	return object, nil
}

// DeleteFromMinIO xóa file từ MinIO
func DeleteFromMinIO(ctx context.Context, objectName string) error {
	cfg := config.Get()

	if !IsMinIOEnabled() {
		return fmt.Errorf("MinIO is not enabled")
	}

	err := minioClient.RemoveObject(ctx, cfg.MinIO.BucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete from MinIO: %w", err)
	}

	logger.Info(fmt.Sprintf("Deleted from MinIO: %s", objectName))
	return nil
}

// GetMinIOURL trả về URL của object trong MinIO
func GetMinIOURL(objectName string) string {
	cfg := config.Get()

	if !IsMinIOEnabled() {
		return ""
	}

	protocol := "http"
	if cfg.MinIO.UseSSL {
		protocol = "https"
	}

	endpoint := strings.TrimPrefix(strings.TrimPrefix(cfg.MinIO.Endpoint, "http://"), "https://")
	if cfg.MinIO.PublicEndpoint != "" {
		endpoint = strings.TrimPrefix(strings.TrimPrefix(cfg.MinIO.PublicEndpoint, "http://"), "https://")
	}

	return fmt.Sprintf("%s://%s/%s/%s", protocol, endpoint, cfg.MinIO.BucketName, objectName)
}
