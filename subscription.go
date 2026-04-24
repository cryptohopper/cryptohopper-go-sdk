package cryptohopper

import (
	"context"
	"fmt"
	"net/url"
)

// SubscriptionAPI is the resource namespace for plans, per-hopper state,
// credits, and billing flows.
type SubscriptionAPI struct {
	client *Client
}

// Hopper returns subscription state for a specific hopper. Requires ``read``.
func (s *SubscriptionAPI) Hopper(ctx context.Context, hopperID any) (map[string]any, error) {
	q := url.Values{"hopper_id": {fmt.Sprint(hopperID)}}
	out := map[string]any{}
	if err := s.client.request(ctx, "GET", "/subscription/hopper", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get returns account-level subscription state. Requires ``read``.
func (s *SubscriptionAPI) Get(ctx context.Context) (map[string]any, error) {
	out := map[string]any{}
	if err := s.client.request(ctx, "GET", "/subscription/get", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Plans returns the list of available subscription plans. Public.
func (s *SubscriptionAPI) Plans(ctx context.Context) ([]map[string]any, error) {
	var out []map[string]any
	if err := s.client.request(ctx, "GET", "/subscription/plans", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Remap moves a subscription slot from one hopper to another. Requires ``manage``.
func (s *SubscriptionAPI) Remap(ctx context.Context, body map[string]any) (map[string]any, error) {
	out := map[string]any{}
	if err := s.client.request(ctx, "POST", "/subscription/remap", nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Assign assigns a subscription slot to a hopper. Requires ``manage``.
func (s *SubscriptionAPI) Assign(ctx context.Context, body map[string]any) (map[string]any, error) {
	out := map[string]any{}
	if err := s.client.request(ctx, "POST", "/subscription/assign", nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCredits returns remaining platform credits. Requires ``read``.
func (s *SubscriptionAPI) GetCredits(ctx context.Context) (map[string]any, error) {
	out := map[string]any{}
	if err := s.client.request(ctx, "GET", "/subscription/getcredits", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// OrderSub starts a subscription purchase. Requires ``user``.
func (s *SubscriptionAPI) OrderSub(ctx context.Context, body map[string]any) (map[string]any, error) {
	out := map[string]any{}
	if err := s.client.request(ctx, "POST", "/subscription/ordersub", nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// StopSubscription cancels / stops an active subscription. Requires ``user``.
func (s *SubscriptionAPI) StopSubscription(ctx context.Context, body map[string]any) error {
	if body == nil {
		body = map[string]any{}
	}
	return s.client.request(ctx, "POST", "/subscription/stopsubscription", nil, body, nil)
}
