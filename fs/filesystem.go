package fs

import (
	"context"
	"io"
	"os"
	"time"
)

// Filesystem is the central universal interface for all librc storage backends.
// Every backend package exposes a New() constructor returning a *Filesystem that
// implements this interface.
type Filesystem interface {
	// Return the name of this remote.
	Name() string

	// Return the root path within the backend.
	Root() string

	// Return a human-readable description.
	String() string

	// Return file information for the object at path.
	Stat(ctx context.Context, path string) (os.FileInfo, error)

	// Open the object at path for reading.
	Open(ctx context.Context, path string) (File, error)

	// Upload data from in to path. size is the exact byte count; pass -1
	// if unknown (backend permitting). modTime sets the modification time.
	Put(ctx context.Context, path string, in io.Reader, size int64, modTime time.Time) error

	// Remove the object at given path.
	Remove(ctx context.Context, path string) error

	// Copy src to dst within the same filesystem.
	Copy(ctx context.Context, src, dst string) error

	// Move src to dst within the same filesystem.
	Move(ctx context.Context, src, dst string) error

	// List the entries in the directory at dir.
	// Pass an empty string for the root directory.
	List(ctx context.Context, dir string) ([]DirEntry, error)

	// Create the directory at dir.
	Mkdir(ctx context.Context, dir string) error

	// Create dir and all missing parent directories.
	MkdirAll(ctx context.Context, dir string) error

	// Remove the empty directory at dir.
	Rmdir(ctx context.Context, dir string) error

	// Remove dir and all of its contents recursively.
	Purge(ctx context.Context, dir string) error
}
