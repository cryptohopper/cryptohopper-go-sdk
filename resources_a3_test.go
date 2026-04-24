package cryptohopper

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestSocial_GetProfile(t *testing.T) {
	var seenURL string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenURL = r.URL.String()
		_, _ = w.Write([]byte(`{"data":{"alias":"pim"}}`))
	})
	if _, err := c.Social.GetProfile(context.Background(), "pim"); err != nil {
		t.Fatal(err)
	}
	if seenURL != "/social/getprofile?alias=pim" {
		t.Errorf("url: %q", seenURL)
	}
}

func TestSocial_CreatePostMapsToBarePost(t *testing.T) {
	var seenPath, seenMethod string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		seenMethod = r.Method
		_, _ = w.Write([]byte(`{"data":{"id":1}}`))
	})
	if _, err := c.Social.CreatePost(context.Background(), map[string]any{"content": "hi"}); err != nil {
		t.Fatal(err)
	}
	if seenMethod != "POST" {
		t.Errorf("method: %q", seenMethod)
	}
	if seenPath != "/social/post" {
		t.Errorf("path: %q", seenPath)
	}
}

func TestSocial_GetConversationMapsToLoadConversation(t *testing.T) {
	var seenURL string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenURL = r.URL.String()
		_, _ = w.Write([]byte(`{"data":[]}`))
	})
	if _, err := c.Social.GetConversation(context.Background(), 42); err != nil {
		t.Fatal(err)
	}
	if seenURL != "/social/loadconversation?conversation_id=42" {
		t.Errorf("url: %q", seenURL)
	}
}

func TestSocial_LikePostsPostID(t *testing.T) {
	var seenBody []byte
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenBody, _ = io.ReadAll(r.Body)
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	if err := c.Social.Like(context.Background(), 99); err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	_ = json.Unmarshal(seenBody, &got)
	if got["post_id"].(float64) != 99 {
		t.Errorf("post_id: %v", got["post_id"])
	}
}

func TestTournaments_ListHitsGetTournaments(t *testing.T) {
	var seenPath string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		_, _ = w.Write([]byte(`{"data":[]}`))
	})
	if _, err := c.Tournaments.List(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if seenPath != "/tournaments/gettournaments" {
		t.Errorf("path: %q", seenPath)
	}
}

func TestTournaments_TournamentLeaderboard(t *testing.T) {
	var seenURL string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenURL = r.URL.String()
		_, _ = w.Write([]byte(`{"data":[]}`))
	})
	if _, err := c.Tournaments.TournamentLeaderboard(context.Background(), 7); err != nil {
		t.Fatal(err)
	}
	if seenURL != "/tournaments/leaderboard_tournament?tournament_id=7" {
		t.Errorf("url: %q", seenURL)
	}
}

func TestTournaments_JoinMergesID(t *testing.T) {
	var seenBody []byte
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenBody, _ = io.ReadAll(r.Body)
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	if err := c.Tournaments.Join(context.Background(), 5, map[string]any{"team": "alpha"}); err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	_ = json.Unmarshal(seenBody, &got)
	if got["tournament_id"].(float64) != 5 {
		t.Errorf("tournament_id: %v", got["tournament_id"])
	}
	if got["team"].(string) != "alpha" {
		t.Errorf("team: %v", got["team"])
	}
}

func TestWebhooks_CreateHitsApiPath(t *testing.T) {
	var seenPath, seenMethod string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		seenMethod = r.Method
		_, _ = w.Write([]byte(`{"data":{"id":1}}`))
	})
	if _, err := c.Webhooks.Create(context.Background(), map[string]any{"url": "https://e.com"}); err != nil {
		t.Fatal(err)
	}
	if seenMethod != "POST" {
		t.Errorf("method: %q", seenMethod)
	}
	if seenPath != "/api/webhook_create" {
		t.Errorf("path: %q", seenPath)
	}
}

func TestApp_InAppPurchaseUnderscored(t *testing.T) {
	var seenPath string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		_, _ = w.Write([]byte(`{"data":{}}`))
	})
	if _, err := c.App.InAppPurchase(context.Background(), map[string]any{"receipt": "abc"}); err != nil {
		t.Fatal(err)
	}
	if seenPath != "/app/in_app_purchase" {
		t.Errorf("path: %q", seenPath)
	}
}
