package cryptohopper

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// newTestClient wires a Client at the httptest server's URL.
func newTestClient(t *testing.T, handler http.HandlerFunc, opts ...ClientOption) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	opts = append(opts, WithBaseURL(srv.URL), WithMaxRetries(0))
	c, err := NewClient("ch_test", opts...)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c, srv
}

func TestNewClient_RequiresAPIKey(t *testing.T) {
	if _, err := NewClient(""); err == nil {
		t.Fatal("expected error for empty apiKey, got nil")
	}
}

func TestRequest_AccessTokenAndUserAgent(t *testing.T) {
	var seenAuth, seenAccessToken, seenUA, seenAccept, seenAppKey string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenAuth = r.Header.Get("Authorization")
		seenAccessToken = r.Header.Get("access-token")
		seenUA = r.Header.Get("User-Agent")
		seenAccept = r.Header.Get("Accept")
		seenAppKey = r.Header.Get("x-api-app-key")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"ok":true}}`))
	})
	var out struct {
		OK bool `json:"ok"`
	}
	if err := c.request(context.Background(), "GET", "/user/get", nil, nil, &out); err != nil {
		t.Fatalf("request: %v", err)
	}
	if seenAccessToken != "ch_test" {
		t.Errorf("access-token: got %q, want ch_test", seenAccessToken)
	}
	if seenAuth != "" {
		t.Errorf("Authorization: got %q, want empty (Cryptohopper v1 uses access-token, not Bearer)", seenAuth)
	}
	if !strings.HasPrefix(seenUA, "cryptohopper-go-sdk/") {
		t.Errorf("User-Agent: got %q", seenUA)
	}
	if seenAccept != "application/json" {
		t.Errorf("Accept: got %q", seenAccept)
	}
	if seenAppKey != "" {
		t.Errorf("x-api-app-key: got %q, want empty", seenAppKey)
	}
	if !out.OK {
		t.Errorf("did not unwrap data envelope: out=%+v", out)
	}
}

func TestRequest_AppKeySetsHeader(t *testing.T) {
	var seenAppKey string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenAppKey = r.Header.Get("x-api-app-key")
		_, _ = w.Write([]byte(`{"data":{}}`))
	}, WithAppKey("client_123"))
	if err := c.request(context.Background(), "GET", "/user/get", nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	if seenAppKey != "client_123" {
		t.Errorf("x-api-app-key: got %q, want client_123", seenAppKey)
	}
}

func TestRequest_PostBodyAndContentType(t *testing.T) {
	var seenCT string
	var seenBody []byte
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenCT = r.Header.Get("Content-Type")
		seenBody, _ = io.ReadAll(r.Body)
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	if err := c.request(context.Background(), "POST", "/x", nil, map[string]any{"foo": 1}, nil); err != nil {
		t.Fatal(err)
	}
	if seenCT != "application/json" {
		t.Errorf("Content-Type: got %q", seenCT)
	}
	if !strings.Contains(string(seenBody), `"foo":1`) {
		t.Errorf("body mismatch: got %q", string(seenBody))
	}
}

func TestRequest_QueryParams(t *testing.T) {
	var seenQS string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenQS = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	q := map[string][]string{"exchange": {"binance"}, "market": {"BTC/USDT"}}
	if err := c.request(context.Background(), "GET", "/exchange/ticker", q, nil, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(seenQS, "exchange=binance") {
		t.Errorf("expected exchange=binance in query, got %q", seenQS)
	}
	if !strings.Contains(seenQS, "market=BTC") {
		t.Errorf("expected market in query, got %q", seenQS)
	}
}

func TestRequest_CryptohopperErrorEnvelope(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"status":403,"code":0,"error":1,"message":"This action requires 'trade' permission scope.","ip_address":"203.0.113.42"}`))
	})
	err := c.request(context.Background(), "GET", "/x", nil, nil, nil)
	var ce *Error
	if !errors.As(err, &ce) {
		t.Fatalf("expected *Error, got %T (%v)", err, err)
	}
	if ce.Code != "FORBIDDEN" {
		t.Errorf("Code: got %q", ce.Code)
	}
	if ce.Status != 403 {
		t.Errorf("Status: got %d", ce.Status)
	}
	if ce.IPAddress != "203.0.113.42" {
		t.Errorf("IPAddress: got %q", ce.IPAddress)
	}
	if ce.Message != "This action requires 'trade' permission scope." {
		t.Errorf("Message: got %q", ce.Message)
	}
}

func TestRequest_ServerCodeFromEnvelope(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Header().Set("Retry-After", "0")
		_, _ = w.Write([]byte(`{"status":429,"code":2,"error":1,"message":"Rate limit reached","ip_address":"203.0.113.42"}`))
	})
	err := c.request(context.Background(), "GET", "/x", nil, nil, nil)
	var ce *Error
	if !errors.As(err, &ce) {
		t.Fatalf("expected *Error, got %v", err)
	}
	if ce.Code != "RATE_LIMITED" {
		t.Errorf("Code: got %q", ce.Code)
	}
	if ce.ServerCode != 2 {
		t.Errorf("ServerCode: got %d, want 2", ce.ServerCode)
	}
}

