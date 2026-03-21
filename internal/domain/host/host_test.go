package host

import (
	"testing"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

func TestMatchesHost(t *testing.T) {
	item := datadogV1.Host{Aliases: []string{"alias-1"}}
	item.SetName("web-01")
	item.SetHostName("web-01.local")
	if !matchesHost(item, "alias-1") || !matchesHost(item, "web-01") || !matchesHost(item, "WEB-01.LOCAL") {
		t.Fatalf("expected host match")
	}
	if matchesHost(item, "db-01") {
		t.Fatalf("did not expect host match")
	}
}
