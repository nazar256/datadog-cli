package monitor

import (
	"context"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	cliruntime "github.com/nazar256/datadog-cli/internal/runtime"
	"github.com/samber/lo"
)

type Service interface {
	List(context.Context, cliruntime.Config, ListParams) (ListResult, error)
	Get(context.Context, cliruntime.Config, int64) (Detail, error)
}

type LiveService struct{}

type ListParams struct {
	Name        string
	Tags        string
	MonitorTags string
	Offset      int64
	Limit       int32
}

type Summary struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	State       string     `json:"state,omitempty"`
	Query       string     `json:"query"`
	Tags        []string   `json:"tags,omitempty"`
	Priority    *int64     `json:"priority,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	ModifiedAt  *time.Time `json:"modified_at,omitempty"`
	Message     string     `json:"message,omitempty"`
	MultiAlert  bool       `json:"multi_alert"`
	DraftStatus string     `json:"draft_status,omitempty"`
}

type Detail = Summary

type ListResult struct {
	Items []Summary `json:"items"`
	Count int       `json:"count"`
}

func (LiveService) List(ctx context.Context, cfg cliruntime.Config, params ListParams) (ListResult, error) {
	client, err := cliruntime.NewClient(ctx, cfg)
	if err != nil {
		return ListResult{}, err
	}
	api := datadogV1.NewMonitorsApi(client.API)
	opt := datadogV1.NewListMonitorsOptionalParameters()
	if params.Name != "" {
		opt.WithName(params.Name)
	}
	if params.Tags != "" {
		opt.WithTags(params.Tags)
	}
	if params.MonitorTags != "" {
		opt.WithMonitorTags(params.MonitorTags)
	}
	if params.Offset > 0 {
		opt.WithIdOffset(params.Offset)
	}
	if params.Limit > 0 {
		opt.WithPageSize(params.Limit)
	}
	items, _, err := api.ListMonitors(client.Ctx, *opt)
	if err != nil {
		return ListResult{}, cliruntime.WrapAPIError(err)
	}
	views := make([]Summary, 0, len(items))
	for _, item := range items {
		views = append(views, mapMonitor(item))
	}
	return ListResult{Items: views, Count: len(views)}, nil
}

func (LiveService) Get(ctx context.Context, cfg cliruntime.Config, id int64) (Detail, error) {
	client, err := cliruntime.NewClient(ctx, cfg)
	if err != nil {
		return Detail{}, err
	}
	api := datadogV1.NewMonitorsApi(client.API)
	item, _, err := api.GetMonitor(client.Ctx, id)
	if err != nil {
		return Detail{}, cliruntime.WrapAPIError(err)
	}
	return mapMonitor(item), nil
}

func mapMonitor(item datadogV1.Monitor) Summary {
	view := Summary{
		ID:         item.GetId(),
		Name:       item.GetName(),
		Type:       string(item.Type),
		State:      string(item.GetOverallState()),
		Query:      item.Query,
		Tags:       item.Tags,
		Message:    item.GetMessage(),
		MultiAlert: item.GetMulti(),
	}
	if item.HasCreated() {
		created := item.GetCreated().UTC()
		view.CreatedAt = &created
	}
	if item.HasModified() {
		modified := item.GetModified().UTC()
		view.ModifiedAt = &modified
	}
	if priority := item.Priority.Get(); priority != nil {
		p := *priority
		view.Priority = &p
	}
	if item.HasDraftStatus() {
		view.DraftStatus = string(item.GetDraftStatus())
	}
	view.Tags = lo.Map(item.Tags, func(tag string, _ int) string { return tag })
	return view
}
