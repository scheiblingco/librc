package webdav

// Vendor identifies the WebDAV server software, which enables
// server-specific behaviour (modification times, checksums, chunked uploads).
type Vendor string

const (
	VendorNextcloud      Vendor = "nextcloud"
	VendorOwncloud       Vendor = "owncloud"
	VendorInfiniteScale  Vendor = "infinitescale"
	VendorSharepoint     Vendor = "sharepoint"
	VendorSharepointNTLM Vendor = "sharepoint-ntlm"
	VendorFastmail       Vendor = "fastmail"
	VendorRclone         Vendor = "rclone"
	VendorOther          Vendor = "other"
)

type Config struct {
	Name string

	// The URL of the WebDAV server (e.g. "https://cloud.example.com/remote.php/dav/files/alice").
	URL string

	// The server software (default: other).
	Vendor Vendor

	// The user name. For NTLM, use the form "Domain\User".
	User string

	// The password.
	Password string

	// A bearer token used instead of User and Password.
	BearerToken string

	// The path below URL to use as the filesystem root.
	Root string
}
