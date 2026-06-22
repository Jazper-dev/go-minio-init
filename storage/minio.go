package storage

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Compile-time assertion that *Client satisfies the Storage interface.
var _ Storage = (*Client)(nil)

// Client is a reusable object-storage client.
type Client struct {
	mc  *minio.Client
	cfg Config
}

// New creates a Client from the given Config.
func New(cfg Config) (*Client, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	mc, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, cfg.Token),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("storage: init client: %w", err)
	}
	return &Client{mc: mc, cfg: cfg}, nil
}

// NewFromEnv creates a Client from environment variables (MINIO_* prefix).
func NewFromEnv() (*Client, error) { return New(ConfigFromEnv()) }

// Raw exposes the underlying *minio.Client for advanced operations not
// covered by this package.
func (c *Client) Raw() *minio.Client { return c.mc }

// Config returns a copy of the active configuration.
func (c *Client) Config() Config { return c.cfg }

// resolveBucket picks the override (if any) over the default, and errors if
// neither is set.
func (c *Client) resolveBucket(override []string) (string, error) {
	b := c.cfg.Bucket
	if len(override) > 0 && override[0] != "" {
		b = override[0]
	}
	if b == "" {
		return "", ErrEmptyBucket
	}
	return b, nil
}

// bucket resolves a single optional override string.
func (c *Client) bucket(override string) (string, error) {
	if override != "" {
		return override, nil
	}
	if c.cfg.Bucket == "" {
		return "", ErrEmptyBucket
	}
	return c.cfg.Bucket, nil
}
