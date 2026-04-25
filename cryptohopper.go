// Package cryptohopper is the official Go SDK for the Cryptohopper API
// (https://www.cryptohopper.com). It wraps the same /v1/* surface the
// cryptohopper CLI and the Node/Python SDKs consume.
//
// # Quickstart
//
//	client, err := cryptohopper.NewClient(os.Getenv("CRYPTOHOPPER_TOKEN"))
//	if err != nil { log.Fatal(err) }
//
//	me, err := client.User.Get(ctx)
//	if err != nil { log.Fatal(err) }
//	fmt.Println(me["email"])
//
//	ticker, err := client.Exchange.Ticker(ctx, "binance", "BTC/USDT")
package cryptohopper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Version is the SDK's semver tag (kept in sync with the latest `v*` git tag).
const Version = "0.4.0-alpha.1"

const (
	defaultBaseURL    = "https://api.cryptohopper.com/v1"
	defaultTimeout    = 30 * time.Second
	defaultMaxRetries = 3
)

// Client is the entry point. Construct once with NewClient and reuse it across
// goroutines; it is safe for concurrent use.
type Client struct {
	apiKey     string
	appKey     string
	baseURL    string
	httpClient *http.Client
	userAgent  string
	maxRetries int

	// Resources.
	User         *UserAPI
	Hoppers      *HoppersAPI
	Exchange     *ExchangeAPI
	Strategy     *StrategyAPI
	Backtest     *BacktestAPI
	Market       *MarketAPI
	Signals      *SignalsAPI
	Arbitrage    *ArbitrageAPI
	MarketMaker  *MarketMakerAPI
	Template     *TemplateAPI
	AI           *AIAPI
	Platform     *PlatformAPI
	Chart        *ChartAPI
	Subscription *SubscriptionAPI
	Social       *SocialAPI
	Tournaments  *TournamentsAPI
	Webhooks     *WebhooksAPI
	App          *AppAPI
}

// ClientOption configures the Client at construction time.
type ClientOption func(*Client)

// WithBaseURL points the client at a staging or local dev server.
// Defaults to https://api.cryptohopper.com/v1.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) { c.baseURL = strings.TrimRight(baseURL, "/") }
}

// WithHTTPClient installs a caller-supplied *http.Client.
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) { c.httpClient = hc }
}

// WithTimeout sets the per-request timeout. Ignored when WithHTTPClient
// is also supplied — configure the timeout on your own client then.
func WithTimeout(d time.Duration) ClientOption {
	return func(c *Client) {
		if c.httpClient == nil {
			c.httpClient = &http.Client{Timeout: d}
		}
	}
}

// WithUserAgent appends a suffix to the default "cryptohopper-go-sdk/<v>"
// User-Agent header.
func WithUserAgent(suffix string) ClientOption {
	return func(c *Client) {
		if suffix != "" {
			c.userAgent = fmt.Sprintf("cryptohopper-go-sdk/%s %s", Version, suffix)
		}
	}
}

// WithAppKey sets the x-api-app-key header to the given OAuth client_id.
// Optional — most users only need the bearer token.
func WithAppKey(appKey string) ClientOption {
	return func(c *Client) { c.appKey = appKey }
}

// WithMaxRetries sets how many times to retry on HTTP 429 (respecting
// Retry-After). Pass 0 to disable. Defaults to 3.
func WithMaxRetries(n int) ClientOption {
	return func(c *Client) { c.maxRetries = n }
}

