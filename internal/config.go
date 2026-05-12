package internal

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the top-level application configuration loaded from YAML.
type Config struct {
	Bot     BotConfig      `yaml:"bot"`
	Servers []ServerConfig `yaml:"servers"`
}

// BotConfig holds bot-specific settings such as the auto-reply message
// and the number of messages to fetch per channel.
type BotConfig struct {
	Message              string `yaml:"message"`
	AnalyzeLastXMessages int    `yaml:"analyze-last-x-messages"`
}

// ServerConfig identifies a single Discord guild and the channels to monitor.
type ServerConfig struct {
	GuildID    string   `yaml:"guild_id"`
	ChannelIDs []string `yaml:"channel_ids"`
}

// AnalysisEntry represents a single line in the analysis output.
// MessageURL is non-empty only for lines that correspond to a Discord message.
type AnalysisEntry struct {
	Text       string
	MessageURL string
}

// LoadConfig reads and parses the YAML configuration file at path.
// It returns an error if the file is missing, malformed, or lacks required fields.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if cfg.Bot.Message == "" {
		return nil, fmt.Errorf("bot.message is required")
	}

	return &cfg, nil
}
