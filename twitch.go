package twitchgo

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/kesuaheli/twitchgo/oauth"
)

const (
	baseURL = "https://api.twitch.tv/helix"
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

// Session is the instance for all Twitch events.
type Session struct {
	mu sync.Mutex

	clientID     string
	clientSecret string
	oauth        *oauth.Client

	ircToken string
	username string
	ircConn  net.Conn
	events   map[IRCMessageCommandName][]interface{}
	eventMu  sync.Mutex
	Prefix   string
}

// New creates a new Twitch instance for API and IRC connections. Can be used to register event
// responses before connecting.
//
// Setting clientID or clientSecret to an empty string will result in not connecting to the Twitch
// API on the call to s.Connect. Setting ircToken to an empty string will result in not connecting
// to the IRC server on the call to s.Connect.
//
// See also [NewAPIOnly], [NewIRCOnly] to create sessions for only one connection type and
// [Session.SetAPI], [Session.SetIRC] to update an existing Twitch session.
func New(clientID, clientSecret, ircToken string) *Session {
	return NewAPIOnly(clientID, clientSecret).SetIRC(ircToken)
}

// NewAPIOnly creates a new Twitch instance only for API connection. Can be used to register event
// responses before connecting. Setting clientID or clientSecret to an empty string will result in
// not connecting to the Twitch API on the call to s.Connect.
//
// See also [New], [NewIRCOnly] to create other session types and [Session.SetAPI],
// [Session.SetIRC] to update an existing Twitch session.
func NewAPIOnly(clientID, clientSecret string) *Session {
	return (&Session{}).SetAPI(clientID, clientSecret)
}

// NewIRCOnly creates a new Twitch instance only for IRC connection. Can be used to register event
// responses before connecting. Setting ircToken to an empty string will result in not connecting to
// the IRC server on the call to s.Connect.
//
// See also [New], [NewAPIOnly] to create other session types and [Session.SetAPI],
// [Session.SetIRC] to update an existing Twitch session.
func NewIRCOnly(ircToken string) *Session {
	return (&Session{}).SetIRC(ircToken)
}

// SetAPI sets the credentials used for a connection to the Twitch API. It will override the
// existing credentials, if previously set. Setting clientID or clientSecret to an empty string will
// result in not connecting to the Twitch API on the call to s.Connect. SetAPI will panic when
// s.Connect was already successfull.
//
// See also [Session.SetIRC] to update IRC credentials.
func (s *Session) SetAPI(clientID, clientSecret string) *Session {
	if s.oauth != nil || s.ircConn != nil {
		panic("Session already connected")
	}

	s.clientID = clientID
	s.clientSecret = clientSecret

	s.oauth = oauth.New(
		"https://id.twitch.tv/oauth2/token",
		clientID,
		clientSecret,
		"",
	)

	return s
}

// SetIRC sets the token used for a connection to the Twitch IRC server. It will override the
// existing token, if previously set. Setting ircToken to an empty string will result in not
// connecting to the IRC server on the call to s.Connect. SetIRC will panic when s.Connect was
// already successfull.
//
// See also [Session.SetAPI] to update API credentials.
func (s *Session) SetIRC(ircToken string) *Session {
	if s.oauth != nil || s.ircConn != nil {
		panic("Session already connected")
	}

	s.ircToken = ircToken
	s.username = ""
	return s
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

// Connect actually starts the connection to twitch.
func (s *Session) Connect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ircConn != nil {
		return ErrAlreadyConnected
	}

	var err error
	adress := fmt.Sprintf("%s:%d", IRCHost, IRCPort)
	s.ircConn, err = net.Dial("tcp", adress)
	if err != nil {
		log.Printf("Dial failed: %+v", err)
		return err
	}

	s.SendCommand("CAP REQ :twitch.tv/commands twitch.tv/membership twitch.tv/tags")
	s.SendCommandf("PASS %s", s.ircToken)
	s.SendCommandf("NICK %s", s.username)

	if err = s.waitForInit(); err != nil {
		s.ircConn.Close()
		return err
	}

	go s.listen()
	return nil
}

func (s *Session) waitForInit() (err error) {
	s.ircConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	var checklist byte
done:
	for {
		var buf []byte
		buf, err = s.readAll()
		if err != nil {
			break
		}
		msgs := strings.Split(string(buf), "\r\n")
		for _, raw := range msgs {
			var check byte
			check, err = s.parseInitMessage(raw)
			checklist |= check
			if err != nil || checklist == 255 {
				break done
			}
		}
	}
	s.ircConn.SetReadDeadline(time.Time{})
	return err
}

func (s *Session) parseInitMessage(raw string) (byte, error) {
	m := s.parseMessage(raw)
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
		m.handle(s)
		return 128, nil
	default:
		if m.Command.Name == IRCMsgCmdNotice && m.Command.Data == "Improperly formatted auth" {
			return 0, ErrInvalidToken
		}
		m.handle(s)
		return 0, nil
	}
}

func (s *Session) Close() {
	s.ircConn.Close()
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
