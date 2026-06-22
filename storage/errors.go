package storage

import (
	"errors"
	"fmt"

	"github.com/minio/minio-go/v7"
)

// Sentinel errors. Use errors.Is to test for them, e.g.:
//
//	if errors.Is(err, storage.ErrNotFound) { ... }
var (
	// ErrInvalidConfig is returned by New when the config is incomplete.
	ErrInvalidConfig = errors.New("storage: invalid config")
	// ErrEmptyBucket is returned when no bucket is configured nor passed.
	ErrEmptyBucket = errors.New("storage: bucket name is empty")
	// ErrNotFound is returned when an object (or bucket) does not exist.
	ErrNotFound = errors.New("storage: object not found")
	// ErrUnsupportedContentType is returned by image helpers when the
	// detected content type is not in the allowed list.
	ErrUnsupportedContentType = errors.New("storage: unsupported content type")
)

// wrapErr normalises low-level S3 errors into the package's sentinel errors
// where possible, while preserving the original error for context.
func wrapErr(err error) error {
	if err == nil {
		return nil
	}
	resp := minio.ToErrorResponse(err)
	switch resp.Code {
	case "NoSuchKey", "NoSuchBucket", "NotFound":
		return fmt.Errorf("%w: %v", ErrNotFound, err)
	}
	return err
}
