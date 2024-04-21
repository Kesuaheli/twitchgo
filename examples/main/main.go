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
	ircToken = "" // Remember to never store your token in production code!
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// creating a new bot with credentials
	bot := twitchgo.NewIRCOnly(ircToken)

	// Adding event listeners
	bot.OnChannelMessage(ChannelMessage)

	// starting the connection
	err := bot.Connect()
	if err != nil {
		fmt.Printf("ERROR connecting Twitch bot: %v", err)
		os.Exit(1)
	}

	// joining channel (this only works if the bot is connected)
	bot.JoinChannel(username)

	fmt.Print("\nPress Ctrl+C to exit\n\n")
	<-ctx.Done()
}

func ChannelMessage(t *twitchgo.Session, c string, u *twitchgo.IRCUser, m string) {
	// Logging the message e.g. the user "username" writes the message "message" in the chat
	// from user "channel":
	// "[#channel] <username> message"
	log.Printf("[%s] <%s> %s", c, u.Nickname, m)
}
