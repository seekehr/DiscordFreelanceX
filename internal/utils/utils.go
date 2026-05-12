package utils

import "github.com/bwmarrin/discordgo"

func GetGuildNameFromID(s *discordgo.Session, guildID string) string {
	guild, err := s.Guild(guildID)
	if err != nil {
		return guildID
	}
	return guild.Name
}

func GetChannelNameFromID(s *discordgo.Session, channelID string) string {
	channel, err := s.Channel(channelID)
	if err != nil {
		return channelID
	}
	return "#" + channel.Name
}
