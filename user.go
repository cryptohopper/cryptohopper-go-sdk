package cryptohopper

import "context"

// UserAPI is the resource namespace for user-profile calls.
//
// The type is UserAPI (not User) because "User" is a natural field name on
// the Client (client.User) and also the expected payload name; the *API
// suffix keeps them cleanly distinct.
type UserAPI struct {
	client *Client
}

// Get fetches the authenticated user's profile. Requires ``user`` scope.
func (u *UserAPI) Get(ctx context.Context) (map[string]any, error) {
	out := map[string]any{}
	if err := u.client.request(ctx, "GET", "/user/get", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
