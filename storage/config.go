package storage

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds the connection settings for the storage client.
//
// It works with any S3-compatible backend (MinIO, AWS S3, Cloudflare R2,
// Backblaze B2, etc.). Only Endpoint, AccessKey and SecretKey are required.
type Config struct {
	// Endpoint is the host[:port] WITHOUT scheme, e.g. "localhost:9000" or
	// "s3.amazonaws.com". Use UseSSL to control http vs https.
	Endpoint string
	// AccessKey / SecretKey are the credentials.
	AccessKey string
	SecretKey string
	// Token is an optional session token (STS / temporary credentials).
	Token string
	// UseSSL switches between http (false) and https (true).
	UseSSL bool
	// Region is optional, e.g. "us-east-1". Leave empty for MinIO defaults.
	Region string
	// Bucket is the default bucket used when a method is not given an override.
	Bucket string
	// PresignExpiry is the default lifetime of presigned URLs. Defaults to 15m.
	PresignExpiry time.Duration
	// PublicBaseURL, if set, is used by PublicURL instead of the endpoint.
	// Handy when objects are served through a CDN. Example:
	// "https://cdn.example.com" -> PublicURL("a/b.png") = "https://cdn.example.com/a/b.png".
	PublicBaseURL string
}

func (c Config) withDefaults() Config {
	if c.PresignExpiry <= 0 {
		c.PresignExpiry = 15 * time.Minute
	}
	return c
}

func (c Config) validate() error {
	if c.Endpoint == "" {
		return fmt.Errorf("%w: Endpoint is required", ErrInvalidConfig)
	}
	if c.AccessKey == "" || c.SecretKey == "" {
		return fmt.Errorf("%w: AccessKey and SecretKey are required", ErrInvalidConfig)
	}
	return nil
}

// ConfigFromEnv builds a Config from environment variables using the "MINIO_"
// prefix. See ConfigFromEnvPrefix for the full list of variables.
func ConfigFromEnv() Config { return ConfigFromEnvPrefix("MINIO_") }

// ConfigFromEnvPrefix builds a Config from environment variables using the
// given prefix. Recognised variables (shown with the default "MINIO_" prefix):
//
//	MINIO_ENDPOINT          host[:port] without scheme  (required)
//	MINIO_ACCESS_KEY        access key                  (required)
//	MINIO_SECRET_KEY        secret key                  (required)
//	MINIO_SESSION_TOKEN     optional STS session token
//	MINIO_USE_SSL           "true"/"false"
//	MINIO_REGION            optional region
//	MINIO_BUCKET            default bucket
//	MINIO_PRESIGN_EXPIRY    Go duration, e.g. "15m", "1h"
//	MINIO_PUBLIC_BASE_URL   optional CDN/public base URL
func ConfigFromEnvPrefix(prefix string) Config {
	c := Config{
		Endpoint:      os.Getenv(prefix + "ENDPOINT"),
		AccessKey:     os.Getenv(prefix + "ACCESS_KEY"),
		SecretKey:     os.Getenv(prefix + "SECRET_KEY"),
		Token:         os.Getenv(prefix + "SESSION_TOKEN"),
		Region:        os.Getenv(prefix + "REGION"),
		Bucket:        os.Getenv(prefix + "BUCKET"),
		PublicBaseURL: os.Getenv(prefix + "PUBLIC_BASE_URL"),
	}
	if v := os.Getenv(prefix + "USE_SSL"); v != "" {
		c.UseSSL, _ = strconv.ParseBool(v)
	}
	if v := os.Getenv(prefix + "PRESIGN_EXPIRY"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			c.PresignExpiry = d
		}
	}
	return c
}
