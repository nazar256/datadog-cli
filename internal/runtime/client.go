package runtime

import (
	"context"
	"fmt"
	"net/http"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
)

type Client struct {
	API    *datadog.APIClient
	Ctx    context.Context
	Config Config
}

func NewClient(parent context.Context, cfg Config) (*Client, error) {
	if err := cfg.RequireAuth(); err != nil {
		return nil, err
	}
	if parent == nil {
		parent = context.Background()
	}

	configuration := datadog.NewConfiguration()
	configuration.HTTPClient = &http.Client{Timeout: cfg.Timeout}
	version := cfg.Version
	if version == "" {
		version = "dev"
	}
	configuration.UserAgent = "ddog-cli/" + version + " (" + cfg.Site + ")"

	apiClient := datadog.NewAPIClient(configuration)
	ctx := context.WithValue(parent, datadog.ContextAPIKeys, map[string]datadog.APIKey{
		"apiKeyAuth": {Key: cfg.APIKey},
		"appKeyAuth": {Key: cfg.AppKey},
	})
	ctx = context.WithValue(ctx, datadog.ContextServerVariables, map[string]string{"site": cfg.Site})

	return &Client{API: apiClient, Ctx: ctx, Config: cfg}, nil
}

func WrapAPIError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("datadog API request failed: %w", err)
}
