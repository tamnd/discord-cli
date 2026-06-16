// Package cli assembles the discord command tree.
package cli

import (
	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/discord-cli/discord"
)

// Build metadata, set via -ldflags at release time.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// NewApp assembles the kit application. Commands are registered as kit
// escape-hatch commands and wired to the Discord API client.
func NewApp() *kit.App {
	id := discord.Domain{}.Info().Identity
	id.Version = Version

	app := kit.New(id)
	(discord.Domain{}).Register(app)

	app.AddCommand(newServerCmd())
	app.AddCommand(newInviteCmd())
	app.AddCommand(newMeCmd())
	app.AddCommand(newServersCmd())
	app.AddCommand(newMessagesCmd())
	app.AddCommand(newVersionCmd())

	return app
}
