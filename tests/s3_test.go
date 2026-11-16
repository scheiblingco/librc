//go:build integration

package tests

import (
	"context"
	"testing"

	"github.com/scheiblingco/librc/remote/s3"
)

// TestS3 requires the s3 service from compose.yml (rclone serve s3) on localhost:9000.
// The service creates the librc-test bucket at startup.
// Credentials: minioadmin / minioadmin123
func TestS3(t *testing.T) {
	ctx := context.Background()

	fsys, err := s3.New(ctx, s3.Config{
		Name:            "s3-test",
		Bucket:          "librc-test",
		Provider:        s3.ProviderOther,
		AccessKeyID:     "minioadmin",
		SecretAccessKey: "minioadmin123",
		Endpoint:        "http://localhost:9000",
		ForcePathStyle:  true,
	})
	if err != nil {
		t.Fatalf("s3.New: %v", err)
	}

	runSuite(t, fsys)
}
