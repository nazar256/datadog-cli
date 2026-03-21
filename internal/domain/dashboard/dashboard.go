package dashboard

import (
	"context"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	cliruntime "github.com/nazar256/datadog-cli/internal/runtime"
)

type Service interface {
	List(context.Context, cliruntime.Config, ListParams) (ListResult, error)
	Get(context.Context, cliruntime.Config, string) (Detail, error)
}

type LiveService struct{}

type ListParams struct {
	Count          int64
	Start          int64
	IncludeShared  bool
	IncludeDeleted bool
}

type Summary struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	LayoutType string     `json:"layout_type,omitempty"`
	Author     string     `json:"author,omitempty"`
	URL        string     `json:"url,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	ModifiedAt *time.Time `json:"modified_at,omitempty"`
}

type Detail struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	LayoutType  string     `json:"layout_type,omitempty"`
	Author      string     `json:"author,omitempty"`
	URL         string     `json:"url,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	ModifiedAt  *time.Time `json:"modified_at,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	WidgetCount int        `json:"widget_count"`
	NotifyList  []string   `json:"notify_list,omitempty"`
}

type ListResult struct {
	Items []Summary `json:"items"`
	Count int       `json:"count"`
}

func (LiveService) List(ctx context.Context, cfg cliruntime.Config, params ListParams) (ListResult, error) {
	client, err := cliruntime.NewClient(ctx, cfg)
	if err != nil {
		return ListResult{}, err
	}
	api := datadogV1.NewDashboardsApi(client.API)
	opt := datadogV1.NewListDashboardsOptionalParameters()
	if params.Count > 0 {
		opt.WithCount(params.Count)
	}
	if params.Start > 0 {
		opt.WithStart(params.Start)
	}
	if params.IncludeShared {
		opt.WithFilterShared(true)
	}
	if params.IncludeDeleted {
		opt.WithFilterDeleted(true)
	}
	resp, _, err := api.ListDashboards(client.Ctx, *opt)
	if err != nil {
		return ListResult{}, cliruntime.WrapAPIError(err)
	}
	items := resp.GetDashboards()
	views := make([]Summary, 0, len(items))
	for _, item := range items {
		views = append(views, mapDashboardSummary(item))
	}
	return ListResult{Items: views, Count: len(views)}, nil
}

func (LiveService) Get(ctx context.Context, cfg cliruntime.Config, id string) (Detail, error) {
	client, err := cliruntime.NewClient(ctx, cfg)
	if err != nil {
		return Detail{}, err
	}
	api := datadogV1.NewDashboardsApi(client.API)
	item, _, err := api.GetDashboard(client.Ctx, id)
	if err != nil {
		return Detail{}, cliruntime.WrapAPIError(err)
	}
	return mapDashboardDetail(item), nil
}

func mapDashboardSummary(item datadogV1.DashboardSummaryDefinition) Summary {
	view := Summary{
		ID:         item.GetId(),
		Title:      item.GetTitle(),
		LayoutType: string(item.GetLayoutType()),
		Author:     item.GetAuthorHandle(),
		URL:        item.GetUrl(),
	}
	if item.HasCreatedAt() {
		created := item.GetCreatedAt().UTC()
		view.CreatedAt = &created
	}
	if item.HasModifiedAt() {
		modified := item.GetModifiedAt().UTC()
		view.ModifiedAt = &modified
	}
	return view
}

func mapDashboardDetail(item datadogV1.Dashboard) Detail {
	view := Detail{
		ID:          item.GetId(),
		Title:       item.Title,
		Description: item.GetDescription(),
		LayoutType:  string(item.LayoutType),
		Author:      item.GetAuthorHandle(),
		URL:         item.GetUrl(),
		WidgetCount: len(item.Widgets),
	}
	if item.HasCreatedAt() {
		created := item.GetCreatedAt().UTC()
		view.CreatedAt = &created
	}
	if item.HasModifiedAt() {
		modified := item.GetModifiedAt().UTC()
		view.ModifiedAt = &modified
	}
	if tags := item.Tags.Get(); tags != nil {
		view.Tags = append([]string{}, (*tags)...)
	}
	if notifyList := item.NotifyList.Get(); notifyList != nil {
		view.NotifyList = append([]string{}, (*notifyList)...)
	}
	return view
}
