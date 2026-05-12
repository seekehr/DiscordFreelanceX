package tasks

import (
	"fmt"
	"sort"

	"github.com/bwmarrin/discordgo"
	"github.com/seekehr/DiscordFreelanceX/internal"
	"github.com/seekehr/DiscordFreelanceX/internal/utils"
)

// AnalyzeLastMessages fetches the most recent messages from every configured
// channel and returns them as structured AnalysisEntry slices, sorted newest-first.
// Guild and channel sections are separated by visual dividers.
func AnalyzeLastMessages(numberofmessages int, s *discordgo.Session, cfg *internal.Config) ([]internal.AnalysisEntry, error) {
	var entries []internal.AnalysisEntry

	for i, server := range cfg.Servers {
		if i > 0 {
			entries = append(entries, internal.AnalysisEntry{
				Text: "============================",
			})
		}

		guildName := utils.GetGuildNameFromID(s, server.GuildID)
		entries = append(entries, internal.AnalysisEntry{
			Text: fmt.Sprintf("Processing guild: %s", guildName),
		})

		for j, channelID := range server.ChannelIDs {
			if j > 0 {
				entries = append(entries, internal.AnalysisEntry{
					Text: "-----------------------------------------------",
				})
			}

			channelName := utils.GetChannelNameFromID(s, channelID)
			messages, err := s.ChannelMessages(channelID, numberofmessages, "", "", "")
			if err != nil {
				entries = append(entries, internal.AnalysisEntry{
					Text: fmt.Sprintf("Failed to read channel %s: %v", channelName, err),
				})
				continue
			}

			sort.Slice(messages, func(a, b int) bool {
				return messages[a].ID > messages[b].ID
			})

			entries = append(entries, internal.AnalysisEntry{
				Text: fmt.Sprintf("%s: fetched %d messages", channelName, len(messages)),
			})
			for _, msg := range messages {
				content := msg.Content
				if content == "" {
					switch {
					case len(msg.Embeds) > 0:
						content = fmt.Sprintf("%s", ParseEmbed(msg.Embeds[0]))
					default:
						content = "[no text content]"
					}
				}

				messageURL := fmt.Sprintf("https://discord.com/channels/%s/%s/%s", server.GuildID, channelID, msg.ID)
				entries = append(entries, internal.AnalysisEntry{
					Text:       fmt.Sprintf("[%s]: %s", msg.Author.Username, content),
					MessageURL: messageURL,
				})
			}
		}
	}

	return entries, nil
}
