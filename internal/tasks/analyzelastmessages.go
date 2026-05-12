package tasks

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/seekehr/DiscordFreelanceX/internal"
)

func AnalyzeLastMessages(numberofmessages int, s *discordgo.Session, cfg *internal.Config) (string, error) {
	var sb strings.Builder

	for _, server := range cfg.Servers {
		sb.WriteString(fmt.Sprintf("Processing guild: %s\n", server.GuildID))

		for _, channelID := range server.ChannelIDs {
			messages, err := s.ChannelMessages(channelID, numberofmessages, "", "", "")
			if err != nil {
				sb.WriteString(fmt.Sprintf("  Failed to read channel %s: %v\n", channelID, err))
				continue
			}

			sb.WriteString(fmt.Sprintf("  Channel %s: fetched %d messages\n", channelID, len(messages)))
			for _, msg := range messages {
				content := msg.Content
				if content == "" {
					switch {
					case len(msg.Embeds) > 0:
						content = fmt.Sprintf("[embed: %s]", ParseEmbed(msg.Embeds[0]))
					case len(msg.Attachments) > 0:
						content = fmt.Sprintf("[attachment: %s]", msg.Attachments[0].Filename)
					case len(msg.StickerItems) > 0:
						content = fmt.Sprintf("[sticker: %s]", msg.StickerItems[0].Name)
					default:
						content = "[no text content]"
					}
				}
				sb.WriteString(fmt.Sprintf("    [%s] %s: %s\n", msg.Timestamp, msg.Author.Username, content))
			}
		}
	}

	return sb.String(), nil
}
