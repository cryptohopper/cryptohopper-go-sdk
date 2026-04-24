package cryptohopper

import (
	"context"
	"fmt"
	"net/url"
)

// StrategyAPI is the resource namespace for user strategies.
type StrategyAPI struct {
	client *Client
}

// List returns every strategy the user owns. Requires ``read``.
func (s *StrategyAPI) List(ctx context.Context) ([]map[string]any, error) {
	var out []map[string]any
	if err := s.client.request(ctx, "GET", "/strategy/strategies", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get fetches a strategy. Requires ``read``.
func (s *StrategyAPI) Get(ctx context.Context, strategyID any) (map[string]any, error) {
	q := url.Values{"strategy_id": {fmt.Sprint(strategyID)}}
	out := map[string]any{}
	if err := s.client.request(ctx, "GET", "/strategy/get", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Create creates a new strategy. Requires ``manage``.
func (s *StrategyAPI) Create(ctx context.Context, body map[string]any) (map[string]any, error) {
	out := map[string]any{}
	if err := s.client.request(ctx, "POST", "/strategy/create", nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Update edits an existing strategy. Requires ``manage``.
func (s *StrategyAPI) Update(ctx context.Context, strategyID any, body map[string]any) (map[string]any, error) {
	payload := mergeID(body, "strategy_id", strategyID)
	out := map[string]any{}
	if err := s.client.request(ctx, "POST", "/strategy/edit", nil, payload, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete removes a strategy. Requires ``manage``.
func (s *StrategyAPI) Delete(ctx context.Context, strategyID any) error {
	return s.client.request(ctx, "POST", "/strategy/delete", nil, map[string]any{"strategy_id": strategyID}, nil)
}
