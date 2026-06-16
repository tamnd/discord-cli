package discord

import (
	"fmt"
	"strings"
)

// ParseInviteCode extracts the invite code from a full Discord URL or returns
// the input as-is if it's already a bare code.
func ParseInviteCode(input string) string {
	input = strings.TrimSpace(input)
	input = strings.TrimPrefix(input, "https://discord.gg/")
	input = strings.TrimPrefix(input, "https://discord.com/invite/")
	input = strings.SplitN(input, "?", 2)[0]
	input = strings.TrimRight(input, "/")
	return input
}

// GuildIconURL returns the CDN URL for a guild icon. Returns "" if hash is nil or empty.
func GuildIconURL(guildID string, hash *string) string {
	if hash == nil || *hash == "" {
		return ""
	}
	ext := "webp"
	if strings.HasPrefix(*hash, "a_") {
		ext = "gif"
	}
	return fmt.Sprintf("https://cdn.discordapp.com/icons/%s/%s.%s", guildID, *hash, ext)
}

// AvatarURL returns the CDN URL for a user avatar. Returns "" if hash is nil or empty.
func AvatarURL(userID string, hash *string) string {
	if hash == nil || *hash == "" {
		return ""
	}
	ext := "webp"
	if strings.HasPrefix(*hash, "a_") {
		ext = "gif"
	}
	return fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.%s", userID, *hash, ext)
}

// InviteToServerRecord converts a raw Invite to a ServerRecord.
func InviteToServerRecord(inv *Invite) ServerRecord {
	rec := ServerRecord{
		Members: inv.ApproxMemberCount,
		Online:  inv.ApproxPresenceCount,
	}
	if inv.Guild != nil {
		g := inv.Guild
		rec.ID = g.ID
		rec.Name = g.Name
		if g.Description != nil {
			rec.Description = *g.Description
		}
		rec.Features = strings.Join(g.Features, ",")
		rec.IconURL = GuildIconURL(g.ID, g.Icon)
		rec.VerificationLevel = g.VerificationLevel
		rec.PremiumTier = g.PremiumTier
		rec.Boosts = g.PremiumSubscriptionCount
		if g.VanityURLCode != nil {
			rec.VanityURL = *g.VanityURLCode
		}
		rec.NSFWLevel = g.NSFWLevel
	}
	if inv.Channel != nil {
		rec.InviteChannel = inv.Channel.Name
	}
	if inv.ExpiresAt != nil {
		rec.InviteExpiresAt = *inv.ExpiresAt
		rec.InviteType = "temporary"
	} else {
		rec.InviteType = "permanent"
	}
	return rec
}

// InviteToInviteRecord converts a raw Invite to an InviteRecord.
func InviteToInviteRecord(inv *Invite) InviteRecord {
	rec := InviteRecord{
		Code:      inv.Code,
		Uses:      inv.Uses,
		MaxUses:   inv.MaxUses,
		MaxAge:    inv.MaxAge,
		Temporary: inv.Temporary,
		Members:   inv.ApproxMemberCount,
		Online:    inv.ApproxPresenceCount,
		InviteURL: "https://discord.gg/" + inv.Code,
	}
	if inv.Guild != nil {
		rec.GuildID = inv.Guild.ID
		rec.GuildName = inv.Guild.Name
	}
	if inv.Channel != nil {
		rec.ChannelID = inv.Channel.ID
		rec.ChannelName = inv.Channel.Name
	}
	if inv.Inviter != nil {
		rec.Inviter = inv.Inviter.Username
	}
	if inv.ExpiresAt != nil {
		rec.ExpiresAt = *inv.ExpiresAt
	}
	return rec
}

// UserToRecord converts a raw User to a UserRecord.
func UserToRecord(u *User) UserRecord {
	rec := UserRecord{
		ID:            u.ID,
		Username:      u.Username,
		Discriminator: u.Discriminator,
		AvatarURL:     AvatarURL(u.ID, u.Avatar),
		Bot:           u.Bot,
		Verified:      u.Verified,
		PublicFlags:   u.PublicFlags,
	}
	if u.GlobalName != nil {
		rec.GlobalName = *u.GlobalName
	}
	if u.Locale != nil {
		rec.Locale = *u.Locale
	}
	return rec
}

// GuildToRecord converts a raw Guild to a GuildRecord.
func GuildToRecord(g Guild) GuildRecord {
	return GuildRecord{
		ID:          g.ID,
		Name:        g.Name,
		IconURL:     GuildIconURL(g.ID, g.Icon),
		Owner:       g.Owner,
		Permissions: g.Permissions,
		Features:    strings.Join(g.Features, ","),
	}
}

// MessageToRecord converts a raw Message to a MessageRecord.
func MessageToRecord(m Message) MessageRecord {
	rec := MessageRecord{
		ID:          m.ID,
		ChannelID:   m.ChannelID,
		AuthorID:    m.Author.ID,
		AuthorName:  m.Author.Username,
		Content:     m.Content,
		Timestamp:   m.Timestamp,
		Pinned:      m.Pinned,
		Type:        m.Type,
		Attachments: len(m.Attachments),
		Embeds:      len(m.Embeds),
		Reactions:   len(m.Reactions),
	}
	if m.Author.GlobalName != nil {
		rec.AuthorGlobalName = *m.Author.GlobalName
	}
	if m.EditedTimestamp != nil {
		rec.EditedTimestamp = *m.EditedTimestamp
	}
	return rec
}
