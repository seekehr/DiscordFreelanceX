package tasks

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/seekehr/DiscordFreelanceX/internal"
	"github.com/seekehr/DiscordFreelanceX/internal/parsers"
	"github.com/seekehr/DiscordFreelanceX/internal/utils"
)

// AnalyzeLastMessages fetches the most recent messages from every configured
// channel and returns them keyed by guild ID, sorted newest-first per channel.
// Forum channels are handled separately by fetching their active threads.
func AnalyzeLastMessages(numberofmessages int, s *discordgo.Session, cfg *internal.Config) (map[string][]internal.AnalysisEntry, error) {
	result := make(map[string][]internal.AnalysisEntry, len(cfg.Servers))

	for _, server := range cfg.Servers {
		var entries []internal.AnalysisEntry

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

			ch, err := s.Channel(channelID)
			if err != nil {
				entries = append(entries, internal.AnalysisEntry{
					Text: fmt.Sprintf("Failed to read channel %s: %v", channelID, err),
				})
				continue
			}

			if ch.Type == discordgo.ChannelTypeGuildForum {
				entries = append(entries, analyzeForumChannel(s, server.GuildID, ch, numberofmessages, cfg.Bot.Keywords)...)
			} else {
				entries = append(entries, analyzeTextChannel(s, server.GuildID, channelID, numberofmessages, cfg.Bot.Keywords)...)
			}
		}

		result[server.GuildID] = entries
	}

	return result, nil
}

func analyzeTextChannel(s *discordgo.Session, guildID, channelID string, limit int, keywords []string) []internal.AnalysisEntry {
	var entries []internal.AnalysisEntry

	channelName := utils.GetChannelNameFromID(s, channelID)
	messages, err := s.ChannelMessages(channelID, limit, "", "", "")
	if err != nil {
		return append(entries, internal.AnalysisEntry{
			Text: fmt.Sprintf("Failed to read channel %s: %v", channelName, err),
		})
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
				content = parsers.ParseEmbed(msg.Embeds[0])
			default:
				content = "[no text content]"
			}
		}

		if !containsKeyword(content, keywords) {
			continue
		}

		messageURL := fmt.Sprintf("https://discord.com/channels/%s/%s/%s", guildID, channelID, msg.ID)
		entries = append(entries, internal.AnalysisEntry{
			Text:       fmt.Sprintf("[%s]: %s", msg.Author.Username, content),
			MessageURL: messageURL,
		})
	}

	return entries
}

func analyzeForumChannel(s *discordgo.Session, guildID string, forum *discordgo.Channel, limit int, keywords []string) []internal.AnalysisEntry {
	var entries []internal.AnalysisEntry

	channelName := "#" + forum.Name
	forumThreads, err := fetchForumThreads(s, forum.ID, limit)
	if err != nil {
		return append(entries, internal.AnalysisEntry{
			Text: fmt.Sprintf("Failed to read forum %s: %v", channelName, err),
		})
	}

	entries = append(entries, internal.AnalysisEntry{
		Text: fmt.Sprintf("%s (forum): fetched %d posts", channelName, len(forumThreads)),
	})

	for _, thread := range forumThreads {
		content := parsers.ParseForumPost(s, thread)
		if !containsKeyword(content, keywords) {
			continue
		}
		threadURL := fmt.Sprintf("https://discord.com/channels/%s/%s", guildID, thread.ID)
		entries = append(entries, internal.AnalysisEntry{
			Text:       content,
			MessageURL: threadURL,
		})
	}

	return entries
}

func containsKeyword(text string, keywords []string) bool {
	if len(keywords) == 0 {
		return true
	}
	lower := strings.ToLower(text)
	for _, kw := range keywords {
		if strings.Contains(lower, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}

// fetchForumThreads uses the user-facing thread search endpoint
// (works with both user and bot tokens, unlike GuildThreadsActive).
func fetchForumThreads(s *discordgo.Session, forumChannelID string, limit int) ([]*discordgo.Channel, error) {
	if limit > 25 {
		limit = 25
	}
	endpoint := discordgo.EndpointChannel(forumChannelID) +
		"/threads/search?sort_by=last_message_time&sort_order=desc&limit=" +
		strconv.Itoa(limit)

	body, err := s.Request("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Threads []*discordgo.Channel `json:"threads"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse forum threads: %w", err)
	}

	return resp.Threads, nil
}
