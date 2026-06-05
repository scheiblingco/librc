package fs

import (
	"context"
	"io"
	"os"
	"time"
)

// File represents an open remote object for reading. It is compatible with
// os.File for sequential read operations.
type File interface {
	io.ReadCloser

	Name() string
	Stat() (os.FileInfo, error)
}

// DirEntry describes a file or directory within a Filesystem. It is compatible
// with os.DirEntry and adds the remote path, size, and modification time
// provided by rclone backends.
type DirEntry interface {
	Name() string
	IsDir() bool
	Size() int64
	ModTime(ctx context.Context) time.Time
	Info() (os.FileInfo, error)
	Type() os.FileMode

	// Remote returns the full remote path of the entry relative to the filesystem root (e.g. "subdir/file.txt").
	Remote() string
}
