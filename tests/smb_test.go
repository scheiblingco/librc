package tests

import (
	"context"
	"testing"

	"github.com/scheiblingco/librc/remote/smb"
)

// TestSMB requires the smb service from compose.yml (dperson/samba on
// localhost:445) with share "share" and credentials testuser/testpass.
func TestSMB(t *testing.T) {
	fsys, err := smb.New(context.Background(), smb.Config{
		Name:     "smb-test",
		Host:     "localhost",
		User:     "testuser",
		Password: "testpass",
		Share:    "share",
	})
	if err != nil {
		t.Fatalf("smb.New: %v", err)
	}
	runSuite(t, fsys)
}
