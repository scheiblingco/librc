package ftp

import "time"

type Config struct {
	// The rclone remote name used internally (e.g. "myftp").
	Name string
	// The FTP server hostname or IP address.
	Host string
	// The FTP port number (default: 21).
	Port int
	// The FTP username.
	User string
	// The	 FTP password.
	Password string
	// TLS enables implicit FTPS (FTP over TLS on the control connection).
	TLS bool
	// Enables explicit FTPS (STARTTLS upgrade on port 21).
	ExplicitTLS bool
	// Limits the number of simultaneous FTP connections (0 = unlimited).
	Concurrency int
	// Closes idle connections after this duration (default: 60s).
	IdleTimeout time.Duration
	// The path within the FTP server to use as the filesystem root.
	// Defaults to the user's home directory when empty.
	Root string
}
