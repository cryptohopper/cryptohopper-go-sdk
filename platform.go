package cryptohopper

import (
	"context"
	"net/url"
)

// PlatformAPI is the resource namespace for public marketing / i18n /
// discovery reads. All endpoints are whitelisted server-side; no auth needed.
type PlatformAPI struct {
	client *Client
}

// LatestBlog returns the latest blog posts. Public.
func (p *PlatformAPI) LatestBlog(ctx context.Context, extra url.Values) ([]map[string]any, error) {
	var out []map[string]any
	if err := p.client.request(ctx, "GET", "/platform/latestblog", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Documentation returns documentation articles. Public.
func (p *PlatformAPI) Documentation(ctx context.Context, extra url.Values) (map[string]any, error) {
	out := map[string]any{}
	if err := p.client.request(ctx, "GET", "/platform/documentation", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// PromoBar returns the active promo bar content. Public.
func (p *PlatformAPI) PromoBar(ctx context.Context) (map[string]any, error) {
	out := map[string]any{}
	if err := p.client.request(ctx, "GET", "/platform/promobar", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SearchDocumentation does full-text search across public documentation.
func (p *PlatformAPI) SearchDocumentation(ctx context.Context, query string) ([]map[string]any, error) {
	q := url.Values{"q": {query}}
	var out []map[string]any
	if err := p.client.request(ctx, "GET", "/platform/searchdocumentation", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Countries returns the full list of countries. Public.
func (p *PlatformAPI) Countries(ctx context.Context) ([]map[string]any, error) {
	var out []map[string]any
	if err := p.client.request(ctx, "GET", "/platform/countries", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CountryAllowlist returns countries the platform currently allows. Public.
func (p *PlatformAPI) CountryAllowlist(ctx context.Context) ([]map[string]any, error) {
	var out []map[string]any
	if err := p.client.request(ctx, "GET", "/platform/countryallowlist", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// IPCountry returns the country resolved from the caller's IP. Public.
func (p *PlatformAPI) IPCountry(ctx context.Context) (map[string]any, error) {
	out := map[string]any{}
	if err := p.client.request(ctx, "GET", "/platform/ipcountry", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Languages returns supported UI languages. Public.
func (p *PlatformAPI) Languages(ctx context.Context) ([]map[string]any, error) {
	var out []map[string]any
	if err := p.client.request(ctx, "GET", "/platform/languages", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BotTypes returns the enumeration of available bot types. Public.
func (p *PlatformAPI) BotTypes(ctx context.Context) ([]map[string]any, error) {
	var out []map[string]any
	if err := p.client.request(ctx, "GET", "/platform/bottypes", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
