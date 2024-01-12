package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
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
	bot.OnChannelCommandMessage("hello", true, HandleCommandHello)

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

func HandleCommandHello(t *twitchgo.Twitch, channel string, u *twitchgo.User, args []string) {
	// Logging the message
	log.Printf("[%s] %s executed hello command with: \"%v\"", channel, u.Nickname, args)

	if len(args) == 0 {
		t.SendMessagef(channel, "Hello %s", u.Nickname)
	} else {
		t.SendMessagef(channel, "Hello %s", strings.Join(args, " "))
	}
}
