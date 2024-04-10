package twitchgo

import (
	"errors"
	"io"
	"net"
	"strings"
)

func (t *IRCSession) listen() {
	for {
		buf, err := t.readAll()
		if errors.Is(err, net.ErrClosed) {
			break
		} else if err != nil {
			break
		}
		msgs := strings.Split(string(buf), "\r\n")
		for _, m := range msgs {
			t.parseMessage(m).handle(t)
		}
	}
}

func (t *IRCSession) readAll() ([]byte, error) {
	buf := make([]byte, 0)
	b := make([]byte, 1024)
	for {
		n, err := t.conn.Read(b)
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

func (t *IRCSession) parseMessage(raw string) *IRCMessage {
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

func (m *IRCMessage) handle(t *IRCSession) {
	if m == nil || t == nil {
		return
	}

	// on ping commands only reply with a pong and exit the handler
	if m.Command.Name == IRCMsgCmdPing {
		t.SendCommand(string(IRCMsgCmdPong))
		return
	}

	handleCallback := ircCallbackEventMap[m.Command.Name]
	if handleCallback == nil {
		return
	}
	for _, c := range t.events[m.Command.Name] {
		handleCallback(t, m, c)
	}

	handleCallback = ircCallbackEventMap["*"]
	if handleCallback == nil {
		return
	}
	for _, c := range t.events["*"] {
		handleCallback(t, m, c)
	}
}
