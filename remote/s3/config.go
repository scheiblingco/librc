package s3

// Provider identifies the S3-compatible storage provider.
type Provider string

const (
	ProviderAWS          Provider = "AWS"
	ProviderMinio        Provider = "Minio"
	ProviderWasabi       Provider = "Wasabi"
	ProviderCloudflareR2 Provider = "Cloudflare"
	ProviderOther        Provider = "Other"
)

// StorageClass specifies the S3 storage class applied to uploaded objects.
type StorageClass string

const (
	StorageClassDefault            StorageClass = ""
	StorageClassStandard           StorageClass = "STANDARD"
	StorageClassReducedRedundancy  StorageClass = "REDUCED_REDUNDANCY"
	StorageClassStandardIA         StorageClass = "STANDARD_IA"
	StorageClassOnezoneIA          StorageClass = "ONEZONE_IA"
	StorageClassGlacier            StorageClass = "GLACIER"
	StorageClassDeepArchive        StorageClass = "DEEP_ARCHIVE"
	StorageClassIntelligentTiering StorageClass = "INTELLIGENT_TIERING"
)

type Config struct {
	Name string

	// The S3 bucket name used as the filesystem root.
	Bucket string

	// Identifies the S3-compatible storage provider.
	Provider Provider

	// Reads AWS credentials from environment variables or the IAM role.
	// When true, AccessKeyID and SecretAccessKey are ignored.
	EnvAuth bool

	// The AWS access key ID.
	AccessKeyID string

	// The AWS secret access key.
	SecretAccessKey string

	// The AWS region (e.g. "us-east-1"). Required for AWS; optional for others.
	Region string

	// The S3 endpoint URL. Required for non-AWS providers such as MinIO.
	Endpoint string

	// The storage class for newly uploaded objects.
	StorageClass StorageClass

	// The canned ACL applied to objects (e.g. "private", "public-read").
	ACL string

	// Forces path-style S3 URLs (required for MinIO and some other providers).
	ForcePathStyle bool

	// Uses the listing's LastModified as each object's mtime instead of HEADing
	// every object for its mtime metadata (much faster for read-only scans; the
	// time is the upload time, not the source file's).
	UseServerModTime bool
}
