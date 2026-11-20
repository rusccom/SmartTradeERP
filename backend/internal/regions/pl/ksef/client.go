package ksef

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	base string
	http *http.Client
}

func New(environment string) (*Client, error) {
	hosts := map[string]string{"test": "https://api-test.ksef.mf.gov.pl/v2", "production": "https://api.ksef.mf.gov.pl/v2"}
	if hosts[environment] == "" {
		return nil, errors.New("select a KSeF environment")
	}
	client := &http.Client{Timeout: 15 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	return &Client{hosts[environment], client}, nil
}

type call struct {
	method, path, token string
	body, out           any
}

func (c *Client) request(ctx context.Context, x call) error {
	req, err := c.build(ctx, x)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return errors.New("KSeF connection failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("KSeF returned HTTP %d", resp.StatusCode)
	}
	if x.out == nil {
		return nil
	}
	return json.NewDecoder(io.LimitReader(resp.Body, 4<<20)).Decode(x.out)
}

func (c *Client) build(ctx context.Context, x call) (*http.Request, error) {
	var body io.Reader
	if x.body != nil {
		data, err := json.Marshal(x.body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, x.method, c.base+x.path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if x.token != "" {
		req.Header.Set("Authorization", "Bearer "+x.token)
	}
	return req, nil
}

type Status struct {
	Code        int      `json:"code"`
	Description string   `json:"description"`
	Details     []string `json:"details"`
}

func (s Status) Message() string {
	msg := strings.TrimSpace(strings.Join(append([]string{s.Description}, s.Details...), " "))
	if len([]rune(msg)) > 600 {
		return string([]rune(msg)[:600])
	}
	return msg
}
