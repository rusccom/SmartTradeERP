package ksef

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"time"
)

type accessToken struct {
	Token string `json:"token"`
}
type challenge struct {
	Challenge string `json:"challenge"`
	Timestamp int64  `json:"timestampMs"`
}
type authorization struct {
	Reference      string      `json:"referenceNumber"`
	Authentication accessToken `json:"authenticationToken"`
}
type Session struct {
	client *Client
	token  string
	certs  []certificate
}

func (c *Client) Authenticate(ctx context.Context, nip, secret string) (*Session, error) {
	var certs []certificate
	if err := c.request(ctx, call{method: "GET", path: "/security/public-key-certificates", out: &certs}); err != nil {
		return nil, err
	}
	cert, err := selectCert(certs, "KsefTokenEncryption")
	if err != nil {
		return nil, err
	}
	auth, err := c.startAuth(ctx, nip, secret, cert)
	if err != nil {
		return nil, err
	}
	if err = c.waitAuth(ctx, auth); err != nil {
		return nil, err
	}
	return c.redeem(ctx, auth.Authentication.Token, certs)
}
func (c *Client) redeem(ctx context.Context, token string, certs []certificate) (*Session, error) {
	var result struct {
		Access accessToken `json:"accessToken"`
	}
	err := c.request(ctx, call{method: "POST", path: "/auth/token/redeem", token: token, out: &result})
	if err == nil && result.Access.Token == "" {
		return nil, errors.New("KSeF did not return an access token")
	}
	return &Session{c, result.Access.Token, certs}, err
}

func (c *Client) startAuth(ctx context.Context, nip, token string, cert certificate) (authorization, error) {
	var ch challenge
	var auth authorization
	if err := c.request(ctx, call{method: "POST", path: "/auth/challenge", out: &ch}); err != nil {
		return auth, err
	}
	if ch.Timestamp <= 0 || ch.Challenge == "" {
		return auth, errors.New("invalid KSeF challenge")
	}
	encrypted, err := cert.wrap([]byte(token + "|" + strconv.FormatInt(ch.Timestamp, 10)))
	if err != nil {
		return auth, err
	}
	body := map[string]any{"challenge": ch.Challenge, "contextIdentifier": map[string]string{"type": "Nip", "value": nip}, "encryptedToken": encrypted, "publicKeyId": cert.ID}
	err = c.request(ctx, call{method: "POST", path: "/auth/ksef-token", body: body, out: &auth})
	if err == nil && (auth.Reference == "" || auth.Authentication.Token == "") {
		return auth, errors.New("KSeF denied authentication")
	}
	return auth, err
}

func (c *Client) waitAuth(ctx context.Context, a authorization) error {
	for i := 0; i < 9; i++ {
		code, err := c.authStatus(ctx, a)
		if err != nil {
			return err
		}
		if code == 200 {
			return nil
		}
		if code != 100 {
			return errors.New("KSeF authentication denied")
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
	return errors.New("KSeF authentication is still processing")
}

func (c *Client) authStatus(ctx context.Context, a authorization) (int, error) {
	var result struct {
		Status Status `json:"status"`
	}
	err := c.request(ctx, call{method: "GET", path: "/auth/" + url.PathEscape(a.Reference), token: a.Authentication.Token, out: &result})
	return result.Status.Code, err
}
