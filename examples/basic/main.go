// Command basic is a small demo / smoke-test for the storage package.
//
// Usage:
//
//	export MINIO_ENDPOINT=localhost:9000
//	export MINIO_ACCESS_KEY=minioadmin
//	export MINIO_SECRET_KEY=minioadmin
//	export MINIO_BUCKET=uploads
//	go run ./examples/basic
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lfssteel/go-minio-init/storage"
)

func main() {
	ctx := context.Background()

	client, err := storage.NewFromEnv()
	if err != nil {
		log.Fatalf("init: %v", err)
	}

	if h := client.Health(ctx); !h.OK {
		log.Fatalf("backend unreachable: %s", h.Err)
	} else {
		fmt.Printf("connected to %s in %s\n", h.Endpoint, h.Latency)
	}

	if err := client.EnsureBucket(ctx); err != nil {
		log.Fatalf("ensure bucket: %v", err)
	}

	res, err := client.UploadBytes(ctx, "demo/hello.txt",
		[]byte("hello from go-minio-init"),
		&storage.UploadOptions{ContentType: "text/plain"})
	if err != nil {
		log.Fatalf("upload: %v", err)
	}
	fmt.Printf("uploaded %s/%s (%d bytes)\n", res.Bucket, res.Key, res.Size)

	url, err := client.PresignGet(ctx, "demo/hello.txt", 5*time.Minute)
	if err != nil {
		log.Fatalf("presign: %v", err)
	}
	fmt.Println("download URL:", url)
}
