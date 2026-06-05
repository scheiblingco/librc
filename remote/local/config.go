package local

type Config struct {
	Name string

	// The root directory on the local filesystem.
	RootPath string

	// Sets symbolic links to be followed and their targets copied.
	CopyLinks bool

	// Silently skip symbolic links instead of reporting an error.
	SkipLinks bool

	// Prevent crossing filesystem boundaries (mount points).
	OneFileSystem bool

	// Disable setting the modification time after upload.
	NoSetModtime bool
}
