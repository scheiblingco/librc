package tests

import (
	"context"
	"testing"

	"github.com/scheiblingco/librc/remote/ftp"
)

// TestFTP requires the ftp service from compose.yml running on localhost:21.
// Credentials: testuser / testpass
func TestFTP(t *testing.T) {
	fsys, err := ftp.New(context.Background(), ftp.Config{
		Name:     "ftp-test",
		Host:     "localhost",
		Port:     21,
		User:     "testuser",
		Password: "testpass",
	})
	if err != nil {
		t.Fatalf("ftp.New: %v", err)
	}
	runSuite(t, fsys)
}
