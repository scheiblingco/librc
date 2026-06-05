// Package types provides sentinel errors shared across all librc backends.
package types

import "errors"

// Sentinel errors returned by Filesystem implementations.
var (
	ErrNotExist   = errors.New("object not found")
	ErrIsDir      = errors.New("is a directory")
	ErrNotDir     = errors.New("is not a directory")
	ErrNotEmpty   = errors.New("directory not empty")
	ErrPermission = errors.New("permission denied")
)
