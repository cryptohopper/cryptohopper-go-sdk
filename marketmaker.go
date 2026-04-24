package cryptohopper

import (
	"context"
	"fmt"
	"net/url"
)

// MarketMakerAPI is the resource namespace for market-maker bot ops + market-
// trend overrides + backlog.
type MarketMakerAPI struct {
	client *Client
}

// Get fetches the market-maker state for a hopper. Requires ``read``.
func (m *MarketMakerAPI) Get(ctx context.Context, extra url.Values) (map[string]any, error) {
	out := map[string]any{}
	if err := m.client.request(ctx, "GET", "/marketmaker/get", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Cancel cancels running market-maker orders. Requires ``trade``.
func (m *MarketMakerAPI) Cancel(ctx context.Context, body map[string]any) error {
	if body == nil {
		body = map[string]any{}
	}
	return m.client.request(ctx, "POST", "/marketmaker/cancel", nil, body, nil)
}

// History returns historical order activity. Requires ``read``.
func (m *MarketMakerAPI) History(ctx context.Context, extra url.Values) ([]map[string]any, error) {
	var out []map[string]any
	if err := m.client.request(ctx, "GET", "/marketmaker/history", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ─── Market-trend overrides ──────────────────────────────────────────────

// GetMarketTrend reads the current market-trend override. Requires ``read``.
func (m *MarketMakerAPI) GetMarketTrend(ctx context.Context, extra url.Values) (map[string]any, error) {
	out := map[string]any{}
	if err := m.client.request(ctx, "GET", "/marketmaker/get-market-trend", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SetMarketTrend sets a market-trend override. Requires ``manage``.
func (m *MarketMakerAPI) SetMarketTrend(ctx context.Context, body map[string]any) (map[string]any, error) {
	out := map[string]any{}
	if err := m.client.request(ctx, "POST", "/marketmaker/set-market-trend", nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteMarketTrend removes the current market-trend override. Requires ``manage``.
func (m *MarketMakerAPI) DeleteMarketTrend(ctx context.Context, body map[string]any) error {
	if body == nil {
		body = map[string]any{}
	}
	return m.client.request(ctx, "POST", "/marketmaker/delete-market-trend", nil, body, nil)
}

// ─── Backlog ─────────────────────────────────────────────────────────────

// Backlogs lists queued/pending market-maker backlog items. Requires ``read``.
func (m *MarketMakerAPI) Backlogs(ctx context.Context, extra url.Values) ([]map[string]any, error) {
	var out []map[string]any
	if err := m.client.request(ctx, "GET", "/marketmaker/get-backlogs", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Backlog fetches a single backlog item. Requires ``read``.
func (m *MarketMakerAPI) Backlog(ctx context.Context, id any) (map[string]any, error) {
	q := url.Values{"backlog_id": {fmt.Sprint(id)}}
	out := map[string]any{}
	if err := m.client.request(ctx, "GET", "/marketmaker/get-backlog", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteBacklog deletes a backlog item. Requires ``manage``.
func (m *MarketMakerAPI) DeleteBacklog(ctx context.Context, id any) error {
	return m.client.request(ctx, "POST", "/marketmaker/delete-backlog", nil, map[string]any{"backlog_id": id}, nil)
}
