package twitchgo

import (
	"fmt"
	"log"
	"strings"
)

// SendCommandf formats according to a format specifier and sends the resulting command to twitch
func (t *IRCSession) SendCommandf(format string, a ...any) {
	t.SendCommand(fmt.Sprintf(format, a...))
}

// SendCommand sends the given command to twitch
func (t *IRCSession) SendCommand(cmd string) {
	cmd = strings.TrimSuffix(cmd, "\n") + "\r\n"
	if len(cmd) == 2 {
		return
	}
	_, err := t.conn.Write([]byte(cmd))
	if err != nil {
		log.Printf("failed to send command '%s': %+v", cmd, err)
		return
	}
	if !strings.HasPrefix(cmd, string(IRCMsgCmdPass)) {
		log.Printf("<< %s", cmd)
	} else {
		log.Printf("<< %s ***", IRCMsgCmdPass)
	}
}

// SendMessagef formats according to a format specifier and sends the resulting message to the given
// channel
func (t *IRCSession) SendMessagef(channel, format string, a ...any) {
	t.SendMessage(channel, fmt.Sprintf(format, a...))
}

// SendMessage sends a message to the given channel
func (t *IRCSession) SendMessage(channel, msg string) {
	channel, _ = strings.CutPrefix(channel, "#")
	t.SendCommandf("%s #%s :%s", IRCMsgCmdPrivmsg, channel, msg)
}

// JoinChannel joins the given channel and receives messages from that channel afterwards
func (t *IRCSession) JoinChannel(channel string) {
	channel, _ = strings.CutPrefix(channel, "#")
	t.SendCommandf("%s #%s", IRCMsgCmdJoin, channel)
}

// LeaveChannel leaves the given channel and nolonger receives messages from that channel afterwards
func (t *IRCSession) LeaveChannel(channel string) {
	channel, _ = strings.CutPrefix(channel, "#")
	t.SendCommandf("%s #%s", IRCMsgCmdPart, channel)
}
