//go:build integration

package tests

import (
	"context"
	"testing"

	"github.com/scheiblingco/librc/remote/webdav"
)

// TestWebDAV requires the webdav service from compose.yml (rclone serve
// webdav) on localhost:18080 with credentials testuser / testpass.
func TestWebDAV(t *testing.T) {
	fsys, err := webdav.New(context.Background(), webdav.Config{
		Name:     "webdav-test",
		URL:      "http://localhost:18080",
		Vendor:   webdav.VendorRclone,
		User:     "testuser",
		Password: "testpass",
	})
	if err != nil {
		t.Fatalf("webdav.New: %v", err)
	}
	runSuite(t, fsys)
}
