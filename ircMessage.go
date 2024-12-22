package twitchgo

// IRCMessage contains the basic data for a message from the IRC server.
type IRCMessage struct {
	Raw     string
	Tags    IRCMessageTags
	Source  *IRCUser
	Command IRCMessageCommand
}
