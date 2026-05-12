package main

import (
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2/widget"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"github.com/seekehr/DiscordFreelanceX/internal"
	"github.com/seekehr/DiscordFreelanceX/internal/gui"
	"github.com/seekehr/DiscordFreelanceX/internal/tasks"
)

func main() {
	cfg, err := internal.LoadConfig("resources/config.yaml")
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}
	fmt.Println("Config loaded successfully")

	token := initaliseEnv()

	guiApp := gui.CreateApp()
	window, entry := gui.CreateAnalysisWindow(guiApp)

	go initialiseDiscord(cfg, token, entry)

	window.ShowAndRun()
}

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

func initialiseDiscord(cfg *internal.Config, token string, entry *widget.Entry) {
	discord, err := discordgo.New(token)
	if err != nil {
		gui.AppendAnalysisText(entry, fmt.Sprintf("Failed to create Discord session: %v", err))
		return
	}

	discord.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentMessageContent

	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		gui.AppendAnalysisText(entry, "Bot is running! Fetching messages...")
		text, err := tasks.AnalyzeLastMessages(cfg.Bot.AnalyzeLastXMessages, s, cfg)
		if err != nil {
			gui.AppendAnalysisText(entry, fmt.Sprintf("Failed to analyze messages: %v", err))
			return
		}
		gui.AppendAnalysisText(entry, text)
	})

	if err := discord.Open(); err != nil {
		gui.AppendAnalysisText(entry, fmt.Sprintf("Failed to open Discord session: %v", err))
	}
}
