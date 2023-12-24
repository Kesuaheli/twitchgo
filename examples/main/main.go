package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kesuaheli/twitchgo"
)

const (
	username = ""
	token    = "" // Remember to never use your store your token in production code!
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// creating a new bot with credentials
	bot := twitchgo.New(username, token)

	// Adding event listeners
	bot.OnChannelMessage(ChannelMessage)

	// starting the connection
	err := bot.Connect()
	if err != nil {
		fmt.Printf("ERROR connecting Twitch bot: %v", err)
		os.Exit(1)
	}

	// joining channel (this only works if the bot is connected)
	bot.JoinChannel("kesuaheli")

	fmt.Print("\nPress Ctrl+C to exit\n\n")
	<-ctx.Done()
}

func ChannelMessage(t *twitchgo.Twitch, c string, u *twitchgo.User, m string) {
	// Logging the message e.g. the user "username" writes the message "message" in the chat
	// from user "channel":
	// "[#channel] <username> message"
	log.Printf("[%s] <%s> %s", c, u.Nickname, m)
}
