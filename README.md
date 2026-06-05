# librc
We found ourselves re-implementing the Rclone library into multiple projects on our end and wanted to create a reusable, easy-to-use interface for using various functionality from Rclone in other Go projects. This library is not meant to be a full replacement, just to be able to use the base functionality in Rclone more easily, and to assist some of our younger colleagues by statically typing the configuration options for the various backends.

All backends aren't implemented, if you're missing one you feel you need please open an issue or, better yet, a pull request with the implementation.

Implementing all backends separately is done to avoid having to import the entire Rclone library.

> [!NOTE]
> AI Usage Disclaimer: Github Copilot was used to assist while writing some of the code in this repository, but the development was driven by a real human* (*not like the TV show).

## Usage

Each backend lives in its own importable package under remote/. All backends implement the `fs.Filesystem` interface in fs/


## Roadmap
- [ ] Implement syncFolder functionality (etc.) for cross-backend sync operations
- [ ] Add support for custom backends
- [ ] Implement the rest of the rclone backends

### Local

```go
import (
    "context"
    "github.com/scheiblingco/librc/remote/local"
)

fs, err := local.New(ctx, local.Config{
    // Internal identifier, set to whatever you like
    // If you use multiple backends of the same type at the same time, this is required to distinguish them in logs and error messages.
    Name: "local",
    // RootPath is the local directory to use as the filesystem root.
    RootPath: "/var/data",

    // CopyLinks follows symlinks and copies the target file/directory. (default: false)
    // CopyLinks: false,
    // SkipLinks skips symlinks instead of reporting them as errors. (default: false)
    // SkipLinks: false,
    // OneFileSystem prevents crossing filesystem boundaries during recursive operations. (default: false)
    // OneFileSystem: false,
    // NoSetModtime disables setting the modification time on uploaded files. (default: false)
    // NoSetModtime: false,
})

// List directory
entries, err := fs.List(ctx, "")

// Read a file
f, err := fs.Open(ctx, "subdir/file.txt")
defer f.Close()
io.Copy(os.Stdout, f)

// Upload a file
data := strings.NewReader("hello world")
err = fs.Put(ctx, "subdir/hello.txt", data, int64(data.Len()), time.Now())
```

### FTP

```go
import "github.com/scheiblingco/librc/remote/ftp"

fs, err := ftp.New(ctx, ftp.Config{
    // Internal identifier, set to whatever you like
    // If you use multiple backends of the same type at the same time, this is required to distinguish them in logs and error messages.
    Name: "myftp",
    // Host is the FTP server hostname or IP address.
    Host: "ftp.example.com",
    // User is the FTP username.
    User: "ftpuser",
    // Password is the FTP password (passed in plain text; obscured internally).
    Password: os.Getenv("FTP_PASSWORD"),

    // Port is the FTP port number. (default: 21)
    // Port: 21,
    // TLS enables implicit FTPS (FTP-over-TLS, typically port 990). (default: false)
    // TLS: false,
    // ExplicitTLS enables explicit FTPS via STARTTLS on the control connection. (default: false)
    // ExplicitTLS: false,
    // Concurrency limits simultaneous FTP connections; 0 means unlimited. (default: 0)
    // Concurrency: 0,
    // IdleTimeout closes idle connections after this duration. (default: 1m)
    // IdleTimeout: 60 * time.Second,
    // Root is the path within the FTP server to treat as the filesystem root.
    // Defaults to the user's home directory when empty.
    // Root: "",
})

entries, err := fs.List(ctx, "uploads/")
err = fs.Put(ctx, "uploads/report.pdf", r, size, time.Now())
```

### SFTP

