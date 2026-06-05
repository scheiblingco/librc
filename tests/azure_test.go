package tests

import (
	"context"
	"testing"

	"github.com/scheiblingco/librc/remote/azure"
)

// TestAzure requires the azurite service from compose.yml running on localhost:10000.
// Uses the well-known Azurite development credentials.
func TestAzure(t *testing.T) {
	ctx := context.Background()

	// Azurite development connection string — safe to commit, not real credentials.
	const azuriteConnStr = "DefaultEndpointsProtocol=http;" +
		"AccountName=devstoreaccount1;" +
		"AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;" +
		"BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"

	fsys, err := azure.New(ctx, azure.Config{
		Name:             "azure-test",
		Container:        "librc-test",
		ConnectionString: azuriteConnStr,
	})
	if err != nil {
		t.Fatalf("azure.New: %v", err)
	}

	// Create the container (Azurite starts empty).
	if err := fsys.Mkdir(ctx, ""); err != nil {
		t.Fatalf("create container: %v", err)
	}

	runSuite(t, fsys)
}
