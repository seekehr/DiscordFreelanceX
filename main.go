package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"github.com/seekehr/DiscordFreelanceX/internal"
)

func main() {
	cfg, err := internal.LoadConfig("resources/config.yaml")
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}
	fmt.Println("Config loaded successfully")
	_ = cfg

	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env file: ", err)
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN is not set in .env")
	}

	fmt.Println("Token loaded successfully")
	discord, err := discordgo.New(token)
	if err != nil {
		log.Fatal("failed to create Discord session: ", err)
	}

	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("Bot is running!")
	})

	if err := discord.Open(); err != nil {
		log.Fatal("failed to open Discord session: ", err)
	}
	defer discord.Close()

	fmt.Println("Press Ctrl+C to exit")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	fmt.Println("Shutting down...")
}
