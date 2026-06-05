package smb

import "time"

type Config struct {
	Name string

	// The SMB server hostname or IP address.
	Host string

	// The SMB port number (default: 445).
	Port int

	// The SMB username.
	User string

	// The SMB password.
	Password string

	// The Windows domain for authentication (optional).
	Domain string

	// The SMB share name (e.g. "share", "homes").
	// It is prepended to Root to form the rclone root path.
	Share string

	// A subpath within the share to use as the filesystem root.
	Root string

	// Closes idle connections after this duration.
	IdleTimeout time.Duration
}
