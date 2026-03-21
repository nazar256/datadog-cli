package logs

import (
	"testing"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
)

func TestMapEntry(t *testing.T) {
	ts := time.Date(2026, 3, 21, 12, 0, 0, 0, time.UTC)
	attrs := datadogV2.NewLogAttributes()
	attrs.SetTimestamp(ts)
	attrs.SetService("api")
	attrs.SetStatus("error")
	attrs.SetHost("web-01")
	attrs.SetMessage("boom")
	attrs.SetTags([]string{"env:prod"})
	entry := datadogV2.NewLog()
	entry.SetId("abc")
	entry.SetAttributes(*attrs)

	view := mapEntry(*entry)
	if view.ID != "abc" || view.Service != "api" || view.Message != "boom" {
		t.Fatalf("unexpected view: %+v", view)
	}
}
