package discord

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("../testdata/" + name)
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return data
}

func testClient(baseURL string) *Client {
	cfg := DefaultConfig()
	cfg.Token = "testtoken"
	cfg.BaseURL = baseURL
	cfg.Retries = 0
	cfg.Timeout = 5 * time.Second
	return NewClient(cfg)
}

func TestFetchInvitePublic(t *testing.T) {
	data := loadFixture(t, "invite_discord_developers.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/invites/") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	}))
	defer srv.Close()

	c := testClient(srv.URL)
	inv, err := c.FetchInvite(context.Background(), "discord-developers")
	if err != nil {
		t.Fatalf("FetchInvite: %v", err)
	}
	if inv.ApproxMemberCount != 152000 {
		t.Errorf("member count = %d, want 152000", inv.ApproxMemberCount)
	}
	if inv.Guild == nil || inv.Guild.Name != "Discord Developers" {
		t.Errorf("guild name mismatch: %+v", inv.Guild)
	}
}

func TestFetchInviteUnknown(t *testing.T) {
	data := loadFixture(t, "api_error_10006.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write(data)
	}))
	defer srv.Close()

	c := testClient(srv.URL)
	_, err := c.FetchInvite(context.Background(), "invalid")
	if err == nil {
		t.Fatal("expected error for unknown invite")
	}
	var ae *APIError
	if apiErr, ok := err.(*APIError); ok {
		ae = apiErr
	}
	if ae == nil {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if ae.Code != 10006 {
		t.Errorf("code = %d, want 10006", ae.Code)
	}
}

