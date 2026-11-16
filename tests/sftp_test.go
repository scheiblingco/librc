//go:build integration

package tests

import (
	"context"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/scheiblingco/librc/remote/sftp"
)

// TestSFTP requires the sftp service from compose.yml running on localhost:2222.
// Credentials: testuser / testpass, writable root: upload/
//
// The test uses ssh-keyscan to retrieve the container's host key automatically,
// so no manual known_hosts setup is required.
func TestSFTP(t *testing.T) {
	fsys, err := sftp.New(context.Background(), sftp.Config{
		Name:           "sftp-test",
		Host:           "localhost",
		Port:           2222,
		User:           "testuser",
		Password:       "testpass",
		KnownHostsFile: scanSFTPHostKey(t, "localhost", 2222),
		Root:           "upload",
	})
	if err != nil {
		t.Fatalf("sftp.New: %v", err)
	}
	runSuite(t, fsys)
}

// scanSFTPHostKey runs ssh-keyscan against host:port, writes the result to a
// temp file, and returns its path. The test is skipped if ssh-keyscan is not
// available or the server is not reachable.
func scanSFTPHostKey(t *testing.T, host string, port int) string {
	t.Helper()
	out, err := exec.Command("ssh-keyscan", "-p", strconv.Itoa(port), host).Output()
	if err != nil || len(out) == 0 {
		t.Skipf("ssh-keyscan -p %d %s failed (is the sftp container running?): %v", port, host, err)
	}
	f, err := os.CreateTemp(t.TempDir(), "known_hosts_*")
	if err != nil {
		t.Fatalf("create known_hosts temp file: %v", err)
	}
	if _, err := f.Write(out); err != nil {
		t.Fatalf("write known_hosts: %v", err)
	}
	f.Close()
	return f.Name()
}
