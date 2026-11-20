package azurefiles

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

// ListShares pages through the service listing with the account's shared key;
// a fake service stands in for Azure Files (Azurite doesn't emulate it).
func TestListShares(t *testing.T) {
	pages := map[string]string{"": "zeta", "p2": "alpha"}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("comp") != "list" || !strings.HasPrefix(r.Header.Get("Authorization"), "SharedKey acct:") {
			http.Error(w, "unexpected request "+r.URL.String(), http.StatusBadRequest)
			return
		}
		next := ""
		if q.Get("marker") == "" {
			next = "p2"
		}
		fmt.Fprintf(w, `<?xml version="1.0" encoding="utf-8"?><EnumerationResults><Shares><Share><Name>%s</Name></Share></Shares><NextMarker>%s</NextMarker></EnumerationResults>`,
			pages[q.Get("marker")], next)
	}))
	defer srv.Close()

	shares, err := ListShares(context.Background(), Config{
		Name:             "list-test",
		ShareName:        "ignored",
		ConnectionString: "DefaultEndpointsProtocol=http;AccountName=acct;AccountKey=a2V5;FileEndpoint=" + srv.URL + "/acct;",
	})
	if err != nil || !slices.Equal(shares, []string{"alpha", "zeta"}) {
		t.Fatalf("ListShares: %v %v", shares, err)
	}
}
