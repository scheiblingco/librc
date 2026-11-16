//go:build integration

package tests

import (
	"context"
	"testing"

	"github.com/scheiblingco/librc/remote/hetzner"
)

// TestHetzner requires the hetzner-mock service from compose.yml (atmoz/sftp
// on localhost:2323). In production, set Host to "" and Username to your
// storage box username to use {Username}.your-storagebox.de:23.
func TestHetzner(t *testing.T) {
	fsys, err := hetzner.New(context.Background(), hetzner.Config{
		Name:           "hetzner-test",
		Username:       "u123456",
		Host:           "localhost",
		Port:           2323,
		Password:       "testpass",
		KnownHostsFile: scanSFTPHostKey(t, "localhost", 2323),
		Root:           "upload",
	})
	if err != nil {
		t.Fatalf("hetzner.New: %v", err)
	}
	runSuite(t, fsys)
}
