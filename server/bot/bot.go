package bot

import (
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

var botSession *discordgo.Session

func Start() error {
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		return fmt.Errorf("ERROR: could not start bot using provided token: %v", err)
	}
	botSession = dg

	err = dg.Open()
	if err != nil {
		return fmt.Errorf("ERROR: could not connect to Discord: %v", err)
	}

	go HandleMessages()

	return nil
}

func HandleMessages() {
	for range time.Tick(time.Second) {
		if botSession != nil {
			messages, err := botSession.ChannelMessages(os.Getenv("DISCORD_CHANNEL_ID"), 1, "", "", "")
			if err != nil {
				fmt.Println("ERROR: could not fetch message from client:", err)
			}
			if messages[0].Content != "" {
				fmt.Println("LOG: new message received: ", messages[0].Content)
			}
		}
	}
}

func SendMessage(message string) error {
	if botSession == nil {
		return fmt.Errorf("ERROR: bot hasn't been initialized yet")
	}

	_, err := botSession.ChannelMessageSend(os.Getenv("DISCORD_CHANNEL_ID"), message)
	if err != nil {
		return fmt.Errorf("ERROR: could not send message to Discord: %v", err)
	}

	return nil
}
