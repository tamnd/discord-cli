package cli

import (
	"fmt"
	"os"

	"github.com/tamnd/discord-cli/discord"
	"github.com/tamnd/discord-cli/pkg/render"
)

// outputFlags holds the shared output format flags.
type outputFlags struct {
	format   string
	fields   []string
	noHeader bool
	tmpl     string
}

func (o *outputFlags) renderer(defaultFormat string) *render.Renderer {
	f := render.Format(o.format)
	if !f.Valid() {
		f = render.Format(defaultFormat)
	}
	return render.New(os.Stdout, f, o.fields, o.noHeader, o.tmpl)
}

// newDiscordClient creates a Discord client from the default config.
func newDiscordClient() *discord.Client {
	return discord.NewClient(discord.DefaultConfig())
}

// requireToken returns an error if DISCORD_TOKEN is not set.
func requireToken() error {
	cfg := discord.DefaultConfig()
	if cfg.Token == "" {
		return fmt.Errorf(`DISCORD_TOKEN is not set.
Get a bot token from https://discord.com/developers/applications
Set it with: export DISCORD_TOKEN=your_token_here`)
	}
	return nil
}

// fatal prints msg to stderr and returns an exitErr.
type exitErr struct {
	code int
	msg  string
}

func (e *exitErr) Error() string { return e.msg }

func fatalf(code int, format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintln(os.Stderr, "discord:", msg)
	return &exitErr{code: code, msg: msg}
}
