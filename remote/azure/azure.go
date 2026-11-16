package azure

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	rcloneazure "github.com/rclone/rclone/backend/azureblob"
	rclonefs "github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/config/configmap"
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

func New(ctx context.Context, cfg Config) (*Filesystem, error) {
	m, err := configMap(cfg)
	if err != nil {
		return nil, err
	}
	inner, err := rcloneazure.NewFs(ctx, cfg.Name, cfg.Container, m)
	if err != nil {
		return nil, fmt.Errorf("azure: %w", err)
	}
	return &Filesystem{inner: inner, cfg: cfg}, nil
}

// configMap turns cfg into rclone options (shared by New and ListContainers)
func configMap(cfg Config) (configmap.Simple, error) {
	m := configmap.Simple{}
	if err := rcloneutil.ApplyDefaults("azureblob", m); err != nil {
		return nil, err
	}
	if cfg.Account != "" {
		m.Set("account", cfg.Account)
	}
	if cfg.Key != "" {
		m.Set("key", cfg.Key)
	}
	if cfg.SASUrl != "" {
		m.Set("sas_url", cfg.SASUrl)
	}
	if cfg.ConnectionString != "" {
		m.Set("connection_string", cfg.ConnectionString)
	}
	if cfg.Endpoint != "" {
		m.Set("endpoint", cfg.Endpoint)
	}
	if cfg.AccessTier != "" {
		m.Set("access_tier", string(cfg.AccessTier))
	}
	return m, nil
}

// ListContainers lists the containers visible on the storage account; cfg.Container is ignored.
// Names come back sorted.
func ListContainers(ctx context.Context, cfg Config) ([]string, error) {
	m, err := configMap(cfg)
	if err != nil {
		return nil, err
	}
	f, err := rcloneazure.NewFs(ctx, cfg.Name, "", m)
	if err != nil {
		return nil, fmt.Errorf("azure: %w", err)
	}
	return rcloneutil.TopLevelDirs(ctx, f)
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
