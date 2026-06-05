package tests

import (
	"context"
	"os"
	"testing"

	"github.com/scheiblingco/librc/remote/azurefiles"
)

// TestAzureFiles requires a real Azure Files endpoint — Azurite does not
// emulate the Azure Files service. Set the following environment variables
// before running, or the test will be skipped:
//
//	AZURE_FILES_ACCOUNT   - storage account name
//	AZURE_FILES_KEY       - storage account key
//	AZURE_FILES_SHARE     - file share name (must already exist)
func TestAzureFiles(t *testing.T) {
	account := os.Getenv("AZURE_FILES_ACCOUNT")
	key := os.Getenv("AZURE_FILES_KEY")
	share := os.Getenv("AZURE_FILES_SHARE")
	if account == "" || key == "" || share == "" {
		t.Skip("set AZURE_FILES_ACCOUNT, AZURE_FILES_KEY, AZURE_FILES_SHARE to run this test")
	}

	fsys, err := azurefiles.New(context.Background(), azurefiles.Config{
		Name:      "azurefiles-test",
		ShareName: share,
		Account:   account,
		Key:       key,
	})
	if err != nil {
		t.Fatalf("azurefiles.New: %v", err)
	}
	runSuite(t, fsys)
}
