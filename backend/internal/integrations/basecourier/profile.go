package basecourier

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const profileURL = "https://api.basecourier.com/api/getProfile.json"

func newProfileClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}
}

func verifyProfile(ctx context.Context, client *http.Client, email, key string) error {
	auth := map[string]string{"login": email, "api_key": key}
	body, err := json.Marshal(map[string]any{"auth": auth})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, profileURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return errRemote
	}
	defer resp.Body.Close()
	return profileResult(resp)
}

func profileResult(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusOK:
		return readProfile(resp.Body)
	case http.StatusUnauthorized, http.StatusForbidden:
		return errDenied
	case http.StatusTooManyRequests:
		return errLimited
	default:
		return errRemote
	}
}

func readProfile(body io.Reader) error {
	var result struct {
		Success *bool                      `json:"success"`
		Data    struct{ Broker *struct{} } `json:"data"`
	}
	const limit = 256 << 10
	raw, err := io.ReadAll(io.LimitReader(body, limit+1))
	if err != nil || len(raw) > limit || json.Unmarshal(raw, &result) != nil || result.Success == nil {
		return errRemote
	}
	if !*result.Success {
		return errDenied
	}
	if result.Data.Broker == nil {
		return errRemote
	}
	return nil
}
