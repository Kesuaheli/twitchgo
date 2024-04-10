package twitchgo

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	IRCHost = "irc.chat.twitch.tv"
	IRCPort = 6667
)

var (
	// ErrAlreadyConnected is returned when t.Connect() is called while a connection is already
	//running.
	ErrAlreadyConnected = errors.New("already connected")

	// ErrInvalidToken is returned when the provided token or username is invalid or improperly
	// formatted and a connection could not be established.
	ErrInvalidToken = errors.New("invalid token")
)

type IRCSession struct {
	sync.Mutex

	token    string
	username string
	conn     net.Conn
	events   map[IRCMessageCommandName][]interface{}
	eventMu  sync.Mutex
	Prefix   string
}

type IRCMessage struct {
	Raw     string
	Tags    IRCMessageTags
	Source  *IRCUser
	Command IRCMessageCommand
}
type IRCUser struct {
	Nickname string
	Host     string
}

type IRCMessageCommand struct {
	Name      IRCMessageCommandName
	Arguments []string
	Data      string
}
type IRCMessageCommandName string

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

// NewIRC creates a new Twitch instance but doesn't actuallay do anything. Can be used to register
// event responses before connecting.
// Start a connection with t.Connect()
func NewIRC(username, token string) (t *IRCSession) {
	if username == "" {
		username = "-"
	}
	t = &IRCSession{
		token:    token,
		username: username,
		events:   make(map[IRCMessageCommandName][]interface{}),
		Prefix:   "!",
	}
	return t
}

// Connect actually starts the connection to twitch.
func (t *IRCSession) Connect() error {
	t.Lock()
	defer t.Unlock()

	if t.conn != nil {
		return ErrAlreadyConnected
	}

	var err error
	adress := fmt.Sprintf("%s:%d", IRCHost, IRCPort)
	t.conn, err = net.Dial("tcp", adress)
	if err != nil {
		log.Printf("Dial failed: %+v", err)
		return err
	}

	t.SendCommand("CAP REQ :twitch.tv/commands twitch.tv/membership twitch.tv/tags")
	t.SendCommandf("PASS %s", t.token)
	t.SendCommandf("NICK %s", t.username)

	if err = t.waitForInit(); err != nil {
		t.conn.Close()
		return err
	}

	go t.listen()
	return nil
}

func (t *IRCSession) waitForInit() (err error) {
	t.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	var checklist byte
done:
	for {
		var buf []byte
		buf, err = t.readAll()
		if err != nil {
			break
		}
		msgs := strings.Split(string(buf), "\r\n")
		for _, raw := range msgs {
			var check byte
			check, err = t.parseInitMessage(raw)
			checklist |= check
			if err != nil || checklist == 255 {
				break done
			}
		}
	}
	t.conn.SetReadDeadline(time.Time{})
	return err
}

func (t *IRCSession) parseInitMessage(raw string) (byte, error) {
	m := t.parseMessage(raw)
	if m == nil {
		return 0, nil
	}
	switch m.Command.Name {
	case IRCMsgCmdCap:
		return 1, nil
	case "001":
		return 2, nil
	case "002":
		return 4, nil
	case "003":
		return 8, nil
	case "004":
		return 16, nil
	case "375":
		return 32, nil
	case "372":
		return 64, nil
	case IRCMsgCmdGlobaluserstate:
		m.handle(t)
		return 128, nil
	default:
		if m.Command.Name == IRCMsgCmdNotice && m.Command.Data == "Improperly formatted auth" {
			return 0, ErrInvalidToken
		}
		m.handle(t)
		return 0, nil
	}
}

func (t *IRCSession) Close() {
	t.conn.Close()
	log.Print("Twitch connection closed!")
}

func (u *IRCUser) String() string {
	if u == nil {
		return "<nil User>"
	}
	if u.Nickname == "" {
		return u.Host
	}
	return u.Nickname
}
