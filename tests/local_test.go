package tests

import (
	"context"
	"testing"

	"github.com/scheiblingco/librc/remote/local"
)

func TestLocal(t *testing.T) {
	fsys, err := local.New(context.Background(), local.Config{
		Name:     "local-test",
		RootPath: t.TempDir(),
	})
	if err != nil {
		t.Fatalf("local.New: %v", err)
	}
	runSuite(t, fsys)
}
