package twitchgo

import (
	"fmt"
	"log"
	"strings"
)

// SendCommandf formats according to a format specifier and sends the resulting command to twitch
func (t *Twitch) SendCommandf(format string, a ...any) {
	t.SendCommand(fmt.Sprintf(format, a...))
}

// SendCommand sends the given command to twitch
func (t *Twitch) SendCommand(cmd string) {
	cmd = strings.TrimSuffix(cmd, "\n") + "\r\n"
	if len(cmd) == 2 {
		return
	}
	_, err := t.conn.Write([]byte(cmd))
	if err != nil {
		log.Printf("failed to send command '%s': %+v", cmd, err)
		return
	}
	if !strings.HasPrefix(cmd, string(MsgCmdPass)) {
		log.Printf("<< %s", cmd)
	} else {
		log.Printf("<< %s ***", MsgCmdPass)
	}
}

// SendMessagef formats according to a format specifier and sends the resulting message to the given
// channel
func (t *Twitch) SendMessagef(channel, format string, a ...any) {
	t.SendMessage(channel, fmt.Sprintf(format, a...))
}

// SendMessage sends a message to the given channel
func (t *Twitch) SendMessage(channel, msg string) {
	channel, _ = strings.CutPrefix(channel, "#")
	t.SendCommandf("%s #%s :%s", MsgCmdPrivmsg, channel, msg)
}

// JoinChannel joins the given channel and receives messages from that channel afterwards
func (t *Twitch) JoinChannel(channel string) {
	channel, _ = strings.CutPrefix(channel, "#")
	t.SendCommandf("%s #%s", MsgCmdJoin, channel)
}
