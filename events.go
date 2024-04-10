package twitchgo

import (
	"strings"
)

var ircCallbackEventMap = make(map[IRCMessageCommandName]func(t *IRCSession, m *IRCMessage, c interface{}))

// OnChannelJoin tells the bot to call the given callback function when a user joins a channel that
// you (the bot) already joined.
func (t *IRCSession) OnChannelJoin(callback IRCChannelJoinCallback) {
	t.events[IRCMsgCmdJoin] = append(t.events[IRCMsgCmdJoin], &callback)
}

func (t *IRCSession) OnChannelNotice(callback func(t *IRCSession)) {
	t.events[IRCMsgCmdNotice] = append(t.events[IRCMsgCmdNotice], &callback)
}

// OnChannelLeave tells the bot to call the given callback function when a user diconnects from a
// channel that you (the bot) already joined.
func (t *IRCSession) OnChannelLeave(callback IRCChannelLeaveCallback) {
	t.events[IRCMsgCmdPart] = append(t.events[IRCMsgCmdPart], &callback)
}

// OnChannelMessage tells the bot to call the given callback function when someone sends a message
// in a channel that you (the bot) already joined.
func (t *IRCSession) OnChannelMessage(callback IRCChannelMessageCallback) {
	t.events[IRCMsgCmdPrivmsg] = append(t.events[IRCMsgCmdPrivmsg], &callback)
}

// OnGlobalUserState is called right after the bot has connected successfully. So this callback
// function is only useful when adding Before calling Connect().
//
// The tags are all the tags the bot user has globally on twitch. Such as display name, global
// badges, user ID,... See "https://dev.twitch.tv/docs/irc/tags/#globaluserstate-tags" for a list of
// all possible tags.
func (t *IRCSession) OnGlobalUserState(callback IRCGlobalUserStateCallback) {
	t.events[IRCMsgCmdGlobaluserstate] = append(t.events[IRCMsgCmdGlobaluserstate], &callback)
}

// OnRoomState is called right after the bot has connected successfully.
func (t *IRCSession) OnRoomState(callback IRCRoomStateCallback) {
	t.events[IRCMsgCmdRoomstate] = append(t.events[IRCMsgCmdRoomstate], &callback)
}

// OnChannelCommandMessage is similar to OnChannelMessage.
//
// OnChannelCommandMessage tells the bot to call the given callback function when someone sends a
// command in a channel that you (the bot) already joined.
// A command is defined by a prefix (usually "!"), e.g. the message "!foo bar" translates to the
// command "foo" with the argument "bar".
func (t *IRCSession) OnChannelCommandMessage(cmd string, ignoreCase bool, callback IRCChannelCommandMessageCallback) {
	if ignoreCase {
		cmd = strings.ToLower(cmd)
	}
	t.OnChannelMessage(func(t *IRCSession, channel string, source *IRCUser, msg string) {
		args := strings.Split(msg, " ")
		msgCommand := args[0]

		msgCommand, hasPrefix := strings.CutPrefix(msgCommand, t.Prefix)
		if !hasPrefix {
			return
		}

		if ignoreCase {
			msgCommand = strings.ToLower(msgCommand)
		}
		if msgCommand != cmd {
			return
		}

		callback(t, channel, source, args[1:])
	})
}

// OnAny is called on any event. This is usefull for debug purposes.
func (t *IRCSession) OnAny(callback IRCAnyCallback) {
	t.events["*"] = append(t.events["*"], &callback)
}

type IRCChannelJoinCallback func(t *IRCSession, channel string, source *IRCUser)
type IRCChannelLeaveCallback func(t *IRCSession, channel string, source *IRCUser)
type IRCChannelMessageCallback func(t *IRCSession, channel string, source *IRCUser, msg string)
type IRCChannelCommandMessageCallback func(t *IRCSession, channel string, source *IRCUser, args []string)
type IRCGlobalUserStateCallback func(t *IRCSession, userTags IRCMessageTags)
type IRCRoomStateCallback func(t *IRCSession, roomTags IRCMessageTags)

type IRCAnyCallback func(t *IRCSession, message IRCMessage)

func init() {
	ircCallbackEventMap[IRCMsgCmdJoin] = func(t *IRCSession, m *IRCMessage, c interface{}) {
		if f, ok := c.(*IRCChannelJoinCallback); ok {
			(*f)(t, m.Command.Arguments[0], m.Source)
		}
	}
	ircCallbackEventMap[IRCMsgCmdPart] = func(t *IRCSession, m *IRCMessage, c interface{}) {
		if f, ok := c.(*IRCChannelLeaveCallback); ok {
			(*f)(t, m.Command.Arguments[0], m.Source)
		}
	}
	ircCallbackEventMap[IRCMsgCmdPrivmsg] = func(t *IRCSession, m *IRCMessage, c interface{}) {
		if f, ok := c.(*IRCChannelMessageCallback); ok {
			(*f)(t, m.Command.Arguments[0], m.Source, m.Command.Data)
		}
	}
	ircCallbackEventMap[IRCMsgCmdGlobaluserstate] = func(t *IRCSession, m *IRCMessage, c interface{}) {
		if f, ok := c.(*IRCGlobalUserStateCallback); ok {
			(*f)(t, m.Tags)
		}
	}
	ircCallbackEventMap[IRCMsgCmdRoomstate] = func(t *IRCSession, m *IRCMessage, c interface{}) {
		if f, ok := c.(*IRCRoomStateCallback); ok {
			(*f)(t, m.Tags)
		}
	}

	// on any
	ircCallbackEventMap["*"] = func(t *IRCSession, m *IRCMessage, c interface{}) {
		if f, ok := c.(*IRCAnyCallback); ok {
			(*f)(t, *m)
		}
	}
}
