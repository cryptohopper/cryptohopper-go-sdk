package cryptohopper

import (
	"context"
	"fmt"
	"net/url"
)

// TemplateAPI is the resource namespace for bot templates (reusable hopper
// configurations).
type TemplateAPI struct {
	client *Client
}

// List returns all templates the user has access to. Requires ``read``.
func (t *TemplateAPI) List(ctx context.Context) ([]map[string]any, error) {
	var out []map[string]any
	if err := t.client.request(ctx, "GET", "/template/templates", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get fetches a template. Requires ``read``.
func (t *TemplateAPI) Get(ctx context.Context, templateID any) (map[string]any, error) {
	q := url.Values{"template_id": {fmt.Sprint(templateID)}}
	out := map[string]any{}
	if err := t.client.request(ctx, "GET", "/template/get", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Basic fetches the basic (lightweight) view of a template. Requires ``read``.
func (t *TemplateAPI) Basic(ctx context.Context, templateID any) (map[string]any, error) {
	q := url.Values{"template_id": {fmt.Sprint(templateID)}}
	out := map[string]any{}
	if err := t.client.request(ctx, "GET", "/template/basic", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Save saves a new template. Requires ``manage``.
func (t *TemplateAPI) Save(ctx context.Context, body map[string]any) (map[string]any, error) {
	out := map[string]any{}
	if err := t.client.request(ctx, "POST", "/template/save-template", nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Update updates an existing template. Requires ``manage``.
func (t *TemplateAPI) Update(ctx context.Context, templateID any, body map[string]any) (map[string]any, error) {
	payload := mergeID(body, "template_id", templateID)
	out := map[string]any{}
	if err := t.client.request(ctx, "POST", "/template/update", nil, payload, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Load applies a template to a hopper. Requires ``manage``.
func (t *TemplateAPI) Load(ctx context.Context, templateID, hopperID any) error {
	body := map[string]any{"template_id": templateID, "hopper_id": hopperID}
	return t.client.request(ctx, "POST", "/template/load", nil, body, nil)
}

// Delete deletes a template. Requires ``manage``.
func (t *TemplateAPI) Delete(ctx context.Context, templateID any) error {
	return t.client.request(ctx, "POST", "/template/delete", nil, map[string]any{"template_id": templateID}, nil)
}