```go
import "github.com/scheiblingco/librc/remote/sftp"

fs, err := sftp.New(ctx, sftp.Config{
    // Internal identifier, set to whatever you like
    // If you use multiple backends of the same type at the same time, this is required to distinguish them in logs and error messages.
    Name: "mysftp",
    // Host is the SFTP server hostname or IP address.
    Host: "sftp.example.com",
    // User is the SSH username.
    User: "deploy",

    // Port is the SSH port number. (default: 22)
    // Port: 22,
    // Password is the SSH password. Used only when no key is configured.
    // Password: os.Getenv("SFTP_PASSWORD"),
    // KeyFile is the path to a PEM-encoded SSH private key file.
    KeyFile: "/home/deploy/.ssh/id_ed25519",
    // KeyPEM is a raw PEM-encoded private key (alternative to KeyFile).
    // KeyPEM: os.Getenv("SFTP_KEY_PEM"),
    // KnownHostsFile is the path to the SSH known_hosts file for host verification.
    // Populate with: ssh-keyscan -p 22 sftp.example.com >> ~/.ssh/known_hosts
    KnownHostsFile: "/home/deploy/.ssh/known_hosts",
    // Concurrency is the number of concurrent SFTP transfer operations. (default: 0, unlimited)
    // Concurrency: 0,
    // IdleTimeout closes idle connections after this duration. (default: 1m)
    // IdleTimeout: 60 * time.Second,
    // Root is the path within the SFTP server to treat as the filesystem root.
    // Root: "",
})

err = fs.MkdirAll(ctx, "releases/v1.2.3")
err = fs.Move(ctx, "releases/v1.2.2/app", "releases/v1.2.3/app")
err = fs.Purge(ctx, "releases/v1.0.0")
```

### S3

```go
import "github.com/scheiblingco/librc/remote/s3"

fs, err := s3.New(ctx, s3.Config{
    // Internal identifier, set to whatever you like
    // If you use multiple backends of the same type at the same time, this is required to distinguish them in logs and error messages.
    Name: "mys3",
    // Bucket is the S3 bucket name used as the filesystem root.
    Bucket: "my-bucket",
    // Provider identifies the S3-compatible storage provider.
    // Use s3.ProviderAWS, s3.ProviderMinio, s3.ProviderWasabi,
    // s3.ProviderCloudflareR2, or s3.ProviderOther.
    Provider: s3.ProviderAWS,
    // Region is the AWS region (e.g. "us-east-1"). Required for AWS.
    Region: "us-east-1",
    // AccessKeyID is the AWS access key ID.
    AccessKeyID: os.Getenv("AWS_ACCESS_KEY_ID"),
    // SecretAccessKey is the AWS secret access key.
    SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),

    // EnvAuth reads credentials from environment variables or the instance IAM role.
    // When true, AccessKeyID and SecretAccessKey are ignored. (default: false)
    // EnvAuth: false,
    // Endpoint overrides the S3 endpoint URL. Required for MinIO and other non-AWS providers.
    // Endpoint: "https://minio.example.com",
    // ForcePathStyle uses path-style S3 URLs (required for MinIO and some other providers). (default: false)
    // ForcePathStyle: false,
    // StorageClass sets the storage class for uploaded objects.
    // Options: s3.StorageClassStandard, s3.StorageClassStandardIA,
    //          s3.StorageClassOnezoneIA, s3.StorageClassGlacier,
    //          s3.StorageClassDeepArchive, s3.StorageClassIntelligentTiering,
    //          s3.StorageClassReducedRedundancy. (default: provider default)
    // StorageClass: s3.StorageClassStandard,
    // ACL is the canned ACL applied to uploaded objects (e.g. "private", "public-read"). (default: "")
    // ACL: "private",
})

entries, err := fs.List(ctx, "prefix/")
err = fs.Copy(ctx, "src/file.txt", "dst/file.txt")
err = fs.Purge(ctx, "old-prefix/")
```

#### MinIO example

```go
fs, err := s3.New(ctx, s3.Config{
    Name:            "minio",
    Bucket:          "my-bucket",
    Provider:        s3.ProviderMinio,
    AccessKeyID:     os.Getenv("MINIO_ACCESS_KEY"),
    SecretAccessKey: os.Getenv("MINIO_SECRET_KEY"),
    Endpoint:        "http://localhost:9000",
    ForcePathStyle:  true,
})
```

### Azure Blob Storage

