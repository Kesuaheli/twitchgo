package twitchgo

import (
	"strings"
)

var ircCallbackEventMap = make(map[IRCMessageCommandName]func(s *Session, m *IRCMessage, c interface{}))

// OnChannelJoin tells the bot to call the given callback function when a user joins a channel that
// you (the bot) already joined.
func (s *Session) OnChannelJoin(callback IRCChannelJoinCallback) {
	s.events[IRCMsgCmdJoin] = append(s.events[IRCMsgCmdJoin], &callback)
}

func (s *Session) OnChannelNotice(callback func(s *Session)) {
	s.events[IRCMsgCmdNotice] = append(s.events[IRCMsgCmdNotice], &callback)
}

// OnChannelLeave tells the bot to call the given callback function when a user diconnects from a
// channel that you (the bot) already joined.
func (s *Session) OnChannelLeave(callback IRCChannelLeaveCallback) {
	s.events[IRCMsgCmdPart] = append(s.events[IRCMsgCmdPart], &callback)
}

// OnChannelMessage tells the bot to call the given callback function when someone sends a message
// in a channel that you (the bot) already joined.
func (s *Session) OnChannelMessage(callback IRCChannelMessageCallback) {
	s.events[IRCMsgCmdPrivmsg] = append(s.events[IRCMsgCmdPrivmsg], &callback)
}

// OnGlobalUserState is called right after the bot has connected successfully. So this callback
// function is only useful when adding Before calling Connect().
//
// The tags are all the tags the bot user has globally on twitch. Such as display name, global
// badges, user ID,... See "https://dev.twitch.tv/docs/irc/tags/#globaluserstate-tags" for a list of
// all possible tags.
func (s *Session) OnGlobalUserState(callback IRCGlobalUserStateCallback) {
	s.events[IRCMsgCmdGlobaluserstate] = append(s.events[IRCMsgCmdGlobaluserstate], &callback)
}

// OnRoomState is called right after the bot has connected successfully.
func (s *Session) OnRoomState(callback IRCRoomStateCallback) {
	s.events[IRCMsgCmdRoomstate] = append(s.events[IRCMsgCmdRoomstate], &callback)
}

// OnChannelCommandMessage is similar to OnChannelMessage.
//
// OnChannelCommandMessage tells the bot to call the given callback function when someone sends a
// command in a channel that you (the bot) already joined.
// A command is defined by a prefix (usually "!"), e.g. the message "!foo bar" translates to the
// command "foo" with the argument "bar".
func (s *Session) OnChannelCommandMessage(cmd string, ignoreCase bool, callback IRCChannelCommandMessageCallback) {
	if ignoreCase {
		cmd = strings.ToLower(cmd)
	}
	s.OnChannelMessage(func(s *Session, channel string, source *IRCUser, msg, msgID string, tags IRCMessageTags) {
		args := strings.Split(msg, " ")
		msgCommand := args[0]

		msgCommand, hasPrefix := strings.CutPrefix(msgCommand, s.Prefix)
		if !hasPrefix {
			return
		}

		if ignoreCase {
			msgCommand = strings.ToLower(msgCommand)
		}
		if msgCommand != cmd {
			return
		}

		callback(s, channel, source, args[1:])
	})
}

// OnAny is called on any event. This is usefull for debug purposes.
func (s *Session) OnAny(callback IRCAnyCallback) {
	s.events["*"] = append(s.events["*"], &callback)
}

type IRCChannelJoinCallback func(s *Session, channel string, source *IRCUser)
type IRCChannelLeaveCallback func(s *Session, channel string, source *IRCUser)
type IRCChannelMessageCallback func(s *Session, channel string, source *IRCUser, msg, msgID string, tags IRCMessageTags)
type IRCChannelCommandMessageCallback func(s *Session, channel string, source *IRCUser, args []string)
type IRCGlobalUserStateCallback func(s *Session, userTags IRCMessageTags)
type IRCRoomStateCallback func(s *Session, roomTags IRCMessageTags)

type IRCAnyCallback func(s *Session, message IRCMessage)

func init() {
	ircCallbackEventMap[IRCMsgCmdJoin] = func(s *Session, m *IRCMessage, c interface{}) {
		if f, ok := c.(*IRCChannelJoinCallback); ok {
			(*f)(s, m.Command.Arguments[0], m.Source)
		}
	}
	ircCallbackEventMap[IRCMsgCmdPart] = func(s *Session, m *IRCMessage, c interface{}) {
		if f, ok := c.(*IRCChannelLeaveCallback); ok {
			(*f)(s, m.Command.Arguments[0], m.Source)
		}
	}
	ircCallbackEventMap[IRCMsgCmdPrivmsg] = func(s *Session, m *IRCMessage, c interface{}) {
		if f, ok := c.(*IRCChannelMessageCallback); ok {
			(*f)(s, m.Command.Arguments[0], m.Source, m.Command.Data, m.Tags.ID, m.Tags)
		}
	}
	ircCallbackEventMap[IRCMsgCmdGlobaluserstate] = func(s *Session, m *IRCMessage, c interface{}) {
		if f, ok := c.(*IRCGlobalUserStateCallback); ok {
			(*f)(s, m.Tags)
		}
	}
	ircCallbackEventMap[IRCMsgCmdRoomstate] = func(s *Session, m *IRCMessage, c interface{}) {
		if f, ok := c.(*IRCRoomStateCallback); ok {
			(*f)(s, m.Tags)
		}
	}

	// on any
	ircCallbackEventMap["*"] = func(s *Session, m *IRCMessage, c interface{}) {
		if f, ok := c.(*IRCAnyCallback); ok {
			(*f)(s, *m)
		}
	}
}
