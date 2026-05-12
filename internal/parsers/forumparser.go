package parsers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// ParseForumPost extracts readable info from a forum thread (post).
// It combines the thread title, applied tag names, and the starter message
// content into a single pipe-delimited string.
// Returns "[empty forum post]" when nothing useful is found.
func ParseForumPost(s *discordgo.Session, thread *discordgo.Channel) string {
	var parts []string

	if thread.Name != "" {
		parts = append(parts, fmt.Sprintf("Title: %s", thread.Name))
	}

	if len(thread.AppliedTags) > 0 {
		tagNames := resolveTagNames(s, thread.ParentID, thread.AppliedTags)
		if len(tagNames) > 0 {
			parts = append(parts, fmt.Sprintf("Tags: %s", strings.Join(tagNames, ", ")))
		}
	}

	msgs, err := s.ChannelMessages(thread.ID, 1, "", "", "")
	if err == nil && len(msgs) > 0 {
		starter := oldestMessage(msgs)
		if starter.Content != "" {
			parts = append(parts, fmt.Sprintf("Details: %s", starter.Content))
		}
		if len(starter.Embeds) > 0 {
			parts = append(parts, ParseEmbed(starter.Embeds[0]))
		}
	}

	if len(parts) == 0 {
		return "[empty forum post]"
	}

	return strings.Join(parts, " | ")
}

// resolveTagNames maps tag IDs to their display names by reading the
// parent forum channel's AvailableTags. Unresolvable IDs are kept as-is.
func resolveTagNames(s *discordgo.Session, forumChannelID string, tagIDs []string) []string {
	parent, err := s.Channel(forumChannelID)
	if err != nil {
		return tagIDs
	}

	lookup := make(map[string]string, len(parent.AvailableTags))
	for _, t := range parent.AvailableTags {
		lookup[t.ID] = t.Name
	}

	names := make([]string, 0, len(tagIDs))
	for _, id := range tagIDs {
		if name, ok := lookup[id]; ok {
			names = append(names, name)
		} else {
			names = append(names, id)
		}
	}
	return names
}

// oldestMessage returns the message with the smallest snowflake ID,
// which is the starter message of a forum post.
func oldestMessage(msgs []*discordgo.Message) *discordgo.Message {
	oldest := msgs[0]
	for _, m := range msgs[1:] {
		if m.ID < oldest.ID {
			oldest = m
		}
	}
	return oldest
}