// NewClient constructs a Client with the given OAuth2 bearer token.
func NewClient(apiKey string, opts ...ClientOption) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("cryptohopper: apiKey is required")
	}

	c := &Client{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		userAgent:  fmt.Sprintf("cryptohopper-go-sdk/%s", Version),
		maxRetries: defaultMaxRetries,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.httpClient == nil {
		c.httpClient = &http.Client{Timeout: defaultTimeout}
	}

	c.User = &UserAPI{client: c}
	c.Hoppers = &HoppersAPI{client: c}
	c.Exchange = &ExchangeAPI{client: c}
	c.Strategy = &StrategyAPI{client: c}
	c.Backtest = &BacktestAPI{client: c}
	c.Market = &MarketAPI{client: c}
	c.Signals = &SignalsAPI{client: c}
	c.Arbitrage = &ArbitrageAPI{client: c}
	c.MarketMaker = &MarketMakerAPI{client: c}
	c.Template = &TemplateAPI{client: c}
	c.AI = &AIAPI{client: c}
	c.Platform = &PlatformAPI{client: c}
	c.Chart = &ChartAPI{client: c}
	c.Subscription = &SubscriptionAPI{client: c}
	c.Social = &SocialAPI{client: c}
	c.Tournaments = &TournamentsAPI{client: c}
	c.Webhooks = &WebhooksAPI{client: c}
	c.App = &AppAPI{client: c}

	return c, nil
}

// Error is returned for every non-2xx response and for transport-level
// failures.
//
//	var ce *cryptohopper.Error
//	if errors.As(err, &ce) && ce.Code == "RATE_LIMITED" {
//	    time.Sleep(ce.RetryAfter)
//	}
type Error struct {
	// Code is derived from the HTTP status (e.g. "UNAUTHORIZED",
	// "RATE_LIMITED") or "NETWORK_ERROR" / "TIMEOUT" for transport issues.
	Code string

	// Status is the HTTP status code. 0 for transport-level failures.
	Status int

	// Message is the server-provided human-readable error message (or a
	// synthetic one for transport failures).
	Message string

	// ServerCode is the numeric `code` field from the Cryptohopper error
	// envelope. Zero when absent. Used to identify rate-limit buckets.
	ServerCode int

	// IPAddress is the client IP the server saw. Useful for debugging
	// OAuth IP-whitelist mismatches. Empty when absent.
	IPAddress string

	// RetryAfter is parsed from the Retry-After header on 429 responses.
	RetryAfter time.Duration
}

// Error implements the error interface.
func (e *Error) Error() string {
	var b strings.Builder
	b.WriteString("cryptohopper: ")
	if e.Code != "" {
		b.WriteString("[")
		b.WriteString(e.Code)
		if e.Status != 0 {
			b.WriteString(" ")
			b.WriteString(strconv.Itoa(e.Status))
		}
		b.WriteString("] ")
	}
	b.WriteString(e.Message)
	if e.IPAddress != "" {
		b.WriteString(" (ip ")
		b.WriteString(e.IPAddress)
		b.WriteString(")")
	}
	return b.String()
}

// cryptohopperErrorEnvelope is the server's flat error body.
type cryptohopperErrorEnvelope struct {
	Status    int    `json:"status"`
	Code      int    `json:"code"`
	Error     int    `json:"error"`
	Message   string `json:"message"`
	IPAddress string `json:"ip_address"`
}

// request performs an HTTP request, auto-retries on 429 up to maxRetries,
// and decodes the `data` envelope into out. Pass nil for out to discard
// the body.
func (c *Client) request(ctx context.Context, method, path string, query url.Values, body, out any) error {
	attempt := 0
	for {
		err := c.doRequest(ctx, method, path, query, body, out)
		var ce *Error
		if errors.As(err, &ce) && ce.Code == "RATE_LIMITED" && attempt < c.maxRetries {
			wait := ce.RetryAfter
			if wait <= 0 {
				wait = time.Duration(1<<attempt) * time.Second
			}
			// time.NewTimer + Stop, not time.After: if ctx fires first,
			// time.After leaves a timer running until `wait` elapses,
			// holding memory + a runtime ticker entry. Stop releases it
			// immediately.
			timer := time.NewTimer(wait)
			select {
			case <-ctx.Done():
				timer.Stop()
				code := "NETWORK_ERROR"
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					code = "TIMEOUT"
				}
				return &Error{Code: code, Message: "context " + ctx.Err().Error() + " while waiting to retry"}
			case <-timer.C:
			}
			attempt++
			continue
		}
		return err
	}
}

