package dashboard

import (
	"testing"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

func TestMapDashboardSummaryAndDetail(t *testing.T) {
	created := time.Date(2026, 3, 21, 9, 0, 0, 0, time.UTC)
	modified := created.Add(time.Hour)
	title := "Ops dashboard"
	author := "ops@example.com"
	id := "abc-def"
	url := "/dashboard/abc-def"
	layout := datadogV1.DASHBOARDLAYOUTTYPE_ORDERED
	summary := datadogV1.DashboardSummaryDefinition{}
	summary.SetId(id)
	summary.SetTitle(title)
	summary.SetAuthorHandle(author)
	summary.SetUrl(url)
	summary.SetLayoutType(layout)
	summary.SetCreatedAt(created)
	summary.SetModifiedAt(modified)

	view := mapDashboardSummary(summary)
	if view.ID != id || view.CreatedAt == nil || view.ModifiedAt == nil {
		t.Fatalf("unexpected summary view: %+v", view)
	}

	detail := datadogV1.NewDashboard(layout, title, []datadogV1.Widget{})
	detail.SetId(id)
	detail.SetAuthorHandle(author)
	detail.SetUrl(url)
	detail.SetCreatedAt(created)
	detail.SetModifiedAt(modified)
	full := mapDashboardDetail(*detail)
	if full.ID != id || full.WidgetCount != 0 || full.CreatedAt == nil || full.ModifiedAt == nil {
		t.Fatalf("unexpected detail view: %+v", full)
	}
}
