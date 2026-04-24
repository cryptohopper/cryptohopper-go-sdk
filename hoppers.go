package cryptohopper

import (
	"context"
	"fmt"
	"net/url"
)

// HoppersAPI is the resource namespace for user trading bots.
type HoppersAPI struct {
	client *Client
}

// HoppersListOptions filters the hoppers returned by List.
type HoppersListOptions struct {
	Exchange string
}

// BuySellInput is the body for Buy / Sell.
type BuySellInput struct {
	HopperID any    `json:"hopper_id"`
	Market   string `json:"market"`
	Amount   any    `json:"amount,omitempty"`
	Price    any    `json:"price,omitempty"`
	// Extra carries any additional fields the endpoint accepts.
	Extra map[string]any `json:"-"`
}

// List returns the authenticated user's hoppers. Requires ``read``.
func (h *HoppersAPI) List(ctx context.Context, opts *HoppersListOptions) ([]map[string]any, error) {
	q := url.Values{}
	if opts != nil && opts.Exchange != "" {
		q.Set("exchange", opts.Exchange)
	}
	var out []map[string]any
	if err := h.client.request(ctx, "GET", "/hopper/list", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get fetches a single hopper. Requires ``read``.
func (h *HoppersAPI) Get(ctx context.Context, hopperID any) (map[string]any, error) {
	q := url.Values{"hopper_id": {fmt.Sprint(hopperID)}}
	out := map[string]any{}
	if err := h.client.request(ctx, "GET", "/hopper/get", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Create creates a new hopper. Requires ``manage``.
func (h *HoppersAPI) Create(ctx context.Context, body map[string]any) (map[string]any, error) {
	out := map[string]any{}
	if err := h.client.request(ctx, "POST", "/hopper/create", nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Update updates a hopper. Requires ``manage``.
func (h *HoppersAPI) Update(ctx context.Context, hopperID any, body map[string]any) (map[string]any, error) {
	payload := mergeID(body, "hopper_id", hopperID)
	out := map[string]any{}
	if err := h.client.request(ctx, "POST", "/hopper/update", nil, payload, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete deletes a hopper. Requires ``manage``.
func (h *HoppersAPI) Delete(ctx context.Context, hopperID any) error {
	return h.client.request(ctx, "POST", "/hopper/delete", nil, map[string]any{"hopper_id": hopperID}, nil)
}

// Positions lists open positions for a hopper. Requires ``read``.
func (h *HoppersAPI) Positions(ctx context.Context, hopperID any) ([]map[string]any, error) {
	q := url.Values{"hopper_id": {fmt.Sprint(hopperID)}}
	var out []map[string]any
	if err := h.client.request(ctx, "GET", "/hopper/positions", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Position fetches a single position. Requires ``read``.
func (h *HoppersAPI) Position(ctx context.Context, hopperID, positionID any) (map[string]any, error) {
	q := url.Values{
		"hopper_id":   {fmt.Sprint(hopperID)},
		"position_id": {fmt.Sprint(positionID)},
	}
	out := map[string]any{}
	if err := h.client.request(ctx, "GET", "/hopper/position", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Orders lists recent orders for a hopper. Additional filters can be passed
// as query-param keys in extra. Requires ``read``.
func (h *HoppersAPI) Orders(ctx context.Context, hopperID any, extra url.Values) ([]map[string]any, error) {
	q := url.Values{"hopper_id": {fmt.Sprint(hopperID)}}
	for k, v := range extra {
		q[k] = v
	}
	var out []map[string]any
	if err := h.client.request(ctx, "GET", "/hopper/orders", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Buy places a buy. Requires ``trade``. Subject to the `order` rate bucket.
func (h *HoppersAPI) Buy(ctx context.Context, input BuySellInput) (map[string]any, error) {
	return h.sendOrder(ctx, "/hopper/buy", input)
}

// Sell places a sell. Requires ``trade``. Subject to the `order` rate bucket.
func (h *HoppersAPI) Sell(ctx context.Context, input BuySellInput) (map[string]any, error) {
	return h.sendOrder(ctx, "/hopper/sell", input)
}

func (h *HoppersAPI) sendOrder(ctx context.Context, path string, input BuySellInput) (map[string]any, error) {
	body := map[string]any{
		"hopper_id": input.HopperID,
		"market":    input.Market,
	}
	if input.Amount != nil {
		body["amount"] = input.Amount
	}
	if input.Price != nil {
		body["price"] = input.Price
	}
	for k, v := range input.Extra {
		body[k] = v
	}
	out := map[string]any{}
	if err := h.client.request(ctx, "POST", path, nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ConfigGet returns a hopper's full config. Requires ``manage``.
func (h *HoppersAPI) ConfigGet(ctx context.Context, hopperID any) (map[string]any, error) {
	q := url.Values{"hopper_id": {fmt.Sprint(hopperID)}}
	out := map[string]any{}
	if err := h.client.request(ctx, "GET", "/hopper/configget", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ConfigUpdate updates a hopper's config. Requires ``manage``.
func (h *HoppersAPI) ConfigUpdate(ctx context.Context, hopperID any, config map[string]any) (map[string]any, error) {
	payload := mergeID(config, "hopper_id", hopperID)
	out := map[string]any{}
	if err := h.client.request(ctx, "POST", "/hopper/configupdate", nil, payload, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ConfigPools lists the config pools for a hopper. Requires ``manage``.
func (h *HoppersAPI) ConfigPools(ctx context.Context, hopperID any) ([]map[string]any, error) {
	q := url.Values{"hopper_id": {fmt.Sprint(hopperID)}}
	var out []map[string]any
	if err := h.client.request(ctx, "GET", "/hopper/configpools", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Panic triggers a panic-sell on the hopper. Requires ``trade``.
func (h *HoppersAPI) Panic(ctx context.Context, hopperID any) error {
	return h.client.request(ctx, "POST", "/hopper/panic", nil, map[string]any{"hopper_id": hopperID}, nil)
}

// mergeID returns a new map that includes the given id under key.
func mergeID(base map[string]any, key string, id any) map[string]any {
	out := make(map[string]any, len(base)+1)
	out[key] = id
	for k, v := range base {
		out[k] = v
	}
	return out
}
