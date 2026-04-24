package cryptohopper

import (
	"context"
	"fmt"
	"net/url"
)

// BacktestAPI is the resource namespace for backtests.
type BacktestAPI struct {
	client *Client
}

// Create starts a new backtest. Requires ``manage``. Rate bucket: backtest.
func (b *BacktestAPI) Create(ctx context.Context, body map[string]any) (map[string]any, error) {
	out := map[string]any{}
	if err := b.client.request(ctx, "POST", "/backtest/new", nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get fetches a backtest. Requires ``read``.
func (b *BacktestAPI) Get(ctx context.Context, backtestID any) (map[string]any, error) {
	q := url.Values{"backtest_id": {fmt.Sprint(backtestID)}}
	out := map[string]any{}
	if err := b.client.request(ctx, "GET", "/backtest/get", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// List returns the user's backtests. Requires ``read``.
func (b *BacktestAPI) List(ctx context.Context, extra url.Values) ([]map[string]any, error) {
	var out []map[string]any
	if err := b.client.request(ctx, "GET", "/backtest/list", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Cancel cancels a running backtest. Requires ``manage``.
func (b *BacktestAPI) Cancel(ctx context.Context, backtestID any) error {
	return b.client.request(ctx, "POST", "/backtest/cancel", nil, map[string]any{"backtest_id": backtestID}, nil)
}

// Restart restarts a backtest. Requires ``manage``.
func (b *BacktestAPI) Restart(ctx context.Context, backtestID any) (map[string]any, error) {
	out := map[string]any{}
	if err := b.client.request(ctx, "POST", "/backtest/restart", nil, map[string]any{"backtest_id": backtestID}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Limits returns the current backtest quota. Requires ``read``.
func (b *BacktestAPI) Limits(ctx context.Context) (map[string]any, error) {
	out := map[string]any{}
	if err := b.client.request(ctx, "GET", "/backtest/limits", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
