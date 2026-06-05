package rcloneutil

import (
	"fmt"

	rclonefs "github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/config/configmap"
)

// ApplyDefaults populates m with the registered default value of every option
// for the named rclone backend (e.g. "sftp", "ftp", "s3", "azureblob").
// Call this before applying user-supplied values so that options the user does
// not set still receive rclone's intended defaults rather than Go zero values.
//
// The backend must already be registered (i.e. its package imported) before
// this is called; all librc backend New() constructors do this because
// they import the corresponding rclone backend package directly.
func ApplyDefaults(backendName string, m configmap.Simple) error {
	ri, err := rclonefs.Find(backendName)
	if err != nil {
		return fmt.Errorf("rcloneutil: unknown backend %q: %w", backendName, err)
	}
	for _, opt := range ri.Options {
		if opt.Default != nil {
			m.Set(opt.Name, fmt.Sprint(opt.Default))
		}
	}
	return nil
}