```go
import "github.com/scheiblingco/librc/remote/azure"

fs, err := azure.New(ctx, azure.Config{
    // Internal identifier, set to whatever you like
    // If you use multiple backends of the same type at the same time, this is required to distinguish them in logs and error messages.
    Name: "myazure",
    // Container is the Azure Blob container used as the filesystem root.
    Container: "my-container",

    // Authentication: use one of ConnectionString, SASUrl, or Account+Key.
    // ConnectionString takes priority over all other auth options.
    // ConnectionString: os.Getenv("AZURE_STORAGE_CONNECTION_STRING"),
    // SASUrl is a Shared Access Signature URL (alternative to Account+Key).
    // SASUrl: os.Getenv("AZURE_SAS_URL"),
    // Account is the Azure storage account name.
    Account: os.Getenv("AZURE_STORAGE_ACCOUNT"),
    // Key is the storage account access key.
    Key: os.Getenv("AZURE_STORAGE_KEY"),

    // Endpoint overrides the Azure Blob endpoint URL (e.g. for Azurite emulator).
    // Endpoint: "http://localhost:10000/devstoreaccount1",
    // AccessTier sets the access tier for uploaded blobs.
    // Options: azure.AccessTierHot, azure.AccessTierCool,
    //          azure.AccessTierCold, azure.AccessTierArchive. (default: account default)
    // AccessTier: azure.AccessTierHot,
})

entries, err := fs.List(ctx, "")
err = fs.Put(ctx, "backups/dump.tar.gz", r, size, time.Now())
```

### Azure Files

```go
import "github.com/scheiblingco/librc/remote/azurefiles"

fs, err := azurefiles.New(ctx, azurefiles.Config{
    // Internal identifier, set to whatever you like
    // If you use multiple backends of the same type at the same time, this is required to distinguish them in logs and error messages.
    Name: "myazurefiles",
    // ShareName is the Azure Files share name.
    ShareName: "myshare",

    // Authentication: use one of ConnectionString, SASUrl, or Account+Key.
    // ConnectionString: os.Getenv("AZURE_STORAGE_CONNECTION_STRING"),
    // SASUrl: os.Getenv("AZURE_FILES_SAS_URL"),
    // Account is the Azure storage account name.
    Account: os.Getenv("AZURE_STORAGE_ACCOUNT"),
    // Key is the storage account access key.
    Key: os.Getenv("AZURE_STORAGE_KEY"),

    // Root is a subpath within the share to use as the filesystem root. (default: share root)
    // Root: "subdir",
    // Endpoint overrides the Azure Files endpoint URL (e.g. for Azurite emulator).
    // Endpoint: "http://localhost:10001/devstoreaccount1",
})

err = fs.MkdirAll(ctx, "team/documents")
err = fs.Put(ctx, "team/documents/report.docx", r, size, time.Now())
```

### SMB / CIFS

```go
import "github.com/scheiblingco/librc/remote/smb"

fs, err := smb.New(ctx, smb.Config{
    // Internal identifier, set to whatever you like
    // If you use multiple backends of the same type at the same time, this is required to distinguish them in logs and error messages.
    Name: "mysmb",
    // Host is the SMB server hostname or IP address.
    Host: "fileserver.corp.example.com",
    // User is the SMB username.
    User: "corpuser",
    // Password is the SMB password (passed in plain text; obscured internally).
    Password: os.Getenv("SMB_PASSWORD"),
    // Share is the SMB share name (required).
    Share: "data",

    // Port is the SMB port number. (default: 445)
    // Port: 445,
    // Domain is the Windows domain for authentication. (default: "")
    // Domain: "CORP",
    // Root is a subpath within the share to use as the filesystem root. (default: share root)
    // Root: "subdir",
    // IdleTimeout closes idle connections after this duration. (default: 1m)
    // IdleTimeout: 60 * time.Second,
})

entries, err := fs.List(ctx, "reports/")
err = fs.Copy(ctx, "reports/q1.xlsx", "archive/2025/q1.xlsx")
```

### Hetzner Storage Box

Hetzner Storage Boxes are accessed via SFTP on port 23. The hostname is derived automatically from the username (`{username}.your-storagebox.de`).

```go
import "github.com/scheiblingco/librc/remote/hetzner"

fs, err := hetzner.New(ctx, hetzner.Config{
    // Internal identifier, set to whatever you like
    // If you use multiple backends of the same type at the same time, this is required to distinguish them in logs and error messages.
    Name: "hetzner",
    // Username is the Hetzner Storage Box username (e.g. "u123456").
    // The SFTP host defaults to {Username}.your-storagebox.de.
    Username: "u123456",
    // Password for SFTP authentication.
    Password: os.Getenv("HETZNER_PASSWORD"),
    // KnownHostsFile is the path to an SSH known_hosts file for host verification.
    // Populate with: ssh-keyscan -p 23 u123456.your-storagebox.de >> ~/.ssh/known_hosts
    KnownHostsFile: "/home/user/.ssh/known_hosts",

    // KeyFile is the path to a PEM-encoded SSH private key (alternative to Password).
    // KeyFile: "/home/user/.ssh/id_ed25519",
    // SubAccount overrides the SSH login username with a Hetzner sub-account name
    // (e.g. "u123456-sub1"). When empty, Username is used. (default: "")
    // SubAccount: "",
    // Root is a path within the storage box to use as the filesystem root. (default: box root)
    // Root: "backups",
    // Host overrides the derived hostname. Useful for testing with a mock SFTP server.
    // Host: "localhost",
    // Port overrides the default SFTP port (23). Useful for testing.
    // Port: 23,
})

err = fs.MkdirAll(ctx, "backups/2025/06")
err = fs.Put(ctx, "backups/2025/06/dump.tar.gz", r, size, time.Now())
```

