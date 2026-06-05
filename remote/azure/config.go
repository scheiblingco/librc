// Package azure provides a librc Filesystem backed by Azure Blob Storage,
// wrapping rclone's azureblob backend.
package azure

// Specifies the Azure Blob Storage access tier.
type AccessTier string

const (
	AccessTierHot     AccessTier = "hot"
	AccessTierCool    AccessTier = "cool"
	AccessTierCold    AccessTier = "cold"
	AccessTierArchive AccessTier = "archive"
)

// Holds the configuration for an Azure Blob Storage backend.
// Authentication priority: ConnectionString > SASUrl > Account+Key.
type Config struct {
	// The rclone remote name used internally (e.g. "myazure").
	Name string
	// The Azure Blob container name used as the filesystem root.
	Container string
	// The Azure storage account name.
	Account string
	// The storage account access key.
	Key string
	// A Shared Access Signature URL (alternative to Account+Key).
	SASUrl string
	// An Azure storage connection string (takes priority over Account+Key).
	ConnectionString string
	// Overrides the Azure Blob endpoint URL.
	Endpoint string
	// Sets the access tier for newly uploaded blobs.
	AccessTier AccessTier
}
