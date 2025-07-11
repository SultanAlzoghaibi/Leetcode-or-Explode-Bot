package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Hello World")
	godotenv.Load()
	discToken := os.Getenv("DISCORD_TOKEN")

	fmt.Println("token:" + discToken)
	sess, err := discordgo.New("Bot " + discToken)

	if err != nil {
		fmt.Println("Error creating Discord session,", err)
	}
	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {

			return
		}

		if m.Content == "ping" {
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		}
	})

	sess.Identify.Intents = discordgo.IntentsGuildMessages

	err = sess.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
	}
	defer sess.Close()
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}