func TestFetchMeSuccess(t *testing.T) {
	data := loadFixture(t, "user_me.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			t.Error("missing Authorization header")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	}))
	defer srv.Close()

	c := testClient(srv.URL)
	u, err := c.FetchMe(context.Background())
	if err != nil {
		t.Fatalf("FetchMe: %v", err)
	}
	if u.Username != "mybot" {
		t.Errorf("username = %q, want mybot", u.Username)
	}
	if !u.Bot {
		t.Error("expected bot=true")
	}
}

func TestFetchMeUnauthorized(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"code": 0, "message": "401: Unauthorized"}`))
	}))
	defer srv.Close()

	c := testClient(srv.URL)
	_, err := c.FetchMe(context.Background())
	if err == nil {
		t.Fatal("expected error for 401")
	}
}

func TestFetchGuilds(t *testing.T) {
	data := loadFixture(t, "guilds_me.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	}))
	defer srv.Close()

	c := testClient(srv.URL)
	guilds, err := c.FetchGuilds(context.Background())
	if err != nil {
		t.Fatalf("FetchGuilds: %v", err)
	}
	if len(guilds) != 2 {
		t.Errorf("want 2 guilds, got %d", len(guilds))
	}
	if guilds[0].Name != "Discord Developers" {
		t.Errorf("guild[0].Name = %q", guilds[0].Name)
	}
}

func TestFetchMessages(t *testing.T) {
	data := loadFixture(t, "messages_sample.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/channels/") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	}))
	defer srv.Close()

	c := testClient(srv.URL)
	msgs, err := c.FetchMessages(context.Background(), "697138785317814292", 50, "", "")
	if err != nil {
		t.Fatalf("FetchMessages: %v", err)
	}
	if len(msgs) != 2 {
		t.Errorf("want 2 messages, got %d", len(msgs))
	}
	if msgs[0].Content != "Hello world" {
		t.Errorf("msg[0].Content = %q", msgs[0].Content)
	}
}

func TestFetchMessagesBeforeParam(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.RawQuery, "before=999") {
			t.Errorf("missing before param, query = %q", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
	}))
	defer srv.Close()

	c := testClient(srv.URL)
	_, err := c.FetchMessages(context.Background(), "chan1", 10, "999", "")
	if err != nil {
		t.Fatalf("FetchMessages: %v", err)
	}
}

func TestParseInviteCode(t *testing.T) {
	cases := []struct{ in, want string }{
		{"discord-developers", "discord-developers"},
		{"https://discord.gg/discord-developers", "discord-developers"},
		{"https://discord.com/invite/abc123", "abc123"},
		{"abc123/", "abc123"},
	}
	for _, tc := range cases {
		got := ParseInviteCode(tc.in)
		if got != tc.want {
			t.Errorf("ParseInviteCode(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestGuildIconURL(t *testing.T) {
	h := "abcdef"
	got := GuildIconURL("123", &h)
	if !strings.Contains(got, "icons/123/abcdef.webp") {
		t.Errorf("icon URL = %q", got)
	}

	anim := "a_abcdef"
	got = GuildIconURL("123", &anim)
	if !strings.Contains(got, ".gif") {
		t.Errorf("animated icon should be .gif, got %q", got)
	}

	got = GuildIconURL("123", nil)
	if got != "" {
		t.Errorf("nil hash should return empty string, got %q", got)
	}
}

func TestAvatarURL(t *testing.T) {
	h := "userhash"
	got := AvatarURL("456", &h)
	if !strings.Contains(got, "avatars/456/userhash.webp") {
		t.Errorf("avatar URL = %q", got)
	}
}

func TestInviteToServerRecord(t *testing.T) {
	data := loadFixture(t, "invite_discord_developers.json")
	var inv Invite
	if err := json.Unmarshal(data, &inv); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	rec := InviteToServerRecord(&inv)
	if rec.Name != "Discord Developers" {
		t.Errorf("Name = %q", rec.Name)
	}
	if rec.Members != 152000 {
		t.Errorf("Members = %d", rec.Members)
	}
	if rec.InviteType != "permanent" {
		t.Errorf("InviteType = %q", rec.InviteType)
	}
	if rec.InviteChannel != "rules" {
		t.Errorf("InviteChannel = %q", rec.InviteChannel)
	}
	if !strings.Contains(rec.Features, "COMMUNITY") {
		t.Errorf("Features = %q", rec.Features)
	}
}

func TestInviteToInviteRecord(t *testing.T) {
	data := loadFixture(t, "invite_discord_developers.json")
	var inv Invite
	if err := json.Unmarshal(data, &inv); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	rec := InviteToInviteRecord(&inv)
	if rec.Code != "discord-developers" {
		t.Errorf("Code = %q", rec.Code)
	}
	if rec.GuildName != "Discord Developers" {
		t.Errorf("GuildName = %q", rec.GuildName)
	}
	if !strings.HasPrefix(rec.InviteURL, "https://discord.gg/") {
		t.Errorf("InviteURL = %q", rec.InviteURL)
	}
}

func TestUserToUserRecord(t *testing.T) {
	data := loadFixture(t, "user_me.json")
	var u User
	if err := json.Unmarshal(data, &u); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	rec := UserToRecord(&u)
	if rec.Username != "mybot" {
		t.Errorf("Username = %q", rec.Username)
	}
	if rec.GlobalName != "My Bot" {
		t.Errorf("GlobalName = %q", rec.GlobalName)
	}
	if rec.AvatarURL == "" {
		t.Error("AvatarURL should not be empty")
	}
}

func TestMessageToRecord(t *testing.T) {
	data := loadFixture(t, "messages_sample.json")
	var msgs []Message
	if err := json.Unmarshal(data, &msgs); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	rec := MessageToRecord(msgs[0])
	if rec.Content != "Hello world" {
		t.Errorf("Content = %q", rec.Content)
	}
	if rec.AuthorName != "someuser" {
		t.Errorf("AuthorName = %q", rec.AuthorName)
	}
	if rec.Reactions != 1 {
		t.Errorf("Reactions = %d, want 1", rec.Reactions)
	}

	rec2 := MessageToRecord(msgs[1])
	if rec2.Attachments != 1 {
		t.Errorf("Attachments = %d, want 1", rec2.Attachments)
	}
	if !rec2.Pinned {
		t.Error("msg[1] should be pinned")
	}
	if rec2.EditedTimestamp == "" {
		t.Error("EditedTimestamp should not be empty for edited message")
	}
}
