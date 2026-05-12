package tasks

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/seekehr/DiscordFreelanceX/internal"
	"github.com/seekehr/DiscordFreelanceX/internal/utils"
)

const (
	lastMsgPath     = "resources/last_msg.txt"
	cooldownHours   = 6
	channelInterval = 5 * time.Second
)

// SendMessages sends cfg.Bot.Message to every configured channel across all guilds.
// It respects a 6-hour cooldown stored as a Unix timestamp in resources/last_msg.txt,
// waits if needed, then sends with a 5-second interval between channels.
// Returns the total number of messages successfully sent.
func SendMessages(s *discordgo.Session, cfg *internal.Config) int {
	waitForCooldown()

	sent := 0
	channels := 0
	guilds := 0

	for _, server := range cfg.Servers {
		guildName := utils.GetGuildNameFromID(s, server.GuildID)
		guildSent := 0

		for i, channelID := range server.ChannelIDs {
			if i > 0 {
				time.Sleep(channelInterval)
			}

			channelName := utils.GetChannelNameFromID(s, channelID)
			_, err := s.ChannelMessageSend(channelID, cfg.Bot.Message)
			if err != nil {
				fmt.Printf("  Failed to send in %s / %s: %v\n", guildName, channelName, err)
				continue
			}

			fmt.Printf("  Sent in %s / %s\n", guildName, channelName)
			guildSent++
			channels++
		}

		if guildSent > 0 {
			guilds++
		}
		sent += guildSent
	}

	fmt.Printf("Sent %d messages in %d channels across %d guilds\n", sent, channels, guilds)
	updateTimestamp()
	return sent
}

// waitForCooldown reads the last-sent timestamp from disk and blocks until
// the 6-hour cooldown has elapsed. If the file is missing or unparseable
// the cooldown is considered expired.
func waitForCooldown() {
	data, err := os.ReadFile(lastMsgPath)
	if err != nil {
		return
	}

	ts, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return
	}

	lastSent := time.Unix(ts, 0)
	elapsed := time.Since(lastSent)
	remaining := (cooldownHours * time.Hour) - elapsed

	if remaining > 0 {
		fmt.Printf("Cooldown active — waiting %s before sending messages\n", remaining.Round(time.Second))
		time.Sleep(remaining)
	}
}

// updateTimestamp writes the current Unix timestamp to resources/last_msg.txt.
func updateTimestamp() {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	if err := os.WriteFile(lastMsgPath, []byte(ts), 0644); err != nil {
		fmt.Printf("Failed to update last_msg.txt: %v\n", err)
	}
}
