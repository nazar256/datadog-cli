package monitor

import (
	"testing"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

func TestMapMonitor(t *testing.T) {
	created := time.Date(2026, 3, 21, 10, 0, 0, 0, time.UTC)
	modified := created.Add(time.Hour)
	priority := int64(2)
	state := datadogV1.MONITOROVERALLSTATES_ALERT
	monitorType := datadogV1.MONITORTYPE_QUERY_ALERT
	item := datadogV1.Monitor{Query: "avg:test{*} > 1", Type: monitorType, Tags: []string{"env:prod"}}
	item.SetId(123)
	item.SetName("CPU high")
	item.SetMessage("Investigate")
	item.SetCreated(created)
	item.SetModified(modified)
	item.SetOverallState(state)
	item.Priority.Set(&priority)

	view := mapMonitor(item)
	if view.ID != 123 || view.Name != "CPU high" || view.State != "Alert" {
		t.Fatalf("unexpected view: %+v", view)
	}
	if view.Priority == nil || *view.Priority != 2 {
		t.Fatalf("unexpected priority: %+v", view.Priority)
	}
	if view.CreatedAt == nil || view.ModifiedAt == nil {
		t.Fatalf("expected timestamps to be set: %+v", view)
	}
}
