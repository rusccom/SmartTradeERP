package connection

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Application struct {
	Client Client
	ID     string
	Secret string
}

func (a Application) Refresh(ctx context.Context, token string) (Credential, error) {
	form := url.Values{"grant_type": {"refresh_token"}, "refresh_token": {token}}
	return a.Grant(ctx, form, true)
}

func (a Application) Send(ctx context.Context, endpoint string, form url.Values, basic bool, dest any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://allegro.pl/auth/oauth/"+endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return Remote
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	if basic {
		req.SetBasicAuth(a.ID, a.Secret)
	}
	return a.Client.Request(req, dest)
}

func (a Application) Grant(ctx context.Context, form url.Values, basic bool) (Credential, error) {
	var result struct {
		Tokens
		Seconds int64 `json:"expires_in"`
	}
	if err := a.Send(ctx, "token", form, basic, &result); err != nil {
		return Credential{}, err
	}
	if result.Access == "" || result.Refresh == "" || result.Seconds <= 0 || result.Seconds > 365*24*3600 {
		return Credential{}, Remote
	}
	expires := time.Now().Add(time.Duration(result.Seconds) * time.Second)
	return Credential{Tokens: result.Tokens, ExpiresAt: &expires}, nil
}
