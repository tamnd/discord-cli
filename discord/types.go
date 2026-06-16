// Package discord provides a client for the Discord REST API v10.
// Public endpoints work without authentication; most require a Bot token
// set in the DISCORD_TOKEN environment variable.
package discord

import "fmt"

// Invite is the raw Discord invite object returned by GET /invites/{code}.
type Invite struct {
	Type                int      `json:"type"`
	Code                string   `json:"code"`
	Channel             *Channel `json:"channel"`
	Guild               *Guild   `json:"guild"`
	Inviter             *User    `json:"inviter"`
	Uses                int      `json:"uses"`
	MaxUses             int      `json:"max_uses"`
	MaxAge              int      `json:"max_age"`
	Temporary           bool     `json:"temporary"`
	ExpiresAt           *string  `json:"expires_at"`
	ApproxMemberCount   int      `json:"approximate_member_count"`
	ApproxPresenceCount int      `json:"approximate_presence_count"`
}

// Guild is a Discord server object.
type Guild struct {
	ID                       string   `json:"id"`
	Name                     string   `json:"name"`
	Description              *string  `json:"description"`
	Icon                     *string  `json:"icon"`
	Banner                   *string  `json:"banner"`
	Splash                   *string  `json:"splash"`
	Features                 []string `json:"features"`
	VerificationLevel        int      `json:"verification_level"`
	VanityURLCode            *string  `json:"vanity_url_code"`
	NSFWLevel                int      `json:"nsfw_level"`
	PremiumSubscriptionCount int      `json:"premium_subscription_count"`
	PremiumTier              int      `json:"premium_tier"`
	Owner                    bool     `json:"owner"`
	Permissions              string   `json:"permissions"`
}

// Channel is a Discord channel object.
type Channel struct {
	ID   string `json:"id"`
	Type int    `json:"type"`
	Name string `json:"name"`
}

// User is a Discord user object.
type User struct {
	ID            string  `json:"id"`
	Username      string  `json:"username"`
	Discriminator string  `json:"discriminator"`
	GlobalName    *string `json:"global_name"`
	Avatar        *string `json:"avatar"`
	Bot           bool    `json:"bot"`
	System        bool    `json:"system"`
	MFAEnabled    bool    `json:"mfa_enabled"`
	Verified      bool    `json:"verified"`
	Locale        *string `json:"locale"`
	PublicFlags   int     `json:"public_flags"`
}

// Message is a Discord message object.
type Message struct {
	ID              string  `json:"id"`
	ChannelID       string  `json:"channel_id"`
	Author          User    `json:"author"`
	Content         string  `json:"content"`
	Timestamp       string  `json:"timestamp"`
	EditedTimestamp *string `json:"edited_timestamp"`
	Pinned          bool    `json:"pinned"`
	Type            int     `json:"type"`
	Attachments     []any   `json:"attachments"`
	Embeds          []any   `json:"embeds"`
	Reactions       []any   `json:"reactions"`
}

// APIError is the Discord API error response envelope.
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("discord API error %d: %s", e.Code, e.Message)
}

// Output record types.

// ServerRecord is the output of the `server` command.
type ServerRecord struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Features          string `json:"features"`
	IconURL           string `json:"icon_url"`
	Members           int    `json:"members"`
	Online            int    `json:"online"`
	VerificationLevel int    `json:"verification_level"`
	PremiumTier       int    `json:"premium_tier"`
	Boosts            int    `json:"boosts"`
	VanityURL         string `json:"vanity_url"`
	NSFWLevel         int    `json:"nsfw_level"`
	InviteChannel     string `json:"invite_channel"`
	InviteExpiresAt   string `json:"invite_expires_at"`
	InviteType        string `json:"invite_type"`
}

// InviteRecord is the output of the `invite` command.
type InviteRecord struct {
	Code        string `json:"code"`
	GuildID     string `json:"guild_id"`
	GuildName   string `json:"guild_name"`
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	Inviter     string `json:"inviter"`
	Uses        int    `json:"uses"`
	MaxUses     int    `json:"max_uses"`
	MaxAge      int    `json:"max_age"`
	Temporary   bool   `json:"temporary"`
	ExpiresAt   string `json:"expires_at"`
	Members     int    `json:"members"`
	Online      int    `json:"online"`
	InviteURL   string `json:"invite_url"`
}

// UserRecord is the output of the `me` command.
type UserRecord struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	GlobalName    string `json:"global_name"`
	Discriminator string `json:"discriminator"`
	AvatarURL     string `json:"avatar_url"`
	Bot           bool   `json:"bot"`
	Verified      bool   `json:"verified"`
	Locale        string `json:"locale"`
	PublicFlags   int    `json:"public_flags"`
}

// GuildRecord is one row in the `servers` output.
type GuildRecord struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IconURL     string `json:"icon_url"`
	Owner       bool   `json:"owner"`
	Permissions string `json:"permissions"`
	Features    string `json:"features"`
}

// MessageRecord is one row in the `messages` output.
type MessageRecord struct {
	ID               string `json:"id"`
	ChannelID        string `json:"channel_id"`
	AuthorID         string `json:"author_id"`
	AuthorName       string `json:"author_name"`
	AuthorGlobalName string `json:"author_global_name"`
	Content          string `json:"content"`
	Timestamp        string `json:"timestamp"`
	EditedTimestamp  string `json:"edited_timestamp"`
	Pinned           bool   `json:"pinned"`
	Type             int    `json:"type"`
	Attachments      int    `json:"attachments"`
	Embeds           int    `json:"embeds"`
	Reactions        int    `json:"reactions"`
}
