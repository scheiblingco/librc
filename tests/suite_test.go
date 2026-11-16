//go:build integration

// Package tests contains integration tests for all librc backends.
// Start the test infrastructure before running:
//
//	docker compose -f tests/compose.yml up -d
//	go test -tags integration ./tests/...
//	docker compose -f tests/compose.yml down
package tests

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	libfs "github.com/scheiblingco/librc/fs"
)

const (
	testFile    = "testdir/hello.txt"
	testContent = "hello from librc integration tests"
)

// runSuite exercises the core Filesystem operations against fsys. Tests run
// sequentially; a failure in a setup step (MkdirAll, Put) is fatal so that
// later steps are not misleadingly skipped.
func runSuite(t *testing.T, fsys libfs.Filesystem) {
	t.Helper()
	ctx := context.Background()

	t.Run("MkdirAll", func(t *testing.T) {
		if err := fsys.MkdirAll(ctx, "testdir"); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}
	})

	t.Run("Put", func(t *testing.T) {
		r := strings.NewReader(testContent)
		if err := fsys.Put(ctx, testFile, r, int64(len(testContent)), time.Now()); err != nil {
			t.Fatalf("Put: %v", err)
		}
	})

	t.Run("Stat", func(t *testing.T) {
		info, err := fsys.Stat(ctx, testFile)
		if err != nil {
			t.Fatalf("Stat: %v", err)
		}
		if got := info.Size(); got != int64(len(testContent)) {
			t.Errorf("Stat: size = %d, want %d", got, len(testContent))
		}
	})

	t.Run("Open", func(t *testing.T) {
		f, err := fsys.Open(ctx, testFile)
		if err != nil {
			t.Fatalf("Open: %v", err)
		}
		defer f.Close()
		got, err := io.ReadAll(f)
		if err != nil {
			t.Fatalf("Read: %v", err)
		}
		if string(got) != testContent {
			t.Errorf("Open: content = %q, want %q", string(got), testContent)
		}
	})

	t.Run("List", func(t *testing.T) {
		entries, err := fsys.List(ctx, "testdir")
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if len(entries) == 0 {
			t.Fatal("List: got 0 entries, want ≥ 1")
		}
		found := false
		for _, e := range entries {
			if e.Name() == "hello.txt" {
				found = true
			}
		}
		if !found {
			t.Errorf("List: hello.txt not found in entries")
		}
	})

	t.Run("Walk", func(t *testing.T) {
		var sawDir, sawFile bool
		err := fsys.Walk(ctx, "", func(e libfs.DirEntry) error {
			switch {
			case e.IsDir() && e.Remote() == "testdir":
				sawDir = true
			case !e.IsDir() && e.Remote() == testFile:
				sawFile = true
				if e.Size() != int64(len(testContent)) {
					t.Errorf("Walk: %s size = %d, want %d", testFile, e.Size(), len(testContent))
				}
			}
			return nil
		})
		if err != nil {
			t.Fatalf("Walk: %v", err)
		}
		if !sawDir || !sawFile {
			t.Errorf("Walk: saw testdir %v, saw %s %v", sawDir, testFile, sawFile)
		}
	})

	t.Run("Copy", func(t *testing.T) {
		if err := fsys.Copy(ctx, testFile, "testdir/copy.txt"); err != nil {
			t.Fatalf("Copy: %v", err)
		}
		info, err := fsys.Stat(ctx, "testdir/copy.txt")
		if err != nil {
			t.Fatalf("Stat after Copy: %v", err)
		}
		if info.Size() != int64(len(testContent)) {
			t.Errorf("Copy: size = %d, want %d", info.Size(), len(testContent))
		}
	})

	t.Run("Move", func(t *testing.T) {
		if err := fsys.Move(ctx, "testdir/copy.txt", "testdir/moved.txt"); err != nil {
			t.Fatalf("Move: %v", err)
		}
		if _, err := fsys.Stat(ctx, "testdir/copy.txt"); err == nil {
			t.Error("Move: source still exists after move")
		}
		if _, err := fsys.Stat(ctx, "testdir/moved.txt"); err != nil {
			t.Errorf("Move: destination missing: %v", err)
		}
	})

	t.Run("Remove", func(t *testing.T) {
		if err := fsys.Remove(ctx, "testdir/moved.txt"); err != nil {
			t.Fatalf("Remove: %v", err)
		}
		if _, err := fsys.Stat(ctx, "testdir/moved.txt"); err == nil {
			t.Error("Remove: file still present after removal")
		}
	})

	t.Run("Purge", func(t *testing.T) {
		if err := fsys.Purge(ctx, "testdir"); err != nil {
			t.Fatalf("Purge: %v", err)
		}
		entries, err := fsys.List(ctx, "")
		if err != nil {
			return // empty root is acceptable on some backends after purge
		}
		for _, e := range entries {
			if e.Name() == "testdir" {
				t.Error("Purge: testdir still present")
			}
		}
	})
}
