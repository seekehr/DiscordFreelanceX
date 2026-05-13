package main

import (
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"github.com/seekehr/DiscordFreelanceX/internal"
	"github.com/seekehr/DiscordFreelanceX/internal/gui"
	"github.com/seekehr/DiscordFreelanceX/internal/tasks"
	"github.com/seekehr/DiscordFreelanceX/internal/utils"
)

func main() {
	cfg, err := internal.LoadConfig("resources/config.yaml")
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}
	fmt.Println("Config loaded successfully")

	token := initaliseEnv()

	guiApp := gui.CreateApp()
	window, td := gui.CreateTabbedWindow(guiApp, cfg)

	go initialiseDiscord(cfg, token, guiApp, td)

	window.ShowAndRun()
}

// initaliseEnv loads the .env file and returns the Discord token.
// Panics if the file is missing or TOKEN is unset.
func initaliseEnv() string {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Sprintf("failed to load .env file: %v", err))
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		panic("TOKEN is not set in .env")
	}

	return token
}

// initialiseDiscord opens a Discord session in the background, and on ready
// fetches messages per guild tab and sets up real-time handlers on the "New" tab.
func initialiseDiscord(cfg *internal.Config, token string, a fyne.App, td *gui.TabbedDisplay) {
	discord, err := discordgo.New(token)
	if err != nil {
		gui.AppendAnalysisText(td.NewRT, fmt.Sprintf("Failed to create Discord session: %v", err))
		return
	}

	discord.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentMessageContent

	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		for _, server := range cfg.ReceiveServers {
			name := utils.GetGuildNameFromID(s, server.GuildID)
			td.RenameGuildTab(server.GuildID, name)
		}

		perGuild, err := tasks.AnalyzeLastMessages(cfg.Bot.AnalyzeLastXMessages, s, cfg)
		if err != nil {
			gui.AppendAnalysisText(td.NewRT, fmt.Sprintf("Failed to analyze messages: %v", err))
			return
		}

		for guildID, entries := range perGuild {
			if rt, ok := td.GuildRTs[guildID]; ok {
				gui.AppendAnalysisEntries(rt, entries)
			}
		}

		gui.AppendAnalysisText(td.NewRT, "Listening for new messages...")
		tasks.AcceptNewMessages(s, cfg, a, td.NewRT)
	})

	if err := discord.Open(); err != nil {
		gui.AppendAnalysisText(td.NewRT, fmt.Sprintf("Failed to open Discord session: %v", err))
	}
	tasks.SendMessages(discord, cfg)
}
