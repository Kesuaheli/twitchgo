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

type Twitch struct {
	sync.Mutex

	token    string
	username string
	conn     net.Conn
	events   map[MessageCommandName]struct {
		msg func(t *Twitch, m *Message)
	}
	eventMu sync.Mutex
}

type Message struct {
	Raw     string
	Tags    MessageTags
	Source  string
	Command MessageCommand
}

type MessageCommand struct {
	Name      MessageCommandName
	Arguments []string
	Data      string
}
type MessageCommandName string

const (
	// Your bot sends this message to join a channel.
	MsgCmdJoin MessageCommandName = "JOIN"
	// Your bot sends this message to specify the bot’s nickname when authenticating with the Twitch
	// IRC server.
	MsgCmdNick MessageCommandName = "NICK"
	// Your bot receives this message from the Twitch IRC server to indicate whether a command
	// succeeded or failed. For example, a moderator tried to ban a user that was already banned.
	MsgCmdNotice MessageCommandName = "NOTICE"
	// Your bot sends this message to leave a channel.
	//
	// Your bot receives this message from the Twitch IRC server when a channel bans it.
	MsgCmdPart MessageCommandName = "PART"
	// Your bot sends this message to specify the bot’s password when authenticating with the Twitch
	// IRC server.
	MsgCmdPass MessageCommandName = "PASS"
	// Your bot receives this message from the Twitch IRC server when the server wants to ensure
	// that your bot is still alive and able to respond to the server’s messages.
	MsgCmdPing MessageCommandName = "PING"
	// Your bot sends this message in reply to the Twitch IRC server’s PING message.
	MsgCmdPong MessageCommandName = "PONG"
	// Your bot sends this message to post a chat message in the channel’s chat room.
	//
	// Your bot receives this message from the Twitch IRC server when a user posts a chat message in
	// the chat room.
	MsgCmdPrivmsg MessageCommandName = "PRIVMSG"

	// Your bot receives this message from the Twitch IRC server when all messages are removed from
	// the chat room, or all messages for a specific user are removed from the chat room.
	MsgCmdClearchat MessageCommandName = "CLEARCHAT"
	// Your bot receives this message from the Twitch IRC server when a specific message is removed
	// from the chat room.
	MsgCmdClearmsg MessageCommandName = "CLEARMSG"
	// Your bot receives this message from the Twitch IRC server when a bot connects to the server.
	MsgCmdGlobaluserstate MessageCommandName = "GLOBALUSERSTATE"
	// Your bot receives this message from the Twitch IRC server when a channel starts or stops host
	// mode.
	MsgCmdHosttarget MessageCommandName = "HOSTTARGET"
	// Your bot receives this message from the Twitch IRC server when the server needs to perform
	// maintenance and is about to disconnect your bot.
	MsgCmdReconnect MessageCommandName = "RECONNECT"
	// Your bot receives this message from the Twitch IRC server when a bot joins a channel or a
	// moderator changes the chat room’s chat settings.
	MsgCmdRoomstate MessageCommandName = "ROOMSTATE"
	// Your bot receives this message from the Twitch IRC server when events like user subscriptions
	// occur.
	MsgCmdUsernotice MessageCommandName = "USERNOTICE"
	// Your bot receives this message from the Twitch IRC server when a user joins a channel or the
	// bot sends a PRIVMSG message.
	MsgCmdUserstate MessageCommandName = "USERSTATE"
	// Your bot receives this message from the Twitch IRC server when a user sends a WHISPER
	// message.
	MsgCmdWhisper MessageCommandName = "WHISPER"
	//
	MsgCmdCap MessageCommandName = "CAP"
	//
	MsgCmdUserList MessageCommandName = "353"
)

// New creates a new Twitch instance but doesn't actuallay do anything. Can be used to register
// event responses before connecting.
// Start a connection with t.Connect()
func New(username, token string) (t *Twitch) {
	t = &Twitch{
		token:    token,
		username: username,
		events:   make(map[MessageCommandName]struct{ msg func(t *Twitch, m *Message) }),
	}
	return t
}

// Connect actually starts the connection to twitch.
func (t *Twitch) Connect() error {
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

func (t *Twitch) waitForInit() (err error) {
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

func (t *Twitch) parseInitMessage(raw string) (byte, error) {
	m := t.parseMessage(raw)
	if m == nil {
		return 0, nil
	}
	switch m.Command.Name {
	case MsgCmdCap:
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
	case "376":
		return 128, nil
	default:
		if m.Command.Name == MsgCmdNotice && m.Command.Data == "Improperly formatted auth" {
			return 0, ErrInvalidToken
		}
		m.handle(t)
		return 0, nil
	}
}

func (t *Twitch) Close() {
	t.conn.Close()
	log.Print("Twitch connection closed!")
}
