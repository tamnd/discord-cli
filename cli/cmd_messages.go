package cli

import (
	"context"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/discord-cli/discord"
)

type messagesCmd struct {
	out    outputFlags
	limit  int
	before string
	after  string
}

func newMessagesCmd() kit.Command {
	c := &messagesCmd{}
	return kit.Command{
		Use:   "messages <channel-id>",
		Short: "Fetch messages from a channel (requires DISCORD_TOKEN)",
		Long: `Fetch recent messages from a Discord channel.

Requires a Discord Bot token set in the DISCORD_TOKEN environment variable.
The bot must be a member of the server and have READ_MESSAGE_HISTORY permission.

Examples:
  discord messages 697138785317814292
  discord messages 697138785317814292 --limit 100 -o jsonl
  discord messages 697138785317814292 --before 1234567890123456789`,
		Args:  kit.ExactArgs(1),
		Flags: c.flags,
		Run:   c.run,
	}
}

func (c *messagesCmd) flags(f *kit.FlagSet) {
	f.StringVarP(&c.out.format, "output", "o", "jsonl", "Output format: table|json|jsonl|csv|tsv|raw")
	f.StringSliceVar(&c.out.fields, "fields", nil, "Select and reorder output columns")
	f.BoolVar(&c.out.noHeader, "no-header", false, "Suppress header row")
	f.StringVar(&c.out.tmpl, "template", "", "Go text/template applied per record")
	f.IntVar(&c.limit, "limit", 50, "Number of messages to fetch (default: 50, max: 100)")
	f.StringVar(&c.before, "before", "", "Fetch messages before this message ID")
	f.StringVar(&c.after, "after", "", "Fetch messages after this message ID")
}

func (c *messagesCmd) run(ctx context.Context, args []string) error {
	if err := requireToken(); err != nil {
		return fatalf(1, "%v", err)
	}
	client := newDiscordClient()
	msgs, err := client.FetchMessages(ctx, args[0], c.limit, c.before, c.after)
	if err != nil {
		return fatalf(1, "%v", err)
	}
	if len(msgs) == 0 {
		return fatalf(3, "no messages found in channel %s", args[0])
	}

	rend := c.out.renderer("jsonl")
	for _, m := range msgs {
		if err := rend.Render(discord.MessageToRecord(m)); err != nil {
			return err
		}
	}
	return nil
}
