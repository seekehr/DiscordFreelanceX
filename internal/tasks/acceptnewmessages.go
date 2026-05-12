package tasks

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"github.com/bwmarrin/discordgo"
	"github.com/seekehr/DiscordFreelanceX/internal"
	"github.com/seekehr/DiscordFreelanceX/internal/gui"
	"github.com/seekehr/DiscordFreelanceX/internal/parsers"
	"github.com/seekehr/DiscordFreelanceX/internal/utils"
)

// AcceptNewMessages registers handlers that display incoming messages and
// new forum posts from all configured channels in the GUI's RichText widget.
// Each incoming event also triggers a Windows toast notification.
func AcceptNewMessages(s *discordgo.Session, cfg *internal.Config, a fyne.App, rt *widget.RichText) {
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

		title := fmt.Sprintf("New message in %s", channelName)
		line := fmt.Sprintf("[%s]: %s", m.Author.Username, content)

		gui.AppendAnalysisEntries(rt, []internal.AnalysisEntry{
			{Text: fmt.Sprintf("NEW in %s | %s", channelName, line), MessageURL: messageURL},
		})
		utils.Notify(a, title, line)
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
		utils.Notify(a, fmt.Sprintf("New forum post in %s", forumName), content)
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
