// Package hetzner provides a librc Filesystem backed by a Hetzner Storage Box,
// using rclone's sftp backend with Hetzner-specific defaults (port 23,
// hostname derived from the storage box username).
package hetzner

type Config struct {
	// The rclone remote name used internally (e.g. "hetzner").
	Name string
	// The Hetzner Storage Box username (e.g. "u123456").
	// The SFTP host defaults to {Username}.your-storagebox.de.
	Username string
	// Overrides the derived SFTP hostname. Useful for testing with a
	// mock SFTP server on localhost.
	Host string
	// Overrides the default Hetzner SFTP port (23). Useful for testing.
	Port int
	// Password for SFTP authentication (alternative to KeyFile).
	Password string
	// The path to a PEM-encoded SSH private key.
	KeyFile string
	// The path to an SSH known_hosts file for host verification.
	// Populate it with: ssh-keyscan -p 23 {Username}.your-storagebox.de
	KnownHostsFile string
	// Overrides the SSH login username with a Hetzner sub-account
	// name (e.g. "u123456-sub1"). When empty, Username is used.
	SubAccount string
	// Root is a path within the storage box to use as the filesystem root.
	Root string
}
