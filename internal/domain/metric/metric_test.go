package metric

import (
	"math"
	"testing"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

func TestMapSeriesUsesLastValidPoint(t *testing.T) {
	metricName := "system.load.1"
	scope := "host:web-01"
	item := datadogV1.MetricsQueryMetadata{Pointlist: [][]*float64{{ptr(1711010000000), nil}, {ptr(1711010060000), ptr(1.5)}, {ptr(1711010120000), ptrNaN()}}}
	item.SetMetric(metricName)
	item.SetScope(scope)
	item.SetInterval(60000)

	view := mapSeries(item)
	if view.LastValue == nil || *view.LastValue != 1.5 {
		t.Fatalf("unexpected last value: %+v", view.LastValue)
	}
	if !view.LastPointTS.Equal(time.UnixMilli(1711010060000).UTC()) {
		t.Fatalf("unexpected last point timestamp: %v", view.LastPointTS)
	}
	if view.PointCount != 1 {
		t.Fatalf("unexpected point count: %d", view.PointCount)
	}
}

func ptr(v float64) *float64 { return &v }
func ptrNaN() *float64 {
	v := math.NaN()
	return &v
}
