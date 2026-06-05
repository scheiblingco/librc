// Provides cross-backend file and directory copy/move operations.
//
// When both source and destination backends implement the Unwrapper interface
// (all librc backends do), rclone's native operations.CopyFile / operations.MoveFile
// are used, which support server-side copy where the provider allows it and
// optimised streaming otherwise. When a backend does not implement Unwrapper the
// file is streamed through the local process via Open + Put.
package transfer

import (
	"context"
	"fmt"
	"path"

	rclonefs "github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/operations"

	libfs "github.com/scheiblingco/librc/fs"
)

// Unwrapper is optionally implemented by librc backends to expose the underlying
// rclone Fs. All built-in backends implement this; third-party backends may not.
type Unwrapper interface {
	Unwrap() rclonefs.Fs
}

// Copies the file at srcPath on src to dstPath on dst.
//
// If both backends implement Unwrapper, rclone's native transfer is used
// (server-side copy where the provider supports it, optimised streaming
// otherwise). If either backend does not implement Unwrapper the file is
// streamed through the local process.
func CopyFile(ctx context.Context, src libfs.Filesystem, srcPath string, dst libfs.Filesystem, dstPath string) error {
	su, srcOK := src.(Unwrapper)
	du, dstOK := dst.(Unwrapper)
	if srcOK && dstOK {
		if err := operations.CopyFile(ctx, du.Unwrap(), su.Unwrap(), dstPath, srcPath); err != nil {
			return fmt.Errorf("transfer: copy %s → %s: %w", srcPath, dstPath, err)
		}
		return nil
	}
	return streamCopyFile(ctx, src, srcPath, dst, dstPath)
}

// Moves the file at srcPath on src to dstPath on dst.
//
// If both backends implement Unwrapper, rclone's native move is used.
// Otherwise the file is streamed to dst and then removed from src.
func MoveFile(ctx context.Context, src libfs.Filesystem, srcPath string, dst libfs.Filesystem, dstPath string) error {
	su, srcOK := src.(Unwrapper)
	du, dstOK := dst.(Unwrapper)
	if srcOK && dstOK {
		if err := operations.MoveFile(ctx, du.Unwrap(), su.Unwrap(), dstPath, srcPath); err != nil {
			return fmt.Errorf("transfer: move %s → %s: %w", srcPath, dstPath, err)
		}
		return nil
	}
	if err := streamCopyFile(ctx, src, srcPath, dst, dstPath); err != nil {
		return err
	}
	return src.Remove(ctx, srcPath)
}

// Copies all files under srcDir on src to dstDir on dst recursively.
//
// Each file is transferred with CopyFile, so native rclone operations are used
// when both backends support Unwrapper. Directory entries are created on dst as
// needed. srcDir and dstDir may be empty strings to refer to the filesystem root.
func CopyDir(ctx context.Context, src libfs.Filesystem, srcDir string, dst libfs.Filesystem, dstDir string) error {
	entries, err := src.List(ctx, srcDir)
	if err != nil {
		return fmt.Errorf("transfer: list %q: %w", srcDir, err)
	}
	for _, e := range entries {
		dstPath := path.Join(dstDir, e.Name())
		if e.IsDir() {
			if err := dst.MkdirAll(ctx, dstPath); err != nil {
				return fmt.Errorf("transfer: mkdir %q: %w", dstPath, err)
			}
			if err := CopyDir(ctx, src, e.Remote(), dst, dstPath); err != nil {
				return err
			}
		} else {
			if err := CopyFile(ctx, src, e.Remote(), dst, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// Moves all files under srcDir on src to dstDir on dst recursively,
// then removes srcDir from src.
//
// Each file is transferred with MoveFile. If any transfer fails, already-moved
// files are not rolled back.
func MoveDir(ctx context.Context, src libfs.Filesystem, srcDir string, dst libfs.Filesystem, dstDir string) error {
	entries, err := src.List(ctx, srcDir)
	if err != nil {
		return fmt.Errorf("transfer: list %q: %w", srcDir, err)
	}
	for _, e := range entries {
		dstPath := path.Join(dstDir, e.Name())
		if e.IsDir() {
			if err := dst.MkdirAll(ctx, dstPath); err != nil {
				return fmt.Errorf("transfer: mkdir %q: %w", dstPath, err)
			}
			if err := MoveDir(ctx, src, e.Remote(), dst, dstPath); err != nil {
				return err
			}
		} else {
			if err := MoveFile(ctx, src, e.Remote(), dst, dstPath); err != nil {
				return err
			}
		}
	}
	return src.Purge(ctx, srcDir)
}

// Transfers a single file from src to dst via Open + Put,
// streaming the bytes through the local process.
func streamCopyFile(ctx context.Context, src libfs.Filesystem, srcPath string, dst libfs.Filesystem, dstPath string) error {
	info, err := src.Stat(ctx, srcPath)
	if err != nil {
		return fmt.Errorf("transfer: stat %q: %w", srcPath, err)
	}
	f, err := src.Open(ctx, srcPath)
	if err != nil {
		return fmt.Errorf("transfer: open %q: %w", srcPath, err)
	}
	defer f.Close()
	if err := dst.Put(ctx, dstPath, f, info.Size(), info.ModTime()); err != nil {
		return fmt.Errorf("transfer: put %q: %w", dstPath, err)
	}
	return nil
}
