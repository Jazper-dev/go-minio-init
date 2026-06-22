// Package storage is a small, reusable object-storage client backed by
// MinIO / S3 (minio-go v7).
//
// Create a Client from a Config (or environment variables) and use the
// upload / download / delete / URL / image helpers. It works with any
// S3-compatible backend: MinIO, AWS S3, Cloudflare R2, Backblaze B2, etc.
//
// # Quick start
//
//	client, err := storage.New(storage.Config{
//		Endpoint:  "localhost:9000",
//		AccessKey: "minioadmin",
//		SecretKey: "minioadmin",
//		Bucket:    "uploads",
//	})
//	if err != nil {
//		return err
//	}
//	if err := client.EnsureBucket(ctx); err != nil {
//		return err
//	}
//	_, err = client.UploadBytes(ctx, "hello.txt", []byte("hi"), nil)
//
// # Testing
//
// Depend on the Storage interface rather than *Client so it can be mocked.
//
// # Errors
//
// Backend "not found" errors are normalised to ErrNotFound; test with
// errors.Is. Other sentinel errors: ErrInvalidConfig, ErrEmptyBucket and
// ErrUnsupportedContentType.
package storage
