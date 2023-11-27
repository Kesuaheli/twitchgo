package twitchgo

import (
	"fmt"
	"log"
	"strings"
)

func (t *Twitch) SendCommandf(format string, a ...any) {
	t.SendCommand(fmt.Sprintf(format, a...))
}

func (t *Twitch) SendCommand(msg string) {
	if len(msg) == 0 {
		return
	}
	msg = fmt.Sprintf("%s\r\n", msg)
	_, err := t.conn.Write([]byte(msg))
	if err != nil {
		log.Printf("failed to send command '%s': %+v", msg, err)
		return
	}
	log.Printf("<< %s", msg)
}

func (t *Twitch) SendMessagef(channel, format string, a ...any) {
	t.SendMessage(channel, fmt.Sprintf(format, a...))
}

func (t *Twitch) SendMessage(channel, msg string) {
	channel, _ = strings.CutPrefix(channel, "#")
	t.SendCommandf("%s #%s :%s", MsgCmdPrivmsg, channel, msg)
}
