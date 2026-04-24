package cryptohopper

import "context"

// AppAPI is the resource namespace for mobile app store receipts and
// in-app purchases.
type AppAPI struct {
	client *Client
}

// Receipt validates an App Store / Play Store receipt.
func (a *AppAPI) Receipt(ctx context.Context, body map[string]any) (map[string]any, error) {
	out := map[string]any{}
	if err := a.client.request(ctx, "POST", "/app/receipt", nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// InAppPurchase records an in-app purchase.
func (a *AppAPI) InAppPurchase(ctx context.Context, body map[string]any) (map[string]any, error) {
	out := map[string]any{}
	if err := a.client.request(ctx, "POST", "/app/in_app_purchase", nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}
