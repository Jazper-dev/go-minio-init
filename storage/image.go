package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/disintegration/imaging"
)

// DefaultAllowedImageTypes is the content-type allow-list used when
// ImageOptions.AllowedTypes is empty.
var DefaultAllowedImageTypes = []string{
	"image/jpeg",
	"image/png",
	"image/gif",
	"image/webp",
}

// ImageOptions controls UploadImage. The zero value uploads the image as-is
// (after content-type validation).
type ImageOptions struct {
	// Bucket overrides the client's default bucket.
	Bucket string
	// MaxWidth / MaxHeight, when > 0, shrink the image to fit within the box
	// while preserving aspect ratio. Images smaller than the box are left
	// untouched. Resizing applies to JPEG and PNG only; other formats are
	// uploaded unchanged.
	MaxWidth  int
	MaxHeight int
	// Quality is the JPEG quality (1-100) used when re-encoding. Default 85.
	Quality int
	// AllowedTypes overrides DefaultAllowedImageTypes.
	AllowedTypes []string
	// Metadata is stored as user metadata.
	Metadata map[string]string
}

// UploadImage validates that r is an allowed image type, optionally resizes it,
// and uploads it under key. The detected content type is stored on the object.
//
// It is intentionally forgiving: if resizing is requested but the image cannot
// be decoded/re-encoded, the original bytes are uploaded unchanged.
func (c *Client) UploadImage(ctx context.Context, key string, r io.Reader, opts *ImageOptions) (*UploadResult, error) {
	if opts == nil {
		opts = &ImageOptions{}
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("storage: read image: %w", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("%w: empty body", ErrUnsupportedContentType)
	}

	ct := http.DetectContentType(data)
	allowed := opts.AllowedTypes
	if len(allowed) == 0 {
		allowed = DefaultAllowedImageTypes
	}
	if !containsString(allowed, ct) {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedContentType, ct)
	}

	if (opts.MaxWidth > 0 || opts.MaxHeight > 0) && (ct == "image/jpeg" || ct == "image/png") {
		if resized, ok := resizeImage(data, ct, opts); ok {
			data = resized
		}
	}

	return c.UploadBytes(ctx, key, data, &UploadOptions{
		ContentType: ct,
		Metadata:    opts.Metadata,
		Bucket:      opts.Bucket,
	})
}

// resizeImage decodes, fits and re-encodes the image. It returns ok=false if
// any step fails, so callers can fall back to the original bytes.
func resizeImage(data []byte, ct string, opts *ImageOptions) ([]byte, bool) {
	img, err := imaging.Decode(bytes.NewReader(data), imaging.AutoOrientation(true))
	if err != nil {
		return nil, false
	}

	w := opts.MaxWidth
	if w <= 0 {
		w = 1 << 15
	}
	h := opts.MaxHeight
	if h <= 0 {
		h = 1 << 15
	}

	b := img.Bounds()
	// Only shrink; never upscale.
	if b.Dx() <= w && b.Dy() <= h {
		return nil, false
	}
	fitted := imaging.Fit(img, w, h, imaging.Lanczos)

	var buf bytes.Buffer
	switch ct {
	case "image/png":
		if err := imaging.Encode(&buf, fitted, imaging.PNG); err != nil {
			return nil, false
		}
	default: // image/jpeg
		q := opts.Quality
		if q <= 0 {
			q = 85
		}
		if err := imaging.Encode(&buf, fitted, imaging.JPEG, imaging.JPEGQuality(q)); err != nil {
			return nil, false
		}
	}
	return buf.Bytes(), true
}

func containsString(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}
