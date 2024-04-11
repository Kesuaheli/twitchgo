package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kesuaheli/twitchgo"
)

const (
	token = "" // Remember to never use your store your token in production code!
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// creating a new bot with credentials
	bot := twitchgo.NewIRC("", token)

	// You can add a listener to the globaluserstate event, which is called right after a succesfull
	// connection. In this event you can read the display name of the bot (and more, of course).
	bot.OnGlobalUserState(GotGlobalUser)

	// starting the connection
	err := bot.Connect()
	if err != nil {
		fmt.Printf("ERROR connecting Twitch bot: %v", err)
		os.Exit(1)
	}

	fmt.Print("\nPress Ctrl+C to exit\n\n")
	<-ctx.Done()
}

func GotGlobalUser(t *twitchgo.IRCSession, tags twitchgo.IRCMessageTags) {
	fmt.Printf("User: %s\n", tags.DisplayName)
}
