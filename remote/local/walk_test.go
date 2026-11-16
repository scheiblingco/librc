package local

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"

	libfs "github.com/scheiblingco/librc/fs"
)

func TestWalk(t *testing.T) {
	root := t.TempDir()
	for _, d := range []string{"a/b/c", "empty"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	files := map[string]string{"top.txt": "1", "a/one.txt": "22", "a/b/c/deep.txt": "333"}
	for name, data := range files {
		if err := os.WriteFile(filepath.Join(root, name), []byte(data), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	fsys, err := New(context.Background(), Config{Name: "walk-test", RootPath: root})
	if err != nil {
		t.Fatal(err)
	}

	var dirs, got []string
	err = fsys.Walk(context.Background(), "", func(e libfs.DirEntry) error {
		if e.IsDir() {
			dirs = append(dirs, e.Remote())
			return nil
		}
		got = append(got, e.Remote())
		if want := int64(len(files[e.Remote()])); e.Size() != want {
			t.Errorf("%s: size %d, want %d", e.Remote(), e.Size(), want)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	slices.Sort(dirs)
	slices.Sort(got)
	if want := []string{"a", "a/b", "a/b/c", "empty"}; !slices.Equal(dirs, want) {
		t.Errorf("dirs %v, want %v", dirs, want)
	}
	if want := []string{"a/b/c/deep.txt", "a/one.txt", "top.txt"}; !slices.Equal(got, want) {
		t.Errorf("files %v, want %v", got, want)
	}

	// A subdirectory walks only what's below it
	got = nil
	if err := fsys.Walk(context.Background(), "a/b", func(e libfs.DirEntry) error {
		got = append(got, e.Remote())
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	slices.Sort(got)
	if want := []string{"a/b/c", "a/b/c/deep.txt"}; !slices.Equal(got, want) {
		t.Errorf("subdir walk %v, want %v", got, want)
	}

	// An error from fn stops the walk and comes back unchanged
	stop := errors.New("stop")
	n := 0
	err = fsys.Walk(context.Background(), "", func(libfs.DirEntry) error {
		n++
		return stop
	})
	if !errors.Is(err, stop) || n != 1 {
		t.Errorf("early stop: err %v after %d entries", err, n)
	}

	if err := fsys.Walk(context.Background(), "missing", func(libfs.DirEntry) error { return nil }); err == nil {
		t.Error("walking a missing directory succeeded")
	}
}
