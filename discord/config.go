package discord

import (
	"os"
	"time"
)

// Config holds the HTTP client configuration for the Discord API.
type Config struct {
	Token     string
	BaseURL   string
	UserAgent string
	Timeout   time.Duration
	Retries   int
}

// DefaultConfig returns a Config populated from the environment and sensible defaults.
func DefaultConfig() Config {
	return Config{
		Token:     os.Getenv("DISCORD_TOKEN"),
		BaseURL:   "https://discord.com/api/v10",
		UserAgent: "discord-cli/0.1.0 (https://github.com/tamnd/discord-cli)",
		Timeout:   30 * time.Second,
		Retries:   3,
	}
}
