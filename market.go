package cryptohopper

import (
	"context"
	"fmt"
	"net/url"
)

// MarketAPI is the resource namespace for public marketplace browse.
type MarketAPI struct {
	client *Client
}

// Signals browses marketplace signals. Public — no auth required.
func (m *MarketAPI) Signals(ctx context.Context, extra url.Values) ([]map[string]any, error) {
	var out []map[string]any
	if err := m.client.request(ctx, "GET", "/market/signals", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Signal fetches a single marketplace signal. Public.
func (m *MarketAPI) Signal(ctx context.Context, signalID any) (map[string]any, error) {
	q := url.Values{"signal_id": {fmt.Sprint(signalID)}}
	out := map[string]any{}
	if err := m.client.request(ctx, "GET", "/market/signal", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Items browses marketplace items (strategies, templates, signals). Public.
func (m *MarketAPI) Items(ctx context.Context, extra url.Values) ([]map[string]any, error) {
	var out []map[string]any
	if err := m.client.request(ctx, "GET", "/market/marketitems", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Item fetches a single marketplace item. Public.
func (m *MarketAPI) Item(ctx context.Context, itemID any) (map[string]any, error) {
	q := url.Values{"item_id": {fmt.Sprint(itemID)}}
	out := map[string]any{}
	if err := m.client.request(ctx, "GET", "/market/marketitem", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Homepage returns the marketplace homepage payload. Public.
func (m *MarketAPI) Homepage(ctx context.Context) (map[string]any, error) {
	out := map[string]any{}
	if err := m.client.request(ctx, "GET", "/market/homepage", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
