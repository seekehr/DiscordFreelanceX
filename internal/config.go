package internal

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Bot     BotConfig      `yaml:"bot"`
	Servers []ServerConfig `yaml:"servers"`
}

type BotConfig struct {
	Message              string `yaml:"message"`
	AnalyzeLastXMessages int    `yaml:"analyze-last-x-messages"`
}

type ServerConfig struct {
	GuildID    string   `yaml:"guild_id"`
	ChannelIDs []string `yaml:"channel_ids"`
}

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
