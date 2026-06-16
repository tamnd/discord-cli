package cli

import (
	"context"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/discord-cli/discord"
)

type serversCmd struct {
	out   outputFlags
	limit int
}

func newServersCmd() kit.Command {
	c := &serversCmd{}
	return kit.Command{
		Use:   "servers",
		Short: "List servers the bot has joined (requires DISCORD_TOKEN)",
		Long: `List all Discord servers (guilds) the authenticated bot has joined.

Requires a Discord Bot token set in the DISCORD_TOKEN environment variable.

Examples:
  discord servers
  discord servers -o jsonl | jq .name
  discord servers --limit 10 -o csv`,
		Args:  kit.NoArgs,
		Flags: c.flags,
		Run:   c.run,
	}
}

func (c *serversCmd) flags(f *kit.FlagSet) {
	f.StringVarP(&c.out.format, "output", "o", "table", "Output format: table|json|jsonl|csv|tsv|raw")
	f.StringSliceVar(&c.out.fields, "fields", nil, "Select and reorder output columns")
	f.BoolVar(&c.out.noHeader, "no-header", false, "Suppress header row")
	f.StringVar(&c.out.tmpl, "template", "", "Go text/template applied per record")
	f.IntVar(&c.limit, "limit", 0, "Max servers to return (0 = all)")
}

func (c *serversCmd) run(ctx context.Context, _ []string) error {
	if err := requireToken(); err != nil {
		return fatalf(1, "%v", err)
	}
	client := newDiscordClient()
	guilds, err := client.FetchGuilds(ctx)
	if err != nil {
		return fatalf(1, "%v", err)
	}
	if len(guilds) == 0 {
		return fatalf(3, "bot has not joined any servers")
	}

	rend := c.out.renderer("table")
	for i, g := range guilds {
		if c.limit > 0 && i >= c.limit {
			break
		}
		if err := rend.Render(discord.GuildToRecord(g)); err != nil {
			return err
		}
	}
	return nil
}
