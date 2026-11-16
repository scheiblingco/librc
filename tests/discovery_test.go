//go:build integration

package tests

import (
	"context"
	"slices"
	"testing"

	"github.com/scheiblingco/librc/remote/azure"
	"github.com/scheiblingco/librc/remote/s3"
	"github.com/scheiblingco/librc/remote/smb"
)

// TestDiscovery lists what the compose services' credentials can see: the
// S3 bucket and SMB share they're started with, and an Azurite container
// this test creates. The configs name a bucket, container and share on
// purpose: listing ignores them.
func TestDiscovery(t *testing.T) {
	ctx := context.Background()

	buckets, err := s3.ListBuckets(ctx, s3.Config{
		Name:            "s3-discovery",
		Bucket:          "ignored",
		Provider:        s3.ProviderOther,
		AccessKeyID:     "minioadmin",
		SecretAccessKey: "minioadmin123",
		Endpoint:        "http://localhost:9000",
		ForcePathStyle:  true,
	})
	if err != nil || !slices.Contains(buckets, "librc-test") {
		t.Errorf("s3.ListBuckets: %v %v", buckets, err)
	}

	azcfg := azure.Config{Name: "azure-discovery", Container: "librc-discovery", ConnectionString: azuriteConnStr}
	fsys, err := azure.New(ctx, azcfg)
	if err != nil {
		t.Fatalf("azure.New: %v", err)
	}
	if err := fsys.Mkdir(ctx, ""); err != nil {
		t.Fatalf("create container: %v", err)
	}
	azcfg.Container = "ignored"
	containers, err := azure.ListContainers(ctx, azcfg)
	if err != nil || !slices.Contains(containers, "librc-discovery") || !slices.IsSorted(containers) {
		t.Errorf("azure.ListContainers: %v %v", containers, err)
	}

	shares, err := smb.ListShares(ctx, smb.Config{
		Name:     "smb-discovery",
		Host:     "localhost",
		User:     "testuser",
		Password: "testpass",
		Share:    "ignored",
		Root:     "ignored",
	})
	t.Logf("buckets %v, containers %v, shares %v", buckets, containers, shares)
	if err != nil || !slices.Contains(shares, "share") {
		t.Errorf("smb.ListShares: %v %v", shares, err)
	}
}
