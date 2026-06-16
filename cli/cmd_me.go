package cli

import (
	"context"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/discord-cli/discord"
)

type meCmd struct {
	out outputFlags
}

func newMeCmd() kit.Command {
	c := &meCmd{}
	return kit.Command{
		Use:   "me",
		Short: "Show the authenticated bot user profile (requires DISCORD_TOKEN)",
		Long: `Fetch the profile of the authenticated Discord bot user.

Requires a Discord Bot token set in the DISCORD_TOKEN environment variable.

Examples:
  DISCORD_TOKEN=Bot.xxx discord me
  discord me -o json`,
		Args:  kit.NoArgs,
		Flags: c.flags,
		Run:   c.run,
	}
}

func (c *meCmd) flags(f *kit.FlagSet) {
	f.StringVarP(&c.out.format, "output", "o", "table", "Output format: table|json|jsonl|csv|tsv|raw")
	f.StringSliceVar(&c.out.fields, "fields", nil, "Select and reorder output columns")
	f.BoolVar(&c.out.noHeader, "no-header", false, "Suppress header row")
	f.StringVar(&c.out.tmpl, "template", "", "Go text/template applied per record")
}

func (c *meCmd) run(ctx context.Context, _ []string) error {
	if err := requireToken(); err != nil {
		return fatalf(1, "%v", err)
	}
	client := newDiscordClient()
	u, err := client.FetchMe(ctx)
	if err != nil {
		return fatalf(1, "%v", err)
	}
	rec := discord.UserToRecord(u)
	return c.out.renderer("table").Render(rec)
}
