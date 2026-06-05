package sftp

import "time"

type Config struct {
	Name string

	// The SFTP server hostname or IP address.
	Host string

	// The SSH port number (default: 22).
	Port int

	// The SSH username.
	User string

	// The SSH password. Used only when no key is provided.
	Password string

	// The path to a PEM-encoded private key file.
	KeyFile string

	// A raw PEM-encoded private key (alternative to KeyFile).
	KeyPEM string

	// The path to the SSH known_hosts file for host verification.
	KnownHostsFile string

	// The number of concurrent SFTP transfer operations.
	Concurrency int

	// The duration after which idle connections are closed.
	IdleTimeout time.Duration

	// The path within the SFTP server to use as the filesystem root.
	// Defaults to the user's home directory when empty.
	Root string
}
