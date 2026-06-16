package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Client talks to the Discord API.
type Client struct {
	cfg Config
	hc  *http.Client
}

// NewClient builds a Client from cfg.
func NewClient(cfg Config) *Client {
	// Strip "Bot " prefix if already present in the token.
	cfg.Token = strings.TrimPrefix(cfg.Token, "Bot ")
	return &Client{
		cfg: cfg,
		hc:  &http.Client{Timeout: cfg.Timeout},
	}
}

// get fetches path from the Discord API. If auth is true, sends the Bot token.
// Returns the response body on HTTP 200.
func (c *Client) get(ctx context.Context, path string, auth bool) ([]byte, error) {
	url := c.cfg.BaseURL + path
	var lastErr error
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt) * 500 * time.Millisecond):
			}
		}
		body, retry, err := c.do(ctx, url, auth)
		if err == nil {
			return body, nil
		}
		lastErr = err
		if !retry {
			return nil, err
		}
	}
	return nil, fmt.Errorf("discord: GET %s: %w", path, lastErr)
}

func (c *Client) do(ctx context.Context, url string, auth bool) ([]byte, bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)
	if auth && c.cfg.Token != "" {
		req.Header.Set("Authorization", "Bot "+c.cfg.Token)
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, true, err
	}
	defer func() { _ = resp.Body.Close() }()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, true, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		reset := resp.Header.Get("X-RateLimit-Reset-After")
		secs, _ := strconv.ParseFloat(reset, 64)
		d := time.Duration(secs*float64(time.Second)) + 100*time.Millisecond
		if d < 100*time.Millisecond {
			d = 100 * time.Millisecond
		}
		time.Sleep(d)
		return nil, true, fmt.Errorf("rate limited")
	}
	if resp.StatusCode >= 500 {
		return nil, true, fmt.Errorf("http %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		var ae APIError
		if jerr := json.Unmarshal(rawBody, &ae); jerr == nil && ae.Code != 0 {
			return nil, false, &ae
		}
		return nil, false, fmt.Errorf("http %d", resp.StatusCode)
	}
	return rawBody, false, nil
}

// FetchInvite fetches an invite by code (no auth required).
func (c *Client) FetchInvite(ctx context.Context, code string) (*Invite, error) {
	body, err := c.get(ctx, "/invites/"+code+"?with_counts=true", false)
	if err != nil {
		return nil, err
	}
	var inv Invite
	if err := json.Unmarshal(body, &inv); err != nil {
		return nil, fmt.Errorf("decode invite: %w", err)
	}
	return &inv, nil
}

// FetchMe fetches the authenticated user's profile (auth required).
func (c *Client) FetchMe(ctx context.Context) (*User, error) {
	body, err := c.get(ctx, "/users/@me", true)
	if err != nil {
		return nil, err
	}
	var u User
	if err := json.Unmarshal(body, &u); err != nil {
		return nil, fmt.Errorf("decode user: %w", err)
	}
	return &u, nil
}

// FetchGuilds fetches the list of guilds the bot has joined (auth required).
func (c *Client) FetchGuilds(ctx context.Context) ([]Guild, error) {
	body, err := c.get(ctx, "/users/@me/guilds?limit=200", true)
	if err != nil {
		return nil, err
	}
	var guilds []Guild
	if err := json.Unmarshal(body, &guilds); err != nil {
		return nil, fmt.Errorf("decode guilds: %w", err)
	}
	return guilds, nil
}

// FetchMessages fetches messages from a channel (auth required).
// before/after are optional message IDs for cursor pagination; pass "" to omit.
func (c *Client) FetchMessages(ctx context.Context, channelID string, limit int, before, after string) ([]Message, error) {
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	path := fmt.Sprintf("/channels/%s/messages?limit=%d", channelID, limit)
	if before != "" {
		path += "&before=" + before
	}
	if after != "" {
		path += "&after=" + after
	}
	body, err := c.get(ctx, path, true)
	if err != nil {
		return nil, err
	}
	var msgs []Message
	if err := json.Unmarshal(body, &msgs); err != nil {
		return nil, fmt.Errorf("decode messages: %w", err)
	}
	return msgs, nil
}
