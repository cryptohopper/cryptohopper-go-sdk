package cryptohopper

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestUser_Get(t *testing.T) {
	var path string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		_, _ = w.Write([]byte(`{"data":{"id":42,"email":"test@example.com","username":"pim","userHash":"abc"}}`))
	})
	u, err := c.User.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if path != "/user/get" {
		t.Errorf("path: %q", path)
	}
	if u["email"] != "test@example.com" {
		t.Errorf("email: %v", u["email"])
	}
}

func TestHoppers_ListWithExchangeFilter(t *testing.T) {
	var seenQuery url.Values
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenQuery = r.URL.Query()
		_, _ = w.Write([]byte(`{"data":[]}`))
	})
	if _, err := c.Hoppers.List(context.Background(), &HoppersListOptions{Exchange: "binance"}); err != nil {
		t.Fatal(err)
	}
	if seenQuery.Get("exchange") != "binance" {
		t.Errorf("exchange: %q", seenQuery.Get("exchange"))
	}
}

func TestHoppers_GetSendsHopperID(t *testing.T) {
	var seenQuery url.Values
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenQuery = r.URL.Query()
		_, _ = w.Write([]byte(`{"data":{"id":42}}`))
	})
	if _, err := c.Hoppers.Get(context.Background(), 42); err != nil {
		t.Fatal(err)
	}
	if seenQuery.Get("hopper_id") != "42" {
		t.Errorf("hopper_id: %q", seenQuery.Get("hopper_id"))
	}
}

func TestHoppers_BuyPostsBody(t *testing.T) {
	var seenPath, seenMethod string
	var seenBody []byte
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		seenMethod = r.Method
		seenBody, _ = io.ReadAll(r.Body)
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	_, err := c.Hoppers.Buy(context.Background(), BuySellInput{
		HopperID: 42,
		Market:   "BTC/USDT",
		Amount:   "0.001",
	})
	if err != nil {
		t.Fatal(err)
	}
	if seenMethod != "POST" {
		t.Errorf("method: %q", seenMethod)
	}
	if seenPath != "/hopper/buy" {
		t.Errorf("path: %q", seenPath)
	}
	var got map[string]any
	if err := json.Unmarshal(seenBody, &got); err != nil {
		t.Fatal(err)
	}
	if got["market"] != "BTC/USDT" {
		t.Errorf("market: %v", got["market"])
	}
}

func TestHoppers_ConfigUpdateMergesID(t *testing.T) {
	var seenBody []byte
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenBody, _ = io.ReadAll(r.Body)
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	if _, err := c.Hoppers.ConfigUpdate(context.Background(), 7, map[string]any{"strategy_id": 99}); err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	_ = json.Unmarshal(seenBody, &got)
	if got["hopper_id"].(float64) != 7 {
		t.Errorf("hopper_id: %v", got["hopper_id"])
	}
	if got["strategy_id"].(float64) != 99 {
		t.Errorf("strategy_id: %v", got["strategy_id"])
	}
}

func TestHoppers_Panic(t *testing.T) {
	var seenPath, seenMethod string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		seenMethod = r.Method
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	if err := c.Hoppers.Panic(context.Background(), 5); err != nil {
		t.Fatal(err)
	}
	if seenMethod != "POST" {
		t.Errorf("method: %q", seenMethod)
	}
	if seenPath != "/hopper/panic" {
		t.Errorf("path: %q", seenPath)
	}
}

func TestExchange_Ticker(t *testing.T) {
	var seenQuery url.Values
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenQuery = r.URL.Query()
		_, _ = w.Write([]byte(`{"data":{"last":42000}}`))
	})
	t2, err := c.Exchange.Ticker(context.Background(), "binance", "BTC/USDT")
	if err != nil {
		t.Fatal(err)
	}
	if seenQuery.Get("exchange") != "binance" || seenQuery.Get("market") != "BTC/USDT" {
		t.Errorf("query: %+v", seenQuery)
	}
	if t2["last"].(float64) != 42000 {
		t.Errorf("last: %v", t2["last"])
	}
}

