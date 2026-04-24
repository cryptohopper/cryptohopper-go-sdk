package cryptohopper

import (
	"context"
	"fmt"
	"net/url"
)

// TournamentsAPI is the resource namespace for trading competitions.
type TournamentsAPI struct {
	client *Client
}

// List returns all tournaments. Requires ``read``.
func (t *TournamentsAPI) List(ctx context.Context, extra url.Values) ([]map[string]any, error) {
	var out []map[string]any
	if err := t.client.request(ctx, "GET", "/tournaments/gettournaments", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Active returns currently-active tournaments. Public.
func (t *TournamentsAPI) Active(ctx context.Context) ([]map[string]any, error) {
	var out []map[string]any
	if err := t.client.request(ctx, "GET", "/tournaments/active", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get fetches a single tournament. Requires ``read``.
func (t *TournamentsAPI) Get(ctx context.Context, tournamentID any) (map[string]any, error) {
	q := url.Values{"tournament_id": {fmt.Sprint(tournamentID)}}
	out := map[string]any{}
	if err := t.client.request(ctx, "GET", "/tournaments/gettournament", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Search searches across tournaments. Requires ``read``.
func (t *TournamentsAPI) Search(ctx context.Context, query string) ([]map[string]any, error) {
	q := url.Values{"q": {query}}
	var out []map[string]any
	if err := t.client.request(ctx, "GET", "/tournaments/search", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Trades returns trades in a tournament. Requires ``read``.
func (t *TournamentsAPI) Trades(ctx context.Context, tournamentID any) ([]map[string]any, error) {
	q := url.Values{"tournament_id": {fmt.Sprint(tournamentID)}}
	var out []map[string]any
	if err := t.client.request(ctx, "GET", "/tournaments/trades", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Stats returns aggregated stats for a tournament. Requires ``read``.
func (t *TournamentsAPI) Stats(ctx context.Context, tournamentID any) (map[string]any, error) {
	q := url.Values{"tournament_id": {fmt.Sprint(tournamentID)}}
	out := map[string]any{}
	if err := t.client.request(ctx, "GET", "/tournaments/stats", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Activity returns the activity feed for a tournament. Requires ``read``.
func (t *TournamentsAPI) Activity(ctx context.Context, tournamentID any) ([]map[string]any, error) {
	q := url.Values{"tournament_id": {fmt.Sprint(tournamentID)}}
	var out []map[string]any
	if err := t.client.request(ctx, "GET", "/tournaments/activity", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Leaderboard returns the cross-tournament leaderboard. Requires ``read``.
func (t *TournamentsAPI) Leaderboard(ctx context.Context, extra url.Values) ([]map[string]any, error) {
	var out []map[string]any
	if err := t.client.request(ctx, "GET", "/tournaments/leaderboard", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TournamentLeaderboard returns the leaderboard for a specific tournament.
// Requires ``read``.
func (t *TournamentsAPI) TournamentLeaderboard(ctx context.Context, tournamentID any) ([]map[string]any, error) {
	q := url.Values{"tournament_id": {fmt.Sprint(tournamentID)}}
	var out []map[string]any
	if err := t.client.request(ctx, "GET", "/tournaments/leaderboard_tournament", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Join joins a tournament. Requires ``manage``.
func (t *TournamentsAPI) Join(ctx context.Context, tournamentID any, body map[string]any) error {
	payload := mergeID(body, "tournament_id", tournamentID)
	return t.client.request(ctx, "POST", "/tournaments/join", nil, payload, nil)
}

// Leave leaves a tournament. Requires ``manage``.
func (t *TournamentsAPI) Leave(ctx context.Context, tournamentID any) error {
	return t.client.request(ctx, "POST", "/tournaments/leave", nil, map[string]any{"tournament_id": tournamentID}, nil)
}
