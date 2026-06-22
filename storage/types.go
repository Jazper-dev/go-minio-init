package storage

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
)

// ObjectInfo describes a stored object.
type ObjectInfo struct {
	Key          string
	Bucket       string
	Size         int64
	ETag         string
	ContentType  string
	LastModified time.Time
	Metadata     map[string]string
}

// UploadResult is returned after a successful upload.
type UploadResult struct {
	Bucket      string
	Key         string
	ETag        string
	Size        int64
	VersionID   string
	ContentType string
}

// UploadOptions are optional settings for an upload. A nil *UploadOptions is
// valid and means "use defaults".
type UploadOptions struct {
	// ContentType, e.g. "image/png". Defaults to "application/octet-stream".
	ContentType string
	// Metadata is stored as user metadata (x-amz-meta-*).
	Metadata map[string]string
	// Bucket overrides the client's default bucket for this call.
	Bucket string
}

// Health reports the result of a connectivity check.
type Health struct {
	OK       bool
	Endpoint string
	Latency  time.Duration
	Err      string
}

// Storage is the behavioural interface implemented by *Client. Depend on this
// in your application code so it can be mocked in tests.
type Storage interface {
	Upload(ctx context.Context, key string, r io.Reader, size int64, opts *UploadOptions) (*UploadResult, error)
	UploadBytes(ctx context.Context, key string, data []byte, opts *UploadOptions) (*UploadResult, error)
	UploadFile(ctx context.Context, key, filePath string, opts *UploadOptions) (*UploadResult, error)
	UploadImage(ctx context.Context, key string, r io.Reader, opts *ImageOptions) (*UploadResult, error)

	Get(ctx context.Context, key string, bucket ...string) (io.ReadCloser, *ObjectInfo, error)
	GetBytes(ctx context.Context, key string, bucket ...string) ([]byte, *ObjectInfo, error)
	Download(ctx context.Context, key, filePath string, bucket ...string) error
	Stat(ctx context.Context, key string, bucket ...string) (*ObjectInfo, error)
	Exists(ctx context.Context, key string, bucket ...string) (bool, error)
	ListObjects(ctx context.Context, prefix string, recursive bool, bucket ...string) ([]ObjectInfo, error)

	Delete(ctx context.Context, key string, bucket ...string) error
	DeleteMany(ctx context.Context, keys []string, bucket ...string) error

	PresignGet(ctx context.Context, key string, expiry time.Duration, bucket ...string) (string, error)
	PresignPut(ctx context.Context, key string, expiry time.Duration, bucket ...string) (string, error)
	PublicURL(key string, bucket ...string) (string, error)

	EnsureBucket(ctx context.Context, bucket ...string) error
	MakeBucket(ctx context.Context, name string) error
	BucketExists(ctx context.Context, name string) (bool, error)
	RemoveBucket(ctx context.Context, name string) error
	SetPublicReadPolicy(ctx context.Context, bucket ...string) error

	Ping(ctx context.Context) error
	Health(ctx context.Context) Health
}

// objInfo converts a minio.ObjectInfo into the package's ObjectInfo.
func objInfo(bucket string, o minio.ObjectInfo) *ObjectInfo {
	return &ObjectInfo{
		Key:          o.Key,
		Bucket:       bucket,
		Size:         o.Size,
		ETag:         o.ETag,
		ContentType:  o.ContentType,
		LastModified: o.LastModified,
		Metadata:     o.UserMetadata,
	}
}
