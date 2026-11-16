package webdav

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	rclonewebdav "github.com/rclone/rclone/backend/webdav"
	rclonefs "github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/config/configmap"
	"github.com/rclone/rclone/fs/config/obscure"
	"github.com/rclone/rclone/fs/operations"

	libfs "github.com/scheiblingco/librc/fs"
	"github.com/scheiblingco/librc/rcloneutil"
)

type Filesystem struct {
	inner rclonefs.Fs
	cfg   Config
}

// compile-time interface check
var _ libfs.Filesystem = (*Filesystem)(nil)

// New connects to the WebDAV server described by cfg. rclone's
// bearer_token_command option is deliberately not exposed: it runs a shell
// command taken from configuration.
func New(ctx context.Context, cfg Config) (*Filesystem, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("webdav: URL is required")
	}
	m := configmap.Simple{}
	if err := rcloneutil.ApplyDefaults("webdav", m); err != nil {
		return nil, err
	}
	m.Set("url", cfg.URL)
	vendor := cfg.Vendor
	if vendor == "" {
		vendor = VendorOther
	}
	m.Set("vendor", string(vendor))
	if cfg.User != "" {
		m.Set("user", cfg.User)
	}
	if cfg.Password != "" {
		obscured, err := obscure.Obscure(cfg.Password)
		if err != nil {
			return nil, fmt.Errorf("webdav: obscure password: %w", err)
		}
		m.Set("pass", obscured)
	}
	if cfg.BearerToken != "" {
		m.Set("bearer_token", cfg.BearerToken)
	}

	inner, err := rclonewebdav.NewFs(ctx, cfg.Name, cfg.Root, m)
	if err != nil {
		return nil, fmt.Errorf("webdav: %w", err)
	}
	return &Filesystem{inner: inner, cfg: cfg}, nil
}

func (f *Filesystem) Unwrap() rclonefs.Fs { return f.inner }

func (f *Filesystem) Name() string   { return f.inner.Name() }
func (f *Filesystem) Root() string   { return f.inner.Root() }
func (f *Filesystem) String() string { return f.inner.String() }

func (f *Filesystem) Stat(ctx context.Context, path string) (os.FileInfo, error) {
	obj, err := f.inner.NewObject(ctx, path)
	if err != nil {
		return nil, err
	}
	return rcloneutil.RcloneObjectFileInfo{Obj: obj}, nil
}

func (f *Filesystem) Open(ctx context.Context, path string) (libfs.File, error) {
	obj, err := f.inner.NewObject(ctx, path)
	if err != nil {
		return nil, err
	}
	rc, err := obj.Open(ctx)
	if err != nil {
		return nil, err
	}
	return &rcloneutil.RcloneFile{RC: rc, Obj: obj}, nil
}

func (f *Filesystem) Put(ctx context.Context, path string, in io.Reader, size int64, modTime time.Time) error {
	info := rcloneutil.RcloneUploadInfo{
		UploadFs:    f.inner,
		RemotePath:  path,
		FileSize:    size,
		FileModTime: modTime,
	}
	_, err := f.inner.Put(ctx, in, info)
	return err
}

func (f *Filesystem) Remove(ctx context.Context, path string) error {
	obj, err := f.inner.NewObject(ctx, path)
	if err != nil {
		return err
	}
	return obj.Remove(ctx)
}

func (f *Filesystem) Copy(ctx context.Context, src, dst string) error {
	srcObj, err := f.inner.NewObject(ctx, src)
	if err != nil {
		return fmt.Errorf("copy src: %w", err)
	}
	dstObj, _ := f.inner.NewObject(ctx, dst)
	_, err = operations.Copy(ctx, f.inner, dstObj, dst, srcObj)
	return err
}

func (f *Filesystem) Move(ctx context.Context, src, dst string) error {
	srcObj, err := f.inner.NewObject(ctx, src)
	if err != nil {
		return fmt.Errorf("move src: %w", err)
	}
	dstObj, _ := f.inner.NewObject(ctx, dst)
	_, err = operations.Move(ctx, f.inner, dstObj, dst, srcObj)
	return err
}

func (f *Filesystem) List(ctx context.Context, dir string) ([]libfs.DirEntry, error) {
	entries, err := f.inner.List(ctx, dir)
	if err != nil {
		return nil, err
	}
	wrapped := rcloneutil.WrapEntries(entries)
	result := make([]libfs.DirEntry, len(wrapped))
	for i, e := range wrapped {
		result[i] = e
	}
	return result, nil
}

func (f *Filesystem) Walk(ctx context.Context, dir string, fn func(libfs.DirEntry) error) error {
	return rcloneutil.Walk(ctx, f.inner, dir, fn)
}

func (f *Filesystem) Mkdir(ctx context.Context, dir string) error {
	return f.inner.Mkdir(ctx, dir)
}

func (f *Filesystem) MkdirAll(ctx context.Context, dir string) error {
	parts := strings.Split(strings.Trim(dir, "/"), "/")
	for i := range parts {
		if err := f.inner.Mkdir(ctx, strings.Join(parts[:i+1], "/")); err != nil {
			return err
		}
	}
	return nil
}

func (f *Filesystem) Rmdir(ctx context.Context, dir string) error {
	return f.inner.Rmdir(ctx, dir)
}

func (f *Filesystem) Purge(ctx context.Context, dir string) error {
	return operations.Purge(ctx, f.inner, dir)
}
