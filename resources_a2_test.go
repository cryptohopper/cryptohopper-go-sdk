package cryptohopper

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func TestAI_AvailableModels(t *testing.T) {
	var seenPath string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		_, _ = w.Write([]byte(`{"data":[]}`))
	})
	if _, err := c.AI.AvailableModels(context.Background()); err != nil {
		t.Fatal(err)
	}
	if seenPath != "/ai/availablemodels" {
		t.Errorf("path: %q", seenPath)
	}
}

func TestAI_GetCreditsKeepsServerPrefix(t *testing.T) {
	var seenPath string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		_, _ = w.Write([]byte(`{"data":{"balance":100}}`))
	})
	out, err := c.AI.GetCredits(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if seenPath != "/ai/getaicredits" {
		t.Errorf("path: %q", seenPath)
	}
	if out["balance"].(float64) != 100 {
		t.Errorf("balance: %v", out["balance"])
	}
}

func TestAI_LLMAnalyze(t *testing.T) {
	var seenPath, seenMethod string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		seenMethod = r.Method
		_, _ = w.Write([]byte(`{"data":{"job_id":7}}`))
	})
	if _, err := c.AI.LLMAnalyze(context.Background(), map[string]any{"strategy_id": 42}); err != nil {
		t.Fatal(err)
	}
	if seenMethod != "POST" {
		t.Errorf("method: %q", seenMethod)
	}
	if seenPath != "/ai/doaillmanalyze" {
		t.Errorf("path: %q", seenPath)
	}
}

func TestPlatform_SearchDocumentation(t *testing.T) {
	var seenQuery url.Values
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenQuery = r.URL.Query()
		_, _ = w.Write([]byte(`{"data":[]}`))
	})
	if _, err := c.Platform.SearchDocumentation(context.Background(), "rsi"); err != nil {
		t.Fatal(err)
	}
	if seenQuery.Get("q") != "rsi" {
		t.Errorf("q: %q", seenQuery.Get("q"))
	}
}

func TestPlatform_BotTypes(t *testing.T) {
	var seenPath string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		_, _ = w.Write([]byte(`{"data":[]}`))
	})
	if _, err := c.Platform.BotTypes(context.Background()); err != nil {
		t.Fatal(err)
	}
	if seenPath != "/platform/bottypes" {
		t.Errorf("path: %q", seenPath)
	}
}

func TestChart_ShareSaveHyphenatedPath(t *testing.T) {
	var seenPath, seenMethod string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		seenMethod = r.Method
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	if _, err := c.Chart.ShareSave(context.Background(), map[string]any{"title": "BTC"}); err != nil {
		t.Fatal(err)
	}
	if seenMethod != "POST" {
		t.Errorf("method: %q", seenMethod)
	}
	if seenPath != "/chart/share-save" {
		t.Errorf("path: %q", seenPath)
	}
}

func TestChart_DeleteSendsChartID(t *testing.T) {
	var seenBody []byte
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenBody, _ = io.ReadAll(r.Body)
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	if err := c.Chart.Delete(context.Background(), 5); err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	_ = json.Unmarshal(seenBody, &got)
	if got["chart_id"].(float64) != 5 {
		t.Errorf("chart_id: %v", got["chart_id"])
	}
}

func TestSubscription_Plans(t *testing.T) {
	var seenPath string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		_, _ = w.Write([]byte(`{"data":[]}`))
	})
	if _, err := c.Subscription.Plans(context.Background()); err != nil {
		t.Fatal(err)
	}
	if seenPath != "/subscription/plans" {
		t.Errorf("path: %q", seenPath)
	}
}

func TestSubscription_HopperSendsID(t *testing.T) {
	var seenQuery url.Values
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenQuery = r.URL.Query()
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	if _, err := c.Subscription.Hopper(context.Background(), 42); err != nil {
		t.Fatal(err)
	}
	if seenQuery.Get("hopper_id") != "42" {
		t.Errorf("hopper_id: %q", seenQuery.Get("hopper_id"))
	}
}

func TestSubscription_StopPostsEmpty(t *testing.T) {
	var seenPath, seenMethod string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		seenMethod = r.Method
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	if err := c.Subscription.StopSubscription(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if seenMethod != "POST" {
		t.Errorf("method: %q", seenMethod)
	}
	if seenPath != "/subscription/stopsubscription" {
		t.Errorf("path: %q", seenPath)
	}
}
