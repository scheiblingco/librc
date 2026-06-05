package hetzner

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	rclonesftp "github.com/rclone/rclone/backend/sftp"
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

// New creates a Filesystem connected to the Hetzner Storage Box described by cfg.
// Hetzner Storage Boxes are accessed via SFTP on port 23. The host defaults to
// {Username}.your-storagebox.de unless cfg.Host is set.
func New(ctx context.Context, cfg Config) (*Filesystem, error) {
	host := cfg.Host
	if host == "" {
		if cfg.Username == "" {
			return nil, fmt.Errorf("hetzner: Username is required when Host is not set")
		}
		host = cfg.Username + ".your-storagebox.de"
	}
	user := cfg.Username
	if cfg.SubAccount != "" {
		user = cfg.SubAccount
	}
	port := 23
	if cfg.Port != 0 {
		port = cfg.Port
	}

	m := configmap.Simple{}
	if err := rcloneutil.ApplyDefaults("sftp", m); err != nil {
		return nil, err
	}
	m.Set("host", host)
	m.Set("port", strconv.Itoa(port))
	if user != "" {
		m.Set("user", user)
	}
	if cfg.Password != "" {
		obscured, err := obscure.Obscure(cfg.Password)
		if err != nil {
			return nil, fmt.Errorf("hetzner: obscure password: %w", err)
		}
		m.Set("pass", obscured)
	}
	if cfg.KeyFile != "" {
		m.Set("key_file", cfg.KeyFile)
	}
	if cfg.KnownHostsFile != "" {
		m.Set("known_hosts_file", cfg.KnownHostsFile)
	}

	inner, err := rclonesftp.NewFs(ctx, cfg.Name, cfg.Root, m)
	if err != nil {
		return nil, fmt.Errorf("hetzner: %w", err)
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
