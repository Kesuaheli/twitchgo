package twitchgo

import (
	"errors"
	"io"
	"net"
	"strings"
)

func (s *Session) listen() {
	for {
		buf, err := s.readAll()
		if errors.Is(err, net.ErrClosed) {
			break
		} else if err != nil {
			break
		}
		msgs := strings.Split(string(buf), "\r\n")
		for _, m := range msgs {
			s.parseMessage(m).handle(s)
		}
	}
}

func (s *Session) readAll() ([]byte, error) {
	buf := make([]byte, 0)
	b := make([]byte, 1024)
	for {
		n, err := s.ircConn.Read(b)
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

func (s *Session) parseMessage(raw string) *IRCMessage {
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
