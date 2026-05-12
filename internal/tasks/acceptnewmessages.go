package tasks

import (
	"fmt"

	"fyne.io/fyne/v2/widget"

	"github.com/bwmarrin/discordgo"
	"github.com/seekehr/DiscordFreelanceX/internal"
	"github.com/seekehr/DiscordFreelanceX/internal/gui"
	"github.com/seekehr/DiscordFreelanceX/internal/parsers"
	"github.com/seekehr/DiscordFreelanceX/internal/utils"
)

// AcceptNewMessages registers handlers that display incoming messages and
// new forum posts from all configured channels in the GUI's RichText widget.
func AcceptNewMessages(s *discordgo.Session, cfg *internal.Config, rt *widget.RichText) {
	watched := buildWatchedChannels(cfg)

	s.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		guildID, ok := watched[m.ChannelID]
		if !ok {
			return
		}

		content := m.Content
		if content == "" {
			switch {
			case len(m.Embeds) > 0:
				content = parsers.ParseEmbed(m.Embeds[0])
			default:
				content = "[no text content]"
			}
		}

		channelName := utils.GetChannelNameFromID(s, m.ChannelID)
		messageURL := fmt.Sprintf("https://discord.com/channels/%s/%s/%s", guildID, m.ChannelID, m.ID)

		gui.AppendAnalysisEntries(rt, []internal.AnalysisEntry{
			{Text: fmt.Sprintf("NEW in %s | [%s]: %s", channelName, m.Author.Username, content), MessageURL: messageURL},
		})
	})

	s.AddHandler(func(_ *discordgo.Session, tc *discordgo.ThreadCreate) {
		if tc.ParentID == "" {
			return
		}
		guildID, ok := watched[tc.ParentID]
		if !ok {
			return
		}

		content := parsers.ParseForumPost(s, tc.Channel)
		threadURL := fmt.Sprintf("https://discord.com/channels/%s/%s", guildID, tc.ID)
		forumName := utils.GetChannelNameFromID(s, tc.ParentID)

		gui.AppendAnalysisEntries(rt, []internal.AnalysisEntry{
			{Text: fmt.Sprintf("NEW POST in %s | %s", forumName, content), MessageURL: threadURL},
		})
	})
}

// buildWatchedChannels creates a channelID -> guildID lookup from the config.
func buildWatchedChannels(cfg *internal.Config) map[string]string {
	m := make(map[string]string, len(cfg.Servers)*4)
	for _, server := range cfg.Servers {
		for _, ch := range server.ChannelIDs {
			m[ch] = server.GuildID
		}
	}
	return m
}
