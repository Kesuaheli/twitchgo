package twitchgo

import (
	"errors"
	"io"
	"net"
	"strings"
)

func (t *Twitch) listen() {
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

func (t *Twitch) readAll() ([]byte, error) {
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

func (t *Twitch) parseMessage(raw string) *Message {
	if len(raw) == 0 {
		return nil
	}

	m := &Message{Raw: raw}

	if raw[0] == '@' {
		i := strings.Index(raw, " ")
		m.Tags = ParseRawTags(raw[1:i])
		raw = raw[i+1:]
	}

	if raw[0] == ':' {
		i := strings.Index(raw, " ")
		source := strings.Split(raw[1:i], "!")
		if len(source) == 2 {
			m.Source = &User{Nickname: source[0], Host: source[0]}
		} else {
			m.Source = &User{Host: source[0]}
		}
		raw = raw[i+1:]
	}

	data := strings.Split(raw, " :")
	args := strings.Split(data[0], " ")

	m.Command.Name = MessageCommandName(args[0])
	if len(args) > 1 {
		m.Command.Arguments = args[1:]
	}
	if len(data) > 1 {
		m.Command.Data = strings.Join(data[1:], " :")
	}

	return m
}

func (m *Message) handle(t *Twitch) {
	if m == nil || t == nil {
		return
	}

	// on ping commands only reply with a pong and exit the handler
	if m.Command.Name == MsgCmdPing {
		t.SendCommand(string(MsgCmdPong))
		return
	}

	handleCallback := callbackEventMap[m.Command.Name]
	if handleCallback == nil {
		return
	}
	for _, c := range t.events[m.Command.Name] {
		handleCallback(t, m, c)
	}

	handleCallback = callbackEventMap["*"]
	if handleCallback == nil {
		return
	}
	for _, c := range t.events["*"] {
		handleCallback(t, m, c)
	}
}
