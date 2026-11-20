package connection

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const Path = "/api/client/integrations/allegro"
const Scopes = "allegro:api:profile:read allegro:api:orders:read"
const userAgent = "SmartTradeERP/1.0 (+https://github.com/rusccom/SmartTradeERP)"

type Client struct{ http *http.Client }

func NewClient() Client {
	return Client{&http.Client{Timeout: 6 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}}
}

func (c Client) Request(req *http.Request, dest any) error {
	req.Header.Set("User-Agent", userAgent)
	res, err := c.http.Do(req)
	if err != nil {
		return Remote
	}
	defer res.Body.Close()
	reader := io.LimitReader(res.Body, 128<<10)
	if res.StatusCode != http.StatusOK {
		return providerError(res.StatusCode, reader)
	}
	if dest != nil && json.NewDecoder(reader).Decode(dest) != nil {
		return Remote
	}
	return nil
}

func providerError(status int, body io.Reader) error {
	var response struct {
		Error string `json:"error"`
	}
	_ = json.NewDecoder(body).Decode(&response)
	if problem := grantProblems[response.Error]; problem != nil {
		return problem
	}
	switch status {
	case http.StatusUnauthorized:
		return Unauthorized
	case http.StatusForbidden:
		return Scope
	case http.StatusTooManyRequests:
		return Limit
	default:
		return Remote
	}
}

var grantProblems = map[string]*Problem{
	"authorization_pending": Pending, "slow_down": Slow, "expired_token": State, "Invalid device code": State,
	"invalid_client": Setup, "invalid_grant": Unauthorized, "access_denied": Denied, "invalid_scope": Scope,
}

func (c Client) read(ctx context.Context, path, token string, dest any) error {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.allegro.pl"+path, nil)
	if err != nil {
		return Remote
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.allegro.public.v1+json")
	return c.Request(req, dest)
}

func (c Client) Verify(ctx context.Context, token string) (Seller, error) {
	var seller Seller
	if err := c.read(ctx, "/me", token, &seller); err != nil {
		return seller, err
	}
	if seller.ID == "" || seller.Login == "" {
		return seller, Remote
	}
	return seller, c.read(ctx, "/order/event-stats", token, nil)
}
