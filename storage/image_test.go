package storage

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestContainsString(t *testing.T) {
	list := []string{"image/jpeg", "image/png"}
	if !containsString(list, "image/png") {
		t.Error("expected png to be found")
	}
	if containsString(list, "image/gif") {
		t.Error("did not expect gif to be found")
	}
}

func pngBytes(t *testing.T, w, h int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, color.RGBA{R: uint8(x), G: uint8(y), B: 100, A: 255})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode: %v", err)
	}
	return buf.Bytes()
}

func TestResizeImageShrinks(t *testing.T) {
	src := pngBytes(t, 200, 100)
	out, ok := resizeImage(src, "image/png", &ImageOptions{MaxWidth: 50, MaxHeight: 50})
	if !ok {
		t.Fatal("expected resize to succeed")
	}
	img, err := png.Decode(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode result: %v", err)
	}
	b := img.Bounds()
	if b.Dx() > 50 || b.Dy() > 50 {
		t.Errorf("resized to %dx%d, want within 50x50", b.Dx(), b.Dy())
	}
}

func TestResizeImageNoUpscale(t *testing.T) {
	src := pngBytes(t, 10, 10)
	if _, ok := resizeImage(src, "image/png", &ImageOptions{MaxWidth: 100, MaxHeight: 100}); ok {
		t.Error("expected no resize for an image smaller than the box")
	}
}
