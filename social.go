package cryptohopper

import (
	"context"
	"fmt"
	"net/url"
)

// SocialAPI is the resource namespace for the Cryptohopper social graph —
// profiles, feed, posts, comments, conversations, follows, engagement.
// Largest resource in the SDK (27 methods).
type SocialAPI struct {
	client *Client
}

// ─── Profiles ────────────────────────────────────────────────────────────

// GetProfile fetches a public profile by alias or id. Requires ``read``.
func (s *SocialAPI) GetProfile(ctx context.Context, aliasOrId any) (map[string]any, error) {
	q := url.Values{"alias": {fmt.Sprint(aliasOrId)}}
	out := map[string]any{}
	if err := s.client.request(ctx, "GET", "/social/getprofile", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// EditProfile updates the authenticated user's profile. Requires ``user``.
func (s *SocialAPI) EditProfile(ctx context.Context, body map[string]any) (map[string]any, error) {
	out := map[string]any{}
	if err := s.client.request(ctx, "POST", "/social/editprofile", nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CheckAlias checks whether an alias is available.
func (s *SocialAPI) CheckAlias(ctx context.Context, alias string) (map[string]any, error) {
	q := url.Values{"alias": {alias}}
	out := map[string]any{}
	if err := s.client.request(ctx, "GET", "/social/checkalias", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ─── Feed / discovery ────────────────────────────────────────────────────

// GetFeed returns the user's personalised feed. Requires ``read``.
func (s *SocialAPI) GetFeed(ctx context.Context, extra url.Values) ([]map[string]any, error) {
	var out []map[string]any
	if err := s.client.request(ctx, "GET", "/social/getfeed", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTrends returns trending topics. Requires ``read``.
func (s *SocialAPI) GetTrends(ctx context.Context) ([]map[string]any, error) {
	var out []map[string]any
	if err := s.client.request(ctx, "GET", "/social/gettrends", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// WhoToFollow returns suggested profiles to follow. Requires ``read``.
func (s *SocialAPI) WhoToFollow(ctx context.Context) ([]map[string]any, error) {
	var out []map[string]any
	if err := s.client.request(ctx, "GET", "/social/whotofollow", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Search searches for posts / users. Requires ``read``.
func (s *SocialAPI) Search(ctx context.Context, query string) ([]map[string]any, error) {
	q := url.Values{"q": {query}}
	var out []map[string]any
	if err := s.client.request(ctx, "GET", "/social/search", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ─── Notifications ───────────────────────────────────────────────────────

// GetNotifications returns notifications for the authenticated user.
// Requires ``notifications``.
func (s *SocialAPI) GetNotifications(ctx context.Context, extra url.Values) ([]map[string]any, error) {
	var out []map[string]any
	if err := s.client.request(ctx, "GET", "/social/getnotifications", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ─── Conversations / messages ────────────────────────────────────────────

// GetConversationList lists the user's DM conversations. Requires ``read``.
func (s *SocialAPI) GetConversationList(ctx context.Context) ([]map[string]any, error) {
	var out []map[string]any
	if err := s.client.request(ctx, "GET", "/social/getconversationlist", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetConversation loads messages for a single conversation. Requires ``read``.
func (s *SocialAPI) GetConversation(ctx context.Context, conversationID any) ([]map[string]any, error) {
	q := url.Values{"conversation_id": {fmt.Sprint(conversationID)}}
	var out []map[string]any
	if err := s.client.request(ctx, "GET", "/social/loadconversation", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SendMessage sends a DM. Requires ``user``.
func (s *SocialAPI) SendMessage(ctx context.Context, body map[string]any) (map[string]any, error) {
	out := map[string]any{}
	if err := s.client.request(ctx, "POST", "/social/sendmessage", nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteMessage deletes a DM. Requires ``user``.
func (s *SocialAPI) DeleteMessage(ctx context.Context, messageID any) error {
	return s.client.request(ctx, "POST", "/social/deletemessage", nil, map[string]any{"message_id": messageID}, nil)
}

// ─── Posts ───────────────────────────────────────────────────────────────

// CreatePost creates a new post. Requires ``user``.
func (s *SocialAPI) CreatePost(ctx context.Context, body map[string]any) (map[string]any, error) {
	out := map[string]any{}
	if err := s.client.request(ctx, "POST", "/social/post", nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetPost fetches a single post. Requires ``read``.
func (s *SocialAPI) GetPost(ctx context.Context, postID any) (map[string]any, error) {
	q := url.Values{"post_id": {fmt.Sprint(postID)}}
	out := map[string]any{}
	if err := s.client.request(ctx, "GET", "/social/getpost", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeletePost deletes a post. Requires ``user``.
func (s *SocialAPI) DeletePost(ctx context.Context, postID any) error {
	return s.client.request(ctx, "POST", "/social/deletepost", nil, map[string]any{"post_id": postID}, nil)
}

// PinPost pins/unpins a post. Requires ``user``.
func (s *SocialAPI) PinPost(ctx context.Context, postID any) error {
	return s.client.request(ctx, "POST", "/social/pinpost", nil, map[string]any{"post_id": postID}, nil)
}

// ─── Comments ────────────────────────────────────────────────────────────

// GetComment fetches a single comment. Requires ``read``.
func (s *SocialAPI) GetComment(ctx context.Context, commentID any) (map[string]any, error) {
	q := url.Values{"comment_id": {fmt.Sprint(commentID)}}
	out := map[string]any{}
	if err := s.client.request(ctx, "GET", "/social/getcomment", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetComments lists comments on a post. Requires ``read``.
func (s *SocialAPI) GetComments(ctx context.Context, postID any) ([]map[string]any, error) {
	q := url.Values{"post_id": {fmt.Sprint(postID)}}
	var out []map[string]any
	if err := s.client.request(ctx, "GET", "/social/getcomments", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteComment deletes a comment. Requires ``user``.
func (s *SocialAPI) DeleteComment(ctx context.Context, commentID any) error {
	return s.client.request(ctx, "POST", "/social/deletecomment", nil, map[string]any{"comment_id": commentID}, nil)
}

// ─── Media ───────────────────────────────────────────────────────────────

// GetMedia fetches a media attachment. Requires ``read``.
func (s *SocialAPI) GetMedia(ctx context.Context, mediaID any) (map[string]any, error) {
	q := url.Values{"media_id": {fmt.Sprint(mediaID)}}
	out := map[string]any{}
	if err := s.client.request(ctx, "GET", "/social/getmedia", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ─── Social graph ────────────────────────────────────────────────────────

// Follow follows/unfollows a profile. Requires ``user``.
func (s *SocialAPI) Follow(ctx context.Context, aliasOrId any) error {
	return s.client.request(ctx, "POST", "/social/follow", nil, map[string]any{"alias": aliasOrId}, nil)
}

// GetFollowers lists followers. Requires ``read``.
func (s *SocialAPI) GetFollowers(ctx context.Context, aliasOrId any) ([]map[string]any, error) {
	q := url.Values{"alias": {fmt.Sprint(aliasOrId)}}
	var out []map[string]any
	if err := s.client.request(ctx, "GET", "/social/followers", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetFollowing checks whether the auth'd user follows the given profile. Requires ``read``.
func (s *SocialAPI) GetFollowing(ctx context.Context, aliasOrId any) (map[string]any, error) {
	q := url.Values{"alias": {fmt.Sprint(aliasOrId)}}
	out := map[string]any{}
	if err := s.client.request(ctx, "GET", "/social/following", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetFollowingProfiles lists profiles the given user follows. Requires ``read``.
func (s *SocialAPI) GetFollowingProfiles(ctx context.Context, aliasOrId any) ([]map[string]any, error) {
	q := url.Values{"alias": {fmt.Sprint(aliasOrId)}}
	var out []map[string]any
	if err := s.client.request(ctx, "GET", "/social/followingprofiles", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ─── Engagement ──────────────────────────────────────────────────────────

// Like likes/unlikes a post. Requires ``user``.
func (s *SocialAPI) Like(ctx context.Context, postID any) error {
	return s.client.request(ctx, "POST", "/social/like", nil, map[string]any{"post_id": postID}, nil)
}

// Repost reposts a post. Requires ``user``.
func (s *SocialAPI) Repost(ctx context.Context, postID any) error {
	return s.client.request(ctx, "POST", "/social/repost", nil, map[string]any{"post_id": postID}, nil)
}

// ─── Moderation ──────────────────────────────────────────────────────────

// BlockUser blocks a user. Requires ``user``.
func (s *SocialAPI) BlockUser(ctx context.Context, aliasOrId any) error {
	return s.client.request(ctx, "POST", "/social/blockuser", nil, map[string]any{"alias": aliasOrId}, nil)
}
