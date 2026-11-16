package rcloneutil

import (
	"context"
	"slices"

	rclonefs "github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/walk"

	libfs "github.com/scheiblingco/librc/fs"
)

// Walk calls fn for every file and directory below dir on f, recursively.
// It uses the backend's flat recursive listing where there is one (S3, Azure
// Blob and other bucket-based remotes) and lists directory by directory
// otherwise. On bucket-based remotes, directories are synthesized from the
// object paths. fn is never called concurrently, but entries come in no
// particular order: a child may be seen before its parent. Returning an error
// from fn stops the walk and Walk returns that error.
func Walk(ctx context.Context, f rclonefs.Fs, dir string, fn func(libfs.DirEntry) error) error {
	return walk.ListR(ctx, f, dir, true, -1, walk.ListAll, func(entries rclonefs.DirEntries) error {
		for _, e := range entries {
			if err := fn(&RcloneDirEntry{Inner: e}); err != nil {
				return err
			}
		}
		return nil
	})
}

// TopLevelDirs returns the sorted names of the directories at the root of f.
// Opened with an empty root, that's the buckets, containers or shares.
func TopLevelDirs(ctx context.Context, f rclonefs.Fs) ([]string, error) {
	entries, err := f.List(ctx, "")
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if _, ok := e.(rclonefs.Directory); ok {
			names = append(names, e.Remote())
		}
	}
	slices.Sort(names)
	return names, nil
}
