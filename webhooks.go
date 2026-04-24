package cryptohopper

import "context"

// WebhooksAPI is the resource namespace for developer webhook registration.
// Maps to the server's /api/webhook_* endpoints.
type WebhooksAPI struct {
	client *Client
}

// Create registers a new webhook. Body should include the URL and
// event types.
func (w *WebhooksAPI) Create(ctx context.Context, body map[string]any) (map[string]any, error) {
	out := map[string]any{}
	if err := w.client.request(ctx, "POST", "/api/webhook_create", nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete deletes a registered webhook.
func (w *WebhooksAPI) Delete(ctx context.Context, webhookID any) error {
	return w.client.request(ctx, "POST", "/api/webhook_delete", nil, map[string]any{"webhook_id": webhookID}, nil)
}
