package discord

import (
	"net/url"
	"strings"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/any-cli/kit/errs"
)

// Domain is the discord kit domain, registered for binary identity and version output.
func init() { kit.Register(Domain{}) }

// Domain is the discord driver.
type Domain struct{}

// Info returns the domain identity.
func (Domain) Info() kit.DomainInfo {
	return kit.DomainInfo{
		Scheme: "discord",
		Hosts:  []string{"discord.com"},
		Identity: kit.Identity{
			Binary: "discord",
			Short:  "Look up Discord servers and read messages from the terminal",
			Long: `Look up Discord servers and read messages from the terminal.

Public commands (server, invite) work without any credentials.
Authenticated commands (me, servers, messages) require a Discord Bot token:

  export DISCORD_TOKEN=your_bot_token_here`,
			Site: "discord.com",
			Repo: "https://github.com/tamnd/discord-cli",
		},
	}
}

// Register is a no-op; commands are added as kit escape-hatch commands in cli/root.go.
func (Domain) Register(app *kit.App) {}

// Classify turns a discord.com URL into a (type, id) pair.
func (Domain) Classify(input string) (uriType, id string, err error) {
	id = refPath(input)
	if id == "" {
		return "", "", errs.Usage("unrecognized discord reference: %q", input)
	}
	return "page", id, nil
}

// Locate returns the URL for a (type, id).
func (Domain) Locate(uriType, id string) (string, error) {
	if uriType != "page" {
		return "", errs.Usage("discord has no resource type %q", uriType)
	}
	return "https://discord.com/" + strings.Trim(id, "/"), nil
}

func refPath(input string) string {
	input = strings.TrimSpace(input)
	if u, err := url.Parse(input); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
		return strings.Trim(u.Path, "/")
	}
	return strings.Trim(input, "/")
}
