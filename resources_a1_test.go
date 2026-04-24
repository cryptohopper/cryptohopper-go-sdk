package cryptohopper

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func TestSignals_List(t *testing.T) {
	var path string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		_, _ = w.Write([]byte(`{"data":[]}`))
	})
	if _, err := c.Signals.List(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if path != "/signals/signals" {
		t.Errorf("path: %q", path)
	}
}

func TestSignals_ChartDataSingleWordPath(t *testing.T) {
	var path string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	if _, err := c.Signals.ChartData(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if path != "/signals/chartdata" {
		t.Errorf("path: %q, want /signals/chartdata", path)
	}
}

func TestArbitrage_ExchangeVsMarketStart(t *testing.T) {
	var seenPath string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	if _, err := c.Arbitrage.ExchangeStart(context.Background(), map[string]any{"hopper_id": 1}); err != nil {
		t.Fatal(err)
	}
	if seenPath != "/arbitrage/exchange" {
		t.Errorf("exchangeStart path: %q", seenPath)
	}
}

func TestArbitrage_MarketCancelHyphenatedPath(t *testing.T) {
	var seenPath, seenMethod string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		seenMethod = r.Method
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	if err := c.Arbitrage.MarketCancel(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if seenMethod != "POST" {
		t.Errorf("method: %q", seenMethod)
	}
	if seenPath != "/arbitrage/market-cancel" {
		t.Errorf("path: %q", seenPath)
	}
}

func TestArbitrage_DeleteBacklog(t *testing.T) {
	var seenBody []byte
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenBody, _ = io.ReadAll(r.Body)
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	if err := c.Arbitrage.DeleteBacklog(context.Background(), 7); err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	_ = json.Unmarshal(seenBody, &got)
	if got["backlog_id"].(float64) != 7 {
		t.Errorf("backlog_id: %v", got["backlog_id"])
	}
}

func TestMarketMaker_GetWithHopperID(t *testing.T) {
	var seenQuery url.Values
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenQuery = r.URL.Query()
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	if _, err := c.MarketMaker.Get(context.Background(), url.Values{"hopper_id": {"1"}}); err != nil {
		t.Fatal(err)
	}
	if seenQuery.Get("hopper_id") != "1" {
		t.Errorf("hopper_id: %q", seenQuery.Get("hopper_id"))
	}
}

func TestMarketMaker_SetMarketTrend(t *testing.T) {
	var seenPath, seenMethod string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		seenMethod = r.Method
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	if _, err := c.MarketMaker.SetMarketTrend(context.Background(), map[string]any{"hopper_id": 1, "trend": "bull"}); err != nil {
		t.Fatal(err)
	}
	if seenMethod != "POST" {
		t.Errorf("method: %q", seenMethod)
	}
	if seenPath != "/marketmaker/set-market-trend" {
		t.Errorf("path: %q", seenPath)
	}
}

func TestTemplate_ListHitsPluralEndpoint(t *testing.T) {
	var seenPath string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		_, _ = w.Write([]byte(`{"data":[]}`))
	})
	if _, err := c.Template.List(context.Background()); err != nil {
		t.Fatal(err)
	}
	if seenPath != "/template/templates" {
		t.Errorf("path: %q", seenPath)
	}
}

func TestTemplate_SaveHitsHyphenatedPath(t *testing.T) {
	var seenPath string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		_, _ = w.Write([]byte(`{"data":{"id":4}}`))
	})
	if _, err := c.Template.Save(context.Background(), map[string]any{"name": "my template"}); err != nil {
		t.Fatal(err)
	}
	if seenPath != "/template/save-template" {
		t.Errorf("path: %q", seenPath)
	}
}

func TestTemplate_LoadSendsBothIDs(t *testing.T) {
	var seenBody []byte
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenBody, _ = io.ReadAll(r.Body)
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	if err := c.Template.Load(context.Background(), 3, 5); err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	_ = json.Unmarshal(seenBody, &got)
	if got["template_id"].(float64) != 3 {
		t.Errorf("template_id: %v", got["template_id"])
	}
	if got["hopper_id"].(float64) != 5 {
		t.Errorf("hopper_id: %v", got["hopper_id"])
	}
}