## Common interface

Because all backends implement `fs.Filesystem`, you can write backend-agnostic code:

```go
import (
    "context"
    "io"
    "time"

    libfs "github.com/scheiblingco/librc/fs"
)

func uploadFile(ctx context.Context, filesystem libfs.Filesystem, path string, r io.Reader, size int64) error {
    return filesystem.Put(ctx, path, r, size, time.Now())
}

func listRecursive(ctx context.Context, filesystem libfs.Filesystem, dir string) ([]libfs.DirEntry, error) {
    return filesystem.List(ctx, dir)
}
```

### Available methods

| Method | Signature |
|---|---|
| `Stat` | `(ctx, path) (os.FileInfo, error)` |
| `Open` | `(ctx, path) (File, error)` |
| `Put` | `(ctx, path, reader, size, modTime) error` |
| `Remove` | `(ctx, path) error` |
| `Copy` | `(ctx, src, dst) error` |
| `Move` | `(ctx, src, dst) error` |
| `List` | `(ctx, dir) ([]DirEntry, error)` |
| `Mkdir` | `(ctx, dir) error` |
| `MkdirAll` | `(ctx, dir) error` |
| `Rmdir` | `(ctx, dir) error` |
| `Purge` | `(ctx, dir) error` |

## Cross-backend transfers

The transfer package moves data between two backends. When both backends are built-in librc backends they expose an `Unwrap() rclonefs.Fs` method, transfer detects this and uses rclone's native CopyFile/MoveFile methods, which support server-side copy where the provider allows it. When a backend does not expose an rclone Fs (e.g. transfer from a stream to a backend), the data is streamed through the local process via Open & Put.

```go
import "github.com/scheiblingco/librc/transfer"
```

### CopyFile — copy a single file between backends

```go
sftp, _ := sftp.New(ctx, sftp.Config{ /* ... */ })
s3,   _ := s3.New(ctx,   s3.Config{  /* ... */ })

// Copy one file from SFTP to S3.
err := transfer.CopyFile(ctx, sftp, "exports/report.csv", s3, "archive/report.csv")
```

### MoveFile — move a single file between backends

```go
// Move a file from local disk to Azure Blob, deleting the local copy on success.
local, _ := local.New(ctx, local.Config{RootPath: "/var/data"})
azure, _ := azure.New(ctx, azure.Config{ /* ... */ })

err := transfer.MoveFile(ctx, local, "tmp/upload.zip", azure, "uploads/upload.zip")
```

### CopyDir — recursively copy a directory between backends

```go
// Mirror an S3 prefix to a Hetzner Storage Box.
s3,      _ := s3.New(ctx,      s3.Config{      /* ... */ })
hetzner, _ := hetzner.New(ctx, hetzner.Config{ /* ... */ })

// Copies everything under "backups/2025/" on S3 to "s3-mirror/2025/" on the storage box.
err := transfer.CopyDir(ctx, s3, "backups/2025/", hetzner, "s3-mirror/2025/")
```

### MoveDir — recursively move a directory between backends

```go
// Move processed files from SFTP to S3, removing the source directory when done.
err := transfer.MoveDir(ctx, sftp, "processed/", s3, "archive/processed/")
```

### Available Methods

| Function | Description |
|---|---|
| `CopyFile(ctx, src, srcPath, dst, dstPath)` | Copy one file from src to dst |
| `MoveFile(ctx, src, srcPath, dst, dstPath)` | Move one file from src to dst |
| `CopyDir(ctx, src, srcDir, dst, dstDir)` | Recursively copy a directory tree |
| `MoveDir(ctx, src, srcDir, dst, dstDir)` | Recursively move a directory tree |
