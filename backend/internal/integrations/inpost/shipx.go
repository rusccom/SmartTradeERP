package inpost

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type organization struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Services []string `json:"services"`
}

func newShipXClient() *http.Client {
	return &http.Client{
		Timeout:       8 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
}

func verifyOrganization(ctx context.Context, client *http.Client, a account, token string) (organization, error) {
	url := a.shipxHost() + "/v1/organizations/" + a.organizationID()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return organization{}, errCarrier
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return organization{}, errCarrier
	}
	defer resp.Body.Close()
	if err := shipxError(resp.StatusCode); err != nil {
		return organization{}, err
	}
	return readOrganization(resp.Body, a.OrganizationID)
}

func shipxError(status int) error {
	switch status {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound:
		return errCredentials
	case http.StatusTooManyRequests:
		return errThrottled
	default:
		return errCarrier
	}
}

func readOrganization(body io.Reader, id int64) (organization, error) {
	var result organization
	err := json.NewDecoder(io.LimitReader(body, 256<<10)).Decode(&result)
	if err != nil || result.ID != id || strings.TrimSpace(result.Name) == "" {
		return organization{}, errCarrier
	}
	if result.Services == nil {
		result.Services = []string{}
	}
	return result, nil
}