func TestRequest_RetryOn429ThenSuccess(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"status":429,"code":1,"error":1,"message":"slow"}`))
			return
		}
		_, _ = w.Write([]byte(`{"data":{"ok":true}}`))
	}))
	t.Cleanup(srv.Close)
	c, err := NewClient("ch_test", WithBaseURL(srv.URL), WithMaxRetries(2))
	if err != nil {
		t.Fatal(err)
	}
	var out struct {
		OK bool `json:"ok"`
	}
	if err := c.request(context.Background(), "GET", "/x", nil, nil, &out); err != nil {
		t.Fatalf("request: %v", err)
	}
	if !out.OK {
		t.Error("expected retried request to succeed")
	}
	if atomic.LoadInt32(&calls) != 2 {
		t.Errorf("calls: got %d, want 2", atomic.LoadInt32(&calls))
	}
}

func TestRequest_GivesUpAfterMaxRetries(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"status":429,"code":1,"error":1,"message":"slow"}`))
	}))
	t.Cleanup(srv.Close)
	c, err := NewClient("ch_test", WithBaseURL(srv.URL), WithMaxRetries(2))
	if err != nil {
		t.Fatal(err)
	}
	err = c.request(context.Background(), "GET", "/x", nil, nil, nil)
	var ce *Error
	if !errors.As(err, &ce) || ce.Code != "RATE_LIMITED" {
		t.Fatalf("expected RATE_LIMITED, got %v", err)
	}
	// Initial attempt + 2 retries = 3 total.
	if got := atomic.LoadInt32(&calls); got != 3 {
		t.Errorf("calls: got %d, want 3", got)
	}
}

func TestRequest_ServerErrorFallback(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("upstream crashed"))
	})
	err := c.request(context.Background(), "GET", "/x", nil, nil, nil)
	var ce *Error
	if !errors.As(err, &ce) || ce.Code != "SERVER_ERROR" {
		t.Fatalf("expected SERVER_ERROR, got %v", err)
	}
}

func TestRequest_ContextCancelIsNetworkError(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := c.request(ctx, "GET", "/x", nil, nil, nil)
	var ce *Error
	if !errors.As(err, &ce) {
		t.Fatalf("expected *Error, got %v", err)
	}
	if ce.Code != "NETWORK_ERROR" && ce.Code != "TIMEOUT" {
		t.Errorf("Code: got %q", ce.Code)
	}
}

func TestParseRetryAfter(t *testing.T) {
	cases := []struct {
		in   string
		want time.Duration
	}{
		{"", 0},
		{"garbage", 0},
		{"0", 0},
		{"3", 3 * time.Second},
		{"1.5", 1500 * time.Millisecond},
		{"-1", 0},
	}
	for _, tc := range cases {
		if got := parseRetryAfter(tc.in); got != tc.want {
			t.Errorf("parseRetryAfter(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestErrorFormatting(t *testing.T) {
	cases := []struct {
		e    *Error
		want string
	}{
		{
			&Error{Code: "RATE_LIMITED", Status: 429, Message: "slow"},
			"cryptohopper: [RATE_LIMITED 429] slow",
		},
		{
			&Error{Code: "FORBIDDEN", Status: 403, Message: "no access", IPAddress: "1.2.3.4"},
			"cryptohopper: [FORBIDDEN 403] no access (ip 1.2.3.4)",
		},
		{
			&Error{Code: "NETWORK_ERROR", Status: 0, Message: "boom"},
			"cryptohopper: [NETWORK_ERROR] boom",
		},
	}
	for _, tc := range cases {
		if got := tc.e.Error(); got != tc.want {
			t.Errorf("Error() = %q, want %q", got, tc.want)
		}
	}
}

func TestOptions_WithBaseURLStripsTrailingSlash(t *testing.T) {
	c, err := NewClient("ch_test", WithBaseURL("https://api-staging.cryptohopper.com/v1/"))
	if err != nil {
		t.Fatal(err)
	}
	if c.baseURL != "https://api-staging.cryptohopper.com/v1" {
		t.Errorf("baseURL: got %q", c.baseURL)
	}
}

func TestOptions_WithUserAgentAppendsSuffix(t *testing.T) {
	c, err := NewClient("ch_test", WithUserAgent("myapp/1.2"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(c.userAgent, " myapp/1.2") {
		t.Errorf("userAgent: got %q", c.userAgent)
	}
}

// Guard against regressions in JSON envelope detection when the body is
// not a JSON object (e.g. plain array response).
func TestRequest_BareArrayResponse(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"id":1},{"id":2}]`))
	})
	var out []map[string]any
	if err := c.request(context.Background(), "GET", "/x", nil, nil, &out); err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Errorf("len: %d", len(out))
	}
}

// Compile-time assertion that Error satisfies the error interface.
var _ error = (*Error)(nil)

// Helper used only for compile-checking the resource field types.
var _ = func() bool {
	var c Client
	var _ *UserAPI = c.User
	var _ *HoppersAPI = c.Hoppers
	var _ *ExchangeAPI = c.Exchange
	var _ *StrategyAPI = c.Strategy
	var _ *BacktestAPI = c.Backtest
	var _ *MarketAPI = c.Market
	// silence unused import warning when json is otherwise unused
	_ = json.Marshal
	return true
}()