func (c *Client) doRequest(ctx context.Context, method, path string, query url.Values, body, out any) error {
	var bodyReader io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return &Error{Code: "UNKNOWN", Message: "failed to marshal request body: " + err.Error()}
		}
		bodyReader = bytes.NewReader(buf)
	}

	full := c.baseURL + path
	if q := query.Encode(); q != "" {
		if strings.Contains(full, "?") {
			full += "&" + q
		} else {
			full += "?" + q
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, full, bodyReader)
	if err != nil {
		return &Error{Code: "UNKNOWN", Message: "failed to build request: " + err.Error()}
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if c.appKey != "" {
		req.Header.Set("x-api-app-key", c.appKey)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		code := "NETWORK_ERROR"
		msg := fmt.Sprintf("could not reach %s (%s)", c.baseURL, err.Error())
		if errors.Is(err, context.DeadlineExceeded) {
			code = "TIMEOUT"
			msg = "request timed out"
		} else if errors.Is(err, context.Canceled) {
			msg = "request cancelled"
		}
		return &Error{Code: code, Message: msg}
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Error{
			Code:    "NETWORK_ERROR",
			Status:  resp.StatusCode,
			Message: "failed to read response body: " + err.Error(),
		}
	}

	if resp.StatusCode >= 400 {
		e := &Error{
			Code:       defaultCodeForStatus(resp.StatusCode),
			Status:     resp.StatusCode,
			Message:    fmt.Sprintf("request failed (%d)", resp.StatusCode),
			RetryAfter: parseRetryAfter(resp.Header.Get("Retry-After")),
		}
		if len(raw) > 0 {
			var env cryptohopperErrorEnvelope
			if jsonErr := json.Unmarshal(raw, &env); jsonErr == nil {
				if env.Message != "" {
					e.Message = env.Message
				}
				if env.Code > 0 {
					e.ServerCode = env.Code
				}
				if env.IPAddress != "" {
					e.IPAddress = env.IPAddress
				}
			}
		}
		return e
	}

	if out == nil {
		return nil
	}
	// Try to unwrap the envelope; fall back to decoding the whole body.
	var probe map[string]json.RawMessage
	if jsonErr := json.Unmarshal(raw, &probe); jsonErr == nil {
		if data, ok := probe["data"]; ok {
			if unErr := json.Unmarshal(data, out); unErr != nil {
				return &Error{Code: "UNKNOWN", Status: resp.StatusCode, Message: "failed to decode response data: " + unErr.Error()}
			}
			return nil
		}
	}
	if unErr := json.Unmarshal(raw, out); unErr != nil {
		return &Error{Code: "UNKNOWN", Status: resp.StatusCode, Message: "failed to decode response: " + unErr.Error()}
	}
	return nil
}

func defaultCodeForStatus(status int) string {
	switch {
	case status == 400:
		return "VALIDATION_ERROR"
	case status == 401:
		return "UNAUTHORIZED"
	case status == 402:
		return "DEVICE_UNAUTHORIZED"
	case status == 403:
		return "FORBIDDEN"
	case status == 404:
		return "NOT_FOUND"
	case status == 409:
		return "CONFLICT"
	case status == 422:
		return "VALIDATION_ERROR"
	case status == 429:
		return "RATE_LIMITED"
	case status == 503:
		return "SERVICE_UNAVAILABLE"
	case status >= 500:
		return "SERVER_ERROR"
	default:
		return "UNKNOWN"
	}
}

// parseRetryAfter honours both delta-seconds and HTTP-date forms per RFC 7231.
func parseRetryAfter(header string) time.Duration {
	if header == "" {
		return 0
	}
	if secs, err := strconv.ParseFloat(header, 64); err == nil && secs >= 0 {
		return time.Duration(secs * float64(time.Second))
	}
	if when, err := http.ParseTime(header); err == nil {
		d := time.Until(when)
		if d < 0 {
			return 0
		}
		return d
	}
	return 0
}

// Params is a tiny alias for callers that want to build query strings
// inline without importing net/url.
type Params = url.Values
