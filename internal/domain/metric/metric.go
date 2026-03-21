package metric

import (
	"context"
	"math"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	cliruntime "github.com/nazar256/datadog-cli/internal/runtime"
	"github.com/nazar256/datadog-cli/internal/timeutil"
)

type Service interface {
	Query(context.Context, cliruntime.Config, QueryParams) (QueryResult, error)
}

type LiveService struct{}

type QueryParams struct {
	Query string
	Range timeutil.Range
}

type Series struct {
	Metric      string     `json:"metric,omitempty"`
	Expression  string     `json:"expression,omitempty"`
	Scope       string     `json:"scope,omitempty"`
	Aggregator  string     `json:"aggregator,omitempty"`
	IntervalMS  int64      `json:"interval_ms,omitempty"`
	PointCount  int        `json:"point_count"`
	Start       *time.Time `json:"start,omitempty"`
	End         *time.Time `json:"end,omitempty"`
	LastPointTS *time.Time `json:"last_point_ts,omitempty"`
	LastValue   *float64   `json:"last_value,omitempty"`
}

type QueryResult struct {
	Query  string    `json:"query"`
	From   time.Time `json:"from"`
	To     time.Time `json:"to"`
	Status string    `json:"status,omitempty"`
	Series []Series  `json:"series"`
	Count  int       `json:"count"`
}

func (LiveService) Query(ctx context.Context, cfg cliruntime.Config, params QueryParams) (QueryResult, error) {
	client, err := cliruntime.NewClient(ctx, cfg)
	if err != nil {
		return QueryResult{}, err
	}
	api := datadogV1.NewMetricsApi(client.API)
	resp, _, err := api.QueryMetrics(client.Ctx, params.Range.From.Unix(), params.Range.To.Unix(), params.Query)
	if err != nil {
		return QueryResult{}, cliruntime.WrapAPIError(err)
	}
	series := make([]Series, 0, len(resp.GetSeries()))
	for _, item := range resp.GetSeries() {
		series = append(series, mapSeries(item))
	}
	return QueryResult{
		Query:  params.Query,
		From:   params.Range.From,
		To:     params.Range.To,
		Status: resp.GetStatus(),
		Series: series,
		Count:  len(series),
	}, nil
}

func mapSeries(item datadogV1.MetricsQueryMetadata) Series {
	view := Series{
		Metric:     item.GetMetric(),
		Expression: item.GetExpression(),
		Scope:      item.GetScope(),
		Aggregator: item.GetAggr(),
		IntervalMS: item.GetInterval(),
		PointCount: countPoints(item.Pointlist),
	}
	if item.HasStart() {
		start := time.UnixMilli(item.GetStart()).UTC()
		view.Start = &start
	}
	if item.HasEnd() {
		end := time.UnixMilli(item.GetEnd()).UTC()
		view.End = &end
	}
	if ts, value, ok := lastPoint(item.Pointlist); ok {
		view.LastPointTS = &ts
		view.LastValue = &value
	}
	return view
}

func countPoints(points [][]*float64) int {
	count := 0
	for _, point := range points {
		if len(point) >= 2 && point[1] != nil && !math.IsNaN(*point[1]) {
			count++
		}
	}
	return count
}

func lastPoint(points [][]*float64) (time.Time, float64, bool) {
	for i := len(points) - 1; i >= 0; i-- {
		point := points[i]
		if len(point) < 2 || point[0] == nil || point[1] == nil || math.IsNaN(*point[1]) {
			continue
		}
		return time.UnixMilli(int64(*point[0])).UTC(), *point[1], true
	}
	return time.Time{}, 0, false
}
