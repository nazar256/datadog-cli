package logs

import (
	"context"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	cliruntime "github.com/nazar256/datadog-cli/internal/runtime"
	"github.com/nazar256/datadog-cli/internal/timeutil"
)

type Service interface {
	Search(context.Context, cliruntime.Config, SearchParams) (SearchResult, error)
}

type LiveService struct{}

type SearchParams struct {
	Query   string
	Range   timeutil.Range
	Limit   int32
	Indexes []string
	SortAsc bool
}

type Entry struct {
	ID        string     `json:"id,omitempty"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
	Service   string     `json:"service,omitempty"`
	Status    string     `json:"status,omitempty"`
	Host      string     `json:"host,omitempty"`
	Message   string     `json:"message,omitempty"`
	Tags      []string   `json:"tags,omitempty"`
}

type SearchResult struct {
	Query string    `json:"query"`
	From  time.Time `json:"from"`
	To    time.Time `json:"to"`
	Items []Entry   `json:"items"`
	Count int       `json:"count"`
}

func (LiveService) Search(ctx context.Context, cfg cliruntime.Config, params SearchParams) (SearchResult, error) {
	client, err := cliruntime.NewClient(ctx, cfg)
	if err != nil {
		return SearchResult{}, err
	}
	api := datadogV2.NewLogsApi(client.API)
	opt := datadogV2.NewListLogsGetOptionalParameters().WithFilterQuery(params.Query).WithFilterFrom(params.Range.From).WithFilterTo(params.Range.To)
	if params.Limit > 0 {
		opt.WithPageLimit(params.Limit)
	}
	if len(params.Indexes) > 0 {
		opt.WithFilterIndexes(params.Indexes)
	}
	if params.SortAsc {
		opt.WithSort(datadogV2.LOGSSORT_TIMESTAMP_ASCENDING)
	} else {
		opt.WithSort(datadogV2.LOGSSORT_TIMESTAMP_DESCENDING)
	}
	resp, _, err := api.ListLogsGet(client.Ctx, *opt)
	if err != nil {
		return SearchResult{}, cliruntime.WrapAPIError(err)
	}
	items := resp.GetData()
	views := make([]Entry, 0, len(items))
	for _, item := range items {
		views = append(views, mapEntry(item))
	}
	return SearchResult{Query: params.Query, From: params.Range.From, To: params.Range.To, Items: views, Count: len(views)}, nil
}

func mapEntry(item datadogV2.Log) Entry {
	attrs := item.GetAttributes()
	view := Entry{
		ID:      item.GetId(),
		Service: attrs.GetService(),
		Status:  attrs.GetStatus(),
		Host:    attrs.GetHost(),
		Message: attrs.GetMessage(),
		Tags:    append([]string{}, attrs.GetTags()...),
	}
	if attrs.HasTimestamp() {
		timestamp := attrs.GetTimestamp().UTC()
		view.Timestamp = &timestamp
	}
	return view
}
