# go-minio-init

A small, reusable object-storage client for Go, backed by MinIO / S3
(`minio-go/v7`). Drop the `storage` package into any project: build a client
from a `Config` (or env vars) and use the upload / download / delete / URL /
image helpers. Works with MinIO, AWS S3, Cloudflare R2, Backblaze B2, and any
other S3-compatible backend.

## Install

```bash
go get github.com/Jazper-dev/go-minio-init/storage
```

> The module path in `go.mod` is `github.com/Jazper-dev/go-minio-init`. Rename it
> to your own repo path before pushing, then run `go mod tidy`.

## Quick start

```go
package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/Jazper-dev/go-minio-init/storage"
)

func main() {
	ctx := context.Background()

	client, err := storage.New(storage.Config{
		Endpoint:  "localhost:9000",
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
		UseSSL:    false,
		Bucket:    "uploads",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create the bucket on startup if missing.
	if err := client.EnsureBucket(ctx); err != nil {
		log.Fatal(err)
	}

	// Upload some bytes.
	res, err := client.UploadBytes(ctx, "hello.txt",
		[]byte("hi"), &storage.UploadOptions{ContentType: "text/plain"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("uploaded %s (%d bytes)", res.Key, res.Size)

	// Get a temporary download URL (valid 10 minutes).
	url, _ := client.PresignGet(ctx, "hello.txt", 10*time.Minute)
	log.Println("download:", url)

	// Upload an image, capped at 1024px wide, validated as a real image.
	_, err = client.UploadImage(ctx, "avatars/u1.jpg",
		strings.NewReader("..."), &storage.ImageOptions{MaxWidth: 1024})
	if err != nil {
		log.Fatal(err)
	}
}
```

## Configure from environment

```go
client, err := storage.NewFromEnv() // reads MINIO_* variables
```

| Variable | Meaning |
| --- | --- |
| `MINIO_ENDPOINT` | host[:port], no scheme (required) |
| `MINIO_ACCESS_KEY` | access key (required) |
| `MINIO_SECRET_KEY` | secret key (required) |
| `MINIO_SESSION_TOKEN` | optional STS token |
| `MINIO_USE_SSL` | `true` / `false` |
| `MINIO_REGION` | optional region |
| `MINIO_BUCKET` | default bucket |
| `MINIO_PRESIGN_EXPIRY` | Go duration, e.g. `15m`, `1h` |
| `MINIO_PUBLIC_BASE_URL` | optional CDN / public base URL |

## Project layout

```
go-minio-init/
├── go.mod
├── LICENSE
├── README.md
├── examples/
│   └── basic/
│       └── main.go        # runnable demo / smoke-test
└── storage/               # the importable package
    ├── doc.go             # package overview (godoc)
    ├── config.go          # Config, defaults, validation, env loading
    ├── minio.go           # Client, New, NewFromEnv, bucket resolution
    ├── types.go           # Storage interface, ObjectInfo, UploadResult, options
    ├── errors.go          # sentinel errors + S3 error normalisation
    ├── upload.go          # Upload, UploadBytes, UploadFile
    ├── download.go        # Get, GetBytes, Download, Stat, Exists, ListObjects
    ├── delete.go          # Delete, DeleteMany
    ├── url.go             # PresignGet, PresignPut, PublicURL
    ├── bucket.go          # EnsureBucket, MakeBucket, BucketExists, RemoveBucket, SetPublicReadPolicy
    ├── health.go          # Ping, Health
    ├── image.go           # UploadImage (content-type validation + optional resize)
    └── *_test.go          # offline unit tests (no server required)
```

Run the demo against a local MinIO:

```bash
go run ./examples/basic
```

Run the tests (no server needed):

```bash
go test ./...
```

## Testing your code

Depend on the `storage.Storage` interface instead of `*storage.Client`, and
swap in your own mock in tests:

```go
type Service struct{ store storage.Storage }
```

## Error handling

```go
_, _, err := client.Get(ctx, "missing.txt")
if errors.Is(err, storage.ErrNotFound) {
	// 404
}
```

Sentinel errors: `ErrInvalidConfig`, `ErrEmptyBucket`, `ErrNotFound`,
`ErrUnsupportedContentType`.

## Notes

- `Upload` accepts `size = -1` when the length is unknown (uses multipart).
- `UploadImage` is forgiving: if resizing fails it uploads the original bytes.
  Resizing applies to JPEG/PNG; other allowed formats are stored unchanged.
- `PublicURL` only resolves to a working link for public buckets
  (`SetPublicReadPolicy`) or when `PublicBaseURL` points at a CDN.
- Need something not exposed here? `client.Raw()` returns the underlying
  `*minio.Client`.
