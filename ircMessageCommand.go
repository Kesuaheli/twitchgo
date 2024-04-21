package twitchgo

// IRCMessageCommand contains the actual command sent with the IRC message.
type IRCMessageCommand struct {
	Name      IRCMessageCommandName
	Arguments []string
	Data      string
}

// IRCMessageCommandName is the type for the name which defines the command sent with the IRC
// message.
type IRCMessageCommandName string

// Some defined command names
const (
	// Your bot sends this message to join a channel.
	IRCMsgCmdJoin IRCMessageCommandName = "JOIN"
	// Your bot sends this message to specify the bot’s nickname when authenticating with the Twitch
	// IRC server.
	IRCMsgCmdNick IRCMessageCommandName = "NICK"
	// Your bot receives this message from the Twitch IRC server to indicate whether a command
	// succeeded or failed. For example, a moderator tried to ban a user that was already banned.
	IRCMsgCmdNotice IRCMessageCommandName = "NOTICE"
	// Your bot sends this message to leave a channel.
	//
	// Your bot receives this message from the Twitch IRC server when a channel bans it.
	IRCMsgCmdPart IRCMessageCommandName = "PART"
	// Your bot sends this message to specify the bot’s password when authenticating with the Twitch
	// IRC server.
	IRCMsgCmdPass IRCMessageCommandName = "PASS"
	// Your bot receives this message from the Twitch IRC server when the server wants to ensure
	// that your bot is still alive and able to respond to the server’s messages.
	IRCMsgCmdPing IRCMessageCommandName = "PING"
	// Your bot sends this message in reply to the Twitch IRC server’s PING message.
	IRCMsgCmdPong IRCMessageCommandName = "PONG"
	// Your bot sends this message to post a chat message in the channel’s chat room.
	//
	// Your bot receives this message from the Twitch IRC server when a user posts a chat message in
	// the chat room.
	IRCMsgCmdPrivmsg IRCMessageCommandName = "PRIVMSG"

	// Your bot receives this message from the Twitch IRC server when all messages are removed from
	// the chat room, or all messages for a specific user are removed from the chat room.
	IRCMsgCmdClearchat IRCMessageCommandName = "CLEARCHAT"
	// Your bot receives this message from the Twitch IRC server when a specific message is removed
	// from the chat room.
	IRCMsgCmdClearmsg IRCMessageCommandName = "CLEARMSG"
	// Your bot receives this message from the Twitch IRC server when a bot connects to the server.
	IRCMsgCmdGlobaluserstate IRCMessageCommandName = "GLOBALUSERSTATE"
	// Your bot receives this message from the Twitch IRC server when a channel starts or stops host
	// mode.
	IRCMsgCmdHosttarget IRCMessageCommandName = "HOSTTARGET"
	// Your bot receives this message from the Twitch IRC server when the server needs to perform
	// maintenance and is about to disconnect your bot.
	IRCMsgCmdReconnect IRCMessageCommandName = "RECONNECT"
	// Your bot receives this message from the Twitch IRC server when a bot joins a channel or a
	// moderator changes the chat room’s chat settings.
	IRCMsgCmdRoomstate IRCMessageCommandName = "ROOMSTATE"
	// Your bot receives this message from the Twitch IRC server when events like user subscriptions
	// occur.
	IRCMsgCmdUsernotice IRCMessageCommandName = "USERNOTICE"
	// Your bot receives this message from the Twitch IRC server when a user joins a channel or the
	// bot sends a PRIVMSG message.
	IRCMsgCmdUserstate IRCMessageCommandName = "USERSTATE"
	// Your bot receives this message from the Twitch IRC server when a user sends a WHISPER
	// message.
	IRCMsgCmdWhisper IRCMessageCommandName = "WHISPER"
	//
	IRCMsgCmdCap IRCMessageCommandName = "CAP"
	//
	IRCMsgCmdUserList IRCMessageCommandName = "353"
)
