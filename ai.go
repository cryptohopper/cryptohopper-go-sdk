package cryptohopper

import (
	"context"
	"fmt"
	"net/url"
)

// AIAPI is the resource namespace for AI assistant features.
type AIAPI struct {
	client *Client
}

// List returns AI assistant items / sessions. Requires ``read``.
func (a *AIAPI) List(ctx context.Context, extra url.Values) ([]map[string]any, error) {
	var out []map[string]any
	if err := a.client.request(ctx, "GET", "/ai/list", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get fetches a single AI item. Requires ``read``.
func (a *AIAPI) Get(ctx context.Context, id any) (map[string]any, error) {
	q := url.Values{"id": {fmt.Sprint(id)}}
	out := map[string]any{}
	if err := a.client.request(ctx, "GET", "/ai/get", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AvailableModels returns LLM models available to the user.
func (a *AIAPI) AvailableModels(ctx context.Context) ([]map[string]any, error) {
	var out []map[string]any
	if err := a.client.request(ctx, "GET", "/ai/availablemodels", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ─── Credits ─────────────────────────────────────────────────────────────

// GetCredits returns the remaining AI credit balance. Requires ``read``.
func (a *AIAPI) GetCredits(ctx context.Context) (map[string]any, error) {
	out := map[string]any{}
	if err := a.client.request(ctx, "GET", "/ai/getaicredits", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreditInvoices returns past invoices for AI-credit purchases. Requires ``read``.
func (a *AIAPI) CreditInvoices(ctx context.Context, extra url.Values) ([]map[string]any, error) {
	var out []map[string]any
	if err := a.client.request(ctx, "GET", "/ai/aicreditinvoices", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreditTransactions returns credit spend/top-up transaction history.
func (a *AIAPI) CreditTransactions(ctx context.Context, extra url.Values) ([]map[string]any, error) {
	var out []map[string]any
	if err := a.client.request(ctx, "GET", "/ai/aicredittransactions", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BuyCredits starts a purchase of additional credits. Requires ``user``.
func (a *AIAPI) BuyCredits(ctx context.Context, body map[string]any) (map[string]any, error) {
	out := map[string]any{}
	if err := a.client.request(ctx, "POST", "/ai/buyaicredits", nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ─── LLM analysis ────────────────────────────────────────────────────────

// LLMAnalyzeOptions returns options/metadata for the LLM analyse endpoint.
func (a *AIAPI) LLMAnalyzeOptions(ctx context.Context) (map[string]any, error) {
	out := map[string]any{}
	if err := a.client.request(ctx, "GET", "/ai/aillmanalyzeoptions", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// LLMAnalyze runs an LLM analysis. Usually async — returns a job id. Requires ``manage``.
func (a *AIAPI) LLMAnalyze(ctx context.Context, body map[string]any) (map[string]any, error) {
	out := map[string]any{}
	if err := a.client.request(ctx, "POST", "/ai/doaillmanalyze", nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// LLMAnalyzeResults fetches the result(s) of an LLM analysis. Requires ``read``.
func (a *AIAPI) LLMAnalyzeResults(ctx context.Context, extra url.Values) (map[string]any, error) {
	out := map[string]any{}
	if err := a.client.request(ctx, "GET", "/ai/aillmanalyzeresults", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// LLMResults returns historical LLM analysis results. Requires ``read``.
func (a *AIAPI) LLMResults(ctx context.Context, extra url.Values) ([]map[string]any, error) {
	var out []map[string]any
	if err := a.client.request(ctx, "GET", "/ai/aillmresults", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
