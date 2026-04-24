package cryptohopper

import (
	"context"
	"net/url"
)

// SignalsAPI is the resource namespace for signal-provider analytics. Distinct
// from MarketAPI's signal browse (which is for marketplace consumers).
type SignalsAPI struct {
	client *Client
}

// List returns the signals this provider has published. Requires ``read``.
func (s *SignalsAPI) List(ctx context.Context, extra url.Values) ([]map[string]any, error) {
	var out []map[string]any
	if err := s.client.request(ctx, "GET", "/signals/signals", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Performance returns performance stats (winrate, avg profit per signal, etc.).
func (s *SignalsAPI) Performance(ctx context.Context, extra url.Values) (map[string]any, error) {
	out := map[string]any{}
	if err := s.client.request(ctx, "GET", "/signals/performance", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Stats returns overall provider stats.
func (s *SignalsAPI) Stats(ctx context.Context) (map[string]any, error) {
	out := map[string]any{}
	if err := s.client.request(ctx, "GET", "/signals/stats", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Distribution returns the distribution of signals across exchanges / markets.
func (s *SignalsAPI) Distribution(ctx context.Context) (map[string]any, error) {
	out := map[string]any{}
	if err := s.client.request(ctx, "GET", "/signals/distribution", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ChartData returns data for charting provider performance over time.
func (s *SignalsAPI) ChartData(ctx context.Context, extra url.Values) (map[string]any, error) {
	out := map[string]any{}
	if err := s.client.request(ctx, "GET", "/signals/chartdata", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
