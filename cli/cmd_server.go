package cli

import (
	"context"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/discord-cli/discord"
)

type serverCmd struct {
	out outputFlags
}

func newServerCmd() kit.Command {
	c := &serverCmd{}
	return kit.Command{
		Use:   "server <invite-code>",
		Short: "Look up a public Discord server via invite code",
		Long: `Fetch public server info by resolving a Discord invite code.

Accepts a bare invite code or a full Discord invite URL.
No authentication required.

Examples:
  discord server discord-developers
  discord server https://discord.gg/reactiflux -o json
  discord server python | jq .members`,
		Args:  kit.ExactArgs(1),
		Flags: c.flags,
		Run:   c.run,
	}
}

func (c *serverCmd) flags(f *kit.FlagSet) {
	f.StringVarP(&c.out.format, "output", "o", "table", "Output format: table|json|jsonl|csv|tsv|raw")
	f.StringSliceVar(&c.out.fields, "fields", nil, "Select and reorder output columns")
	f.BoolVar(&c.out.noHeader, "no-header", false, "Suppress header row")
	f.StringVar(&c.out.tmpl, "template", "", "Go text/template applied per record")
}

func (c *serverCmd) run(ctx context.Context, args []string) error {
	code := discord.ParseInviteCode(args[0])
	client := newDiscordClient()
	inv, err := client.FetchInvite(ctx, code)
	if err != nil {
		return fatalf(1, "%v", err)
	}
	rec := discord.InviteToServerRecord(inv)
	return c.out.renderer("table").Render(rec)
}
