// Provides shared adapter types used by all librc backend
// packages to satisfy the fs.File, fs.DirEntry, and os.FileInfo interfaces in
// terms of rclone's internal types.
package rcloneutil

import (
	"context"
	"io"
	"os"
	"path"
	"time"

	rclonefs "github.com/rclone/rclone/fs"
	rclonehash "github.com/rclone/rclone/fs/hash"
)

// Implements fs.File by wrapping an rclone Object and its open ReadCloser.
type RcloneFile struct {
	RC  io.ReadCloser
	Obj rclonefs.Object
}

func (f *RcloneFile) Read(p []byte) (int, error) { return f.RC.Read(p) }
func (f *RcloneFile) Close() error               { return f.RC.Close() }
func (f *RcloneFile) Name() string               { return f.Obj.Remote() }
func (f *RcloneFile) Stat() (os.FileInfo, error) { return RcloneObjectFileInfo{Obj: f.Obj}, nil }

// Implements fs.DirEntry by wrapping an rclone DirEntry.
type RcloneDirEntry struct {
	Inner rclonefs.DirEntry
}

func (d *RcloneDirEntry) Name() string { return path.Base(d.Inner.Remote()) }
func (d *RcloneDirEntry) IsDir() bool {
	_, ok := d.Inner.(rclonefs.Directory)
	return ok
}
func (d *RcloneDirEntry) Type() os.FileMode {
	if d.IsDir() {
		return os.ModeDir
	}
	return 0
}
func (d *RcloneDirEntry) Info() (os.FileInfo, error) {
	if obj, ok := d.Inner.(rclonefs.Object); ok {
		return RcloneObjectFileInfo{Obj: obj}, nil
	}
	return RcloneDirFileInfo{Dir: d.Inner}, nil
}
func (d *RcloneDirEntry) Remote() string                        { return d.Inner.Remote() }
func (d *RcloneDirEntry) Size() int64                           { return d.Inner.Size() }
func (d *RcloneDirEntry) ModTime(ctx context.Context) time.Time { return d.Inner.ModTime(ctx) }

// Implements os.FileInfo for an rclone Object.
type RcloneObjectFileInfo struct {
	Obj rclonefs.Object
}

func (i RcloneObjectFileInfo) Name() string       { return path.Base(i.Obj.Remote()) }
func (i RcloneObjectFileInfo) Size() int64        { return i.Obj.Size() }
func (i RcloneObjectFileInfo) Mode() os.FileMode  { return 0644 }
func (i RcloneObjectFileInfo) ModTime() time.Time { return i.Obj.ModTime(context.Background()) }
func (i RcloneObjectFileInfo) IsDir() bool        { return false }
func (i RcloneObjectFileInfo) Sys() any           { return i.Obj }

// Implements os.FileInfo for an rclone Directory entry.
type RcloneDirFileInfo struct {
	Dir rclonefs.DirEntry
}

func (i RcloneDirFileInfo) Name() string       { return path.Base(i.Dir.Remote()) }
func (i RcloneDirFileInfo) Size() int64        { return i.Dir.Size() }
func (i RcloneDirFileInfo) Mode() os.FileMode  { return os.ModeDir | 0755 }
func (i RcloneDirFileInfo) ModTime() time.Time { return i.Dir.ModTime(context.Background()) }
func (i RcloneDirFileInfo) IsDir() bool        { return true }
func (i RcloneDirFileInfo) Sys() any           { return i.Dir }

// Implements rclonefs.ObjectInfo for use with rclonefs.Fs.Put.
// Carries the destination path, size, and modification time for an upload.
type RcloneUploadInfo struct {
	UploadFs    rclonefs.Fs
	RemotePath  string
	FileSize    int64
	FileModTime time.Time
}

func (i RcloneUploadInfo) Fs() rclonefs.Info                                         { return i.UploadFs }
func (i RcloneUploadInfo) String() string                                            { return i.RemotePath }
func (i RcloneUploadInfo) Remote() string                                            { return i.RemotePath }
func (i RcloneUploadInfo) ModTime(_ context.Context) time.Time                       { return i.FileModTime }
func (i RcloneUploadInfo) Size() int64                                               { return i.FileSize }
func (i RcloneUploadInfo) Hash(_ context.Context, _ rclonehash.Type) (string, error) { return "", nil }
func (i RcloneUploadInfo) Storable() bool                                            { return true }

// Converts a slice of rclone DirEntries into RcloneDirEntry values
// satisfying fs.DirEntry.
func WrapEntries(entries rclonefs.DirEntries) []*RcloneDirEntry {
	result := make([]*RcloneDirEntry, len(entries))
	for i, e := range entries {
		result[i] = &RcloneDirEntry{Inner: e}
	}
	return result
}
