package azurefiles

// Holds the configuration for an Azure Files backend.
// Authentication priority: ConnectionString > SASUrl > Account+Key.
type Config struct {
	// The rclone remote name used internally (e.g. "myazurefiles").
	Name string
	// The Azure Files share name (e.g. "myshare").
	ShareName string
	// The subpath within the share to use as the filesystem root.
	Root string
	// The Azure storage account name.
	Account string
	// The storage account access key.
	Key string
	// An Azure storage connection string (takes priority over Account+Key).
	ConnectionString string
	// A Shared Access Signature URL (alternative to Account+Key).
	SASUrl string
	// Endpoint overrides the Azure Files endpoint URL.
	Endpoint string
}
