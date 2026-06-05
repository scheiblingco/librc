package tests

import (
	"context"
	"testing"

	"github.com/scheiblingco/librc/remote/s3"
)

// TestS3 requires the minio service from compose.yml running on localhost:9000.
// The createbucket service pre-creates the librc-test bucket.
// Credentials: minioadmin / minioadmin123
func TestS3(t *testing.T) {
	ctx := context.Background()

	fsys, err := s3.New(ctx, s3.Config{
		Name:            "s3-test",
		Bucket:          "librc-test",
		Provider:        s3.ProviderMinio,
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
