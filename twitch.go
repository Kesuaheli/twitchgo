package twitchgo

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

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

// SetAuthRefreshToken sets a custom refresh token to use for the API calls.
func (s *Session) SetAuthRefreshToken(refreshToken string) *Session {
	if s.oauth == nil {
		panic("Session has no API auth")
	}
	s.oauth.SetRefreshToken(refreshToken)
	return s
}

// SetIRC sets the token used for a connection to the Twitch IRC server. It will override the
// existing token, if previously set. Setting ircToken to an empty string will result in not
// connecting to the IRC server on the call to s.Connect. SetIRC will panic when s.Connect was
// already successfull.
//
// See also [Session.SetAPI] to update API credentials.
func (s *Session) SetIRC(ircToken string) *Session {
	if s.ircConn != nil {
		panic("Session already connected")
	}

	s.ircToken = ircToken
	s.events = make(map[IRCMessageCommandName][]interface{})
	s.Prefix = "!"
	return s
}

// Connect actually starts the connection to the Twitch IRC server.
func (s *Session) Connect() error {
	if s.ircToken == "" {
		return nil
	}

	if s.ircConn != nil {
		return ErrAlreadyConnected
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var err error
	address := fmt.Sprintf("%s:%d", IRCHost, IRCPort)
	s.ircConn, err = net.Dial("tcp", address)
	if err != nil {
		log.Printf("Dial failed: %+v", err)
		return err
	}

	s.SendCommand("CAP REQ :twitch.tv/commands twitch.tv/membership twitch.tv/tags")
	s.SendCommandf("PASS %s", s.ircToken)
	s.SendCommand("NICK -")

	if err = waitForInit(s); err != nil {
		s.ircConn.Close()
		return err
	}

	go listen(s)
	return nil
}

// Close closes the connection to the Twitch IRC server.
func (s *Session) Close() {
	if s.ircConn != nil {
		s.ircConn.Close()
	}
	log.Print("Twitch connection closed!")
}
