package cryptohopper

import (
	"context"
	"fmt"
	"net/url"
)

// ChartAPI is the resource namespace for saved chart layouts + shared charts.
type ChartAPI struct {
	client *Client
}

// List returns the user's saved charts. Requires ``read``.
func (c *ChartAPI) List(ctx context.Context) ([]map[string]any, error) {
	var out []map[string]any
	if err := c.client.request(ctx, "GET", "/chart/list", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get fetches a single saved chart. Requires ``read``.
func (c *ChartAPI) Get(ctx context.Context, chartID any) (map[string]any, error) {
	q := url.Values{"chart_id": {fmt.Sprint(chartID)}}
	out := map[string]any{}
	if err := c.client.request(ctx, "GET", "/chart/get", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Save saves a new chart layout. Requires ``manage``.
func (c *ChartAPI) Save(ctx context.Context, body map[string]any) (map[string]any, error) {
	out := map[string]any{}
	if err := c.client.request(ctx, "POST", "/chart/save", nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete removes a saved chart. Requires ``manage``.
func (c *ChartAPI) Delete(ctx context.Context, chartID any) error {
	return c.client.request(ctx, "POST", "/chart/delete", nil, map[string]any{"chart_id": chartID}, nil)
}

// ShareSave saves a shared (public-link) chart. Requires ``manage``.
func (c *ChartAPI) ShareSave(ctx context.Context, body map[string]any) (map[string]any, error) {
	out := map[string]any{}
	if err := c.client.request(ctx, "POST", "/chart/share-save", nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ShareGet fetches a shared chart by its share id / key. Public.
func (c *ChartAPI) ShareGet(ctx context.Context, shareID string) (map[string]any, error) {
	q := url.Values{"share_id": {shareID}}
	out := map[string]any{}
	if err := c.client.request(ctx, "GET", "/chart/share-get", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
