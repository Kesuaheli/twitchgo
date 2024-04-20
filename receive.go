package twitchgo

import (
	"errors"
	"io"
	"net"
	"strings"
	"time"
)

// waitForInit waits up to 5 seconds for a login response from the Twitch IRC server.
func waitForInit(s *Session) (err error) {
	s.ircConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	defer s.ircConn.SetReadDeadline(time.Time{})

	for {
		var buf []byte
		buf, err = readAll(s.ircConn)
		if err != nil {
			return err
		}
		for _, raw := range strings.Split(string(buf), "\r\n") {
			m := parseMessage(raw)
			if m.Command.Name == IRCMsgCmdGlobaluserstate {
				return nil
			} else if m.Command.Name == IRCMsgCmdNotice && m.Command.Data == "Improperly formatted auth" {
				return ErrInvalidToken
			}
		}
	}
}

func parseInitMessage(s *Session, raw string) (byte, error) {
	m := parseMessage(raw)
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

func listen(s *Session) {
	for {
		buf, err := readAll(s.ircConn)
		if errors.Is(err, net.ErrClosed) {
			break
		} else if err != nil {
			break
		}
		msgs := strings.Split(string(buf), "\r\n")
		for _, m := range msgs {
			parseMessage(m).handle(s)
		}
	}
}

func readAll(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 0)
	b := make([]byte, 1024)
	for {
		n, err := conn.Read(b)
		if err == io.EOF {
			break
		} else if err != nil {
			return []byte{}, err
		}
		buf = append(buf, b[:n]...)
		if buf[len(buf)-1] == '\n' {
			break
		}
	}
	return buf, nil
}

func parseMessage(raw string) *IRCMessage {
	if len(raw) == 0 {
		return nil
	}

	m := &IRCMessage{Raw: raw}

	if raw[0] == '@' {
		i := strings.Index(raw, " ")
		m.Tags = ParseRawIRCTags(raw[1:i])
		raw = raw[i+1:]
	}

	if raw[0] == ':' {
		i := strings.Index(raw, " ")
		source := strings.Split(raw[1:i], "!")
		if len(source) == 2 {
			m.Source = &IRCUser{Nickname: source[0], Host: source[0]}
		} else {
			m.Source = &IRCUser{Host: source[0]}
		}
		raw = raw[i+1:]
	}

	data := strings.Split(raw, " :")
	args := strings.Split(data[0], " ")

	m.Command.Name = IRCMessageCommandName(args[0])
	if len(args) > 1 {
		m.Command.Arguments = args[1:]
	}
	if len(data) > 1 {
		m.Command.Data = strings.Join(data[1:], " :")
	}

	return m
}

func (m *IRCMessage) handle(s *Session) {
	if m == nil || s == nil {
		return
	}

	// on ping commands only reply with a pong and exit the handler
	if m.Command.Name == IRCMsgCmdPing {
		s.SendCommand(string(IRCMsgCmdPong))
		return
	}

	handleCallback := ircCallbackEventMap[m.Command.Name]
	if handleCallback == nil {
		return
	}
	for _, c := range s.events[m.Command.Name] {
		handleCallback(s, m, c)
	}

	handleCallback = ircCallbackEventMap["*"]
	if handleCallback == nil {
		return
	}
	for _, c := range s.events["*"] {
		handleCallback(s, m, c)
	}
}
