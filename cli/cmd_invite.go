package cli

import (
	"context"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/discord-cli/discord"
)

type inviteCmd struct {
	out outputFlags
}

func newInviteCmd() kit.Command {
	c := &inviteCmd{}
	return kit.Command{
		Use:   "invite <code>",
		Short: "Get details about a Discord invite",
		Long: `Fetch detailed invite metadata including expiry, uses, and inviter.

Accepts a bare invite code or a full Discord invite URL.
No authentication required.

Examples:
  discord invite discord-developers
  discord invite abc123 -o jsonl`,
		Args:  kit.ExactArgs(1),
		Flags: c.flags,
		Run:   c.run,
	}
}

func (c *inviteCmd) flags(f *kit.FlagSet) {
	f.StringVarP(&c.out.format, "output", "o", "table", "Output format: table|json|jsonl|csv|tsv|raw")
	f.StringSliceVar(&c.out.fields, "fields", nil, "Select and reorder output columns")
	f.BoolVar(&c.out.noHeader, "no-header", false, "Suppress header row")
	f.StringVar(&c.out.tmpl, "template", "", "Go text/template applied per record")
}

func (c *inviteCmd) run(ctx context.Context, args []string) error {
	code := discord.ParseInviteCode(args[0])
	client := newDiscordClient()
	inv, err := client.FetchInvite(ctx, code)
	if err != nil {
		return fatalf(1, "%v", err)
	}
	rec := discord.InviteToInviteRecord(inv)
	return c.out.renderer("table").Render(rec)
}
