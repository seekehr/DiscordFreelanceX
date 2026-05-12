package utils

import "github.com/bwmarrin/discordgo"

// GetGuildNameFromID resolves a guild snowflake ID to its human-readable name
// via the Discord API. Falls back to the raw ID on failure.
func GetGuildNameFromID(s *discordgo.Session, guildID string) string {
	guild, err := s.Guild(guildID)
	if err != nil {
		return guildID
	}
	return guild.Name
}

// GetChannelNameFromID resolves a channel snowflake ID to its display name
// (prefixed with #) via the Discord API. Falls back to the raw ID on failure.
func GetChannelNameFromID(s *discordgo.Session, channelID string) string {
	channel, err := s.Channel(channelID)
	if err != nil {
		return channelID
	}
	return "#" + channel.Name
}
