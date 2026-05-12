package tasks

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func ParseEmbed(embed *discordgo.MessageEmbed) string {
	var parts []string

	if embed.Author != nil && embed.Author.Name != "" {
		parts = append(parts, fmt.Sprintf("Author: %s", embed.Author.Name))
	}
	if embed.Title != "" {
		parts = append(parts, fmt.Sprintf("Title: %s", embed.Title))
	}
	if embed.Description != "" {
		parts = append(parts, fmt.Sprintf("Description: %s", embed.Description))
	}
	for _, field := range embed.Fields {
		parts = append(parts, fmt.Sprintf("%s: %s", field.Name, field.Value))
	}
	if embed.Footer != nil && embed.Footer.Text != "" {
		parts = append(parts, fmt.Sprintf("Footer: %s", embed.Footer.Text))
	}

	if len(parts) == 0 {
		return "[empty embed]"
	}

	return strings.Join(parts, " | ")
}