func TestExchange_CandlesWithOptions(t *testing.T) {
	var seenQuery url.Values
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenQuery = r.URL.Query()
		_, _ = w.Write([]byte(`{"data":[]}`))
	})
	_, err := c.Exchange.Candles(context.Background(), "binance", "BTC/USDT", "1h", &CandleOptions{From: 1700000000, To: 1700003600})
	if err != nil {
		t.Fatal(err)
	}
	if seenQuery.Get("timeframe") != "1h" {
		t.Errorf("timeframe: %q", seenQuery.Get("timeframe"))
	}
	if seenQuery.Get("from") != "1700000000" {
		t.Errorf("from: %q", seenQuery.Get("from"))
	}
}

func TestExchange_Exchanges_NoParams(t *testing.T) {
	var seenURL string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenURL = r.URL.String()
		_, _ = w.Write([]byte(`{"data":[]}`))
	})
	if _, err := c.Exchange.Exchanges(context.Background()); err != nil {
		t.Fatal(err)
	}
	if seenURL != "/exchange/exchanges" {
		t.Errorf("url: %q", seenURL)
	}
}

func TestStrategy_ListHitsPluralEndpoint(t *testing.T) {
	var seenPath string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		_, _ = w.Write([]byte(`{"data":[]}`))
	})
	if _, err := c.Strategy.List(context.Background()); err != nil {
		t.Fatal(err)
	}
	if seenPath != "/strategy/strategies" {
		t.Errorf("path: %q", seenPath)
	}
}

func TestStrategy_UpdateHitsEdit(t *testing.T) {
	var seenPath string
	var seenBody []byte
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		seenBody, _ = io.ReadAll(r.Body)
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	if _, err := c.Strategy.Update(context.Background(), 5, map[string]any{"name": "renamed"}); err != nil {
		t.Fatal(err)
	}
	if seenPath != "/strategy/edit" {
		t.Errorf("path: %q", seenPath)
	}
	if !strings.Contains(string(seenBody), `"strategy_id":5`) {
		t.Errorf("body missing strategy_id: %s", string(seenBody))
	}
	if !strings.Contains(string(seenBody), `"name":"renamed"`) {
		t.Errorf("body missing name: %s", string(seenBody))
	}
}

func TestBacktest_CreateHitsNew(t *testing.T) {
	var seenPath, seenMethod string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		seenMethod = r.Method
		_, _ = w.Write([]byte(`{"data":{"id":1}}`))
	})
	if _, err := c.Backtest.Create(context.Background(), map[string]any{"hopper_id": 42}); err != nil {
		t.Fatal(err)
	}
	if seenPath != "/backtest/new" {
		t.Errorf("path: %q", seenPath)
	}
	if seenMethod != "POST" {
		t.Errorf("method: %q", seenMethod)
	}
}

func TestBacktest_Limits(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":{"remaining":3,"limit":5}}`))
	})
	out, err := c.Backtest.Limits(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if out["remaining"].(float64) != 3 {
		t.Errorf("remaining: %v", out["remaining"])
	}
}

func TestMarket_ItemsUsesMarketitemsEndpoint(t *testing.T) {
	var seenPath string
	var seenQuery url.Values
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		seenQuery = r.URL.Query()
		_, _ = w.Write([]byte(`{"data":[]}`))
	})
	if _, err := c.Market.Items(context.Background(), url.Values{"type": {"strategy"}}); err != nil {
		t.Fatal(err)
	}
	if seenPath != "/market/marketitems" {
		t.Errorf("path: %q", seenPath)
	}
	if seenQuery.Get("type") != "strategy" {
		t.Errorf("type: %q", seenQuery.Get("type"))
	}
}

func TestMarket_Signal(t *testing.T) {
	var seenQuery url.Values
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenQuery = r.URL.Query()
		_, _ = w.Write([]byte(`{"data":{"id":99}}`))
	})
	if _, err := c.Market.Signal(context.Background(), 99); err != nil {
		t.Fatal(err)
	}
	if seenQuery.Get("signal_id") != "99" {
		t.Errorf("signal_id: %q", seenQuery.Get("signal_id"))
	}
}
