package twitchgo

import (
	"strings"
)

var callbackEventMap = make(map[MessageCommandName]func(t *Twitch, m *Message, c interface{}))

// OnChannelJoin tells the bot to call the given callback function when a user joins a channel that
// you (the bot) already joined.
func (t *Twitch) OnChannelJoin(callback ChannelJoinCallback) {
	t.events[MsgCmdJoin] = append(t.events[MsgCmdJoin], &callback)
}

func (t *Twitch) OnChannelNotice(callback func(t *Twitch)) {
	t.events[MsgCmdNotice] = append(t.events[MsgCmdNotice], &callback)
}

// OnChannelLeave tells the bot to call the given callback function when a user diconnects from a
// channel that you (the bot) already joined.
func (t *Twitch) OnChannelLeave(callback ChannelLeaveCallback) {
	t.events[MsgCmdPart] = append(t.events[MsgCmdPart], &callback)
}

// OnChannelMessage tells the bot to call the given callback function when someone sends a message
// in a channel that you (the bot) already joined.
func (t *Twitch) OnChannelMessage(callback ChannelMessageCallback) {
	t.events[MsgCmdPrivmsg] = append(t.events[MsgCmdPrivmsg], &callback)
}

// OnGlobalUserState is called right after the bot has connected successfully. So this callback
// function is only useful when adding Before calling Connect().
//
// The tags are all the tags the bot user has globally on twitch. Such as display name, global
// badges, user ID,... See "https://dev.twitch.tv/docs/irc/tags/#globaluserstate-tags" for a list of
// all possible tags.
func (t *Twitch) OnGlobalUserState(callback GlobalUserStateCallback) {
	t.events[MsgCmdGlobaluserstate] = append(t.events[MsgCmdGlobaluserstate], &callback)
}

// OnRoomState is called right after the bot has connected successfully.
func (t *Twitch) OnRoomState(callback RoomStateCallback) {
	t.events[MsgCmdRoomstate] = append(t.events[MsgCmdRoomstate], &callback)
}

// OnChannelCommandMessage is similar to OnChannelMessage.
//
// OnChannelCommandMessage tells the bot to call the given callback function when someone sends a
// command in a channel that you (the bot) already joined.
// A command is defined by a prefix (usually "!"), e.g. the message "!foo bar" translates to the
// command "foo" with the argument "bar".
func (t *Twitch) OnChannelCommandMessage(cmd string, callback ChannelCommandMessageCallback) {
	t.OnChannelMessage(func(t *Twitch, channel string, source *User, msg string) {
		var ok bool
		if msg, ok = strings.CutPrefix(msg, t.Prefix); !ok {
			return
		}
		if msg, ok = strings.CutPrefix(msg, cmd); !ok {
			return
		}

		msg = strings.Trim(msg, " ")
		args := []string{}
		if msg != "" {
			args = strings.Split(msg, " ")
		}
		callback(t, channel, source, args)
	})
}

// OnAny is called on any event. This is usefull for debug purposes.
func (t *Twitch) OnAny(callback AnyCallback) {
	t.events["*"] = append(t.events["*"], &callback)
}

type ChannelJoinCallback func(t *Twitch, channel string, source *User)
type ChannelLeaveCallback func(t *Twitch, channel string, source *User)
type ChannelMessageCallback func(t *Twitch, channel string, source *User, msg string)
type ChannelCommandMessageCallback func(t *Twitch, channel string, source *User, args []string)
type GlobalUserStateCallback func(t *Twitch, userTags MessageTags)
type RoomStateCallback func(t *Twitch, roomTags MessageTags)

type AnyCallback func(t *Twitch, message Message)

func init() {
	callbackEventMap[MsgCmdJoin] = func(t *Twitch, m *Message, c interface{}) {
		if f, ok := c.(*ChannelJoinCallback); ok {
			(*f)(t, m.Command.Arguments[0], m.Source)
		}
	}
	callbackEventMap[MsgCmdPart] = func(t *Twitch, m *Message, c interface{}) {
		if f, ok := c.(*ChannelLeaveCallback); ok {
			(*f)(t, m.Command.Arguments[0], m.Source)
		}
	}
	callbackEventMap[MsgCmdPrivmsg] = func(t *Twitch, m *Message, c interface{}) {
		if f, ok := c.(*ChannelMessageCallback); ok {
			(*f)(t, m.Command.Arguments[0], m.Source, m.Command.Data)
		}
	}
	callbackEventMap[MsgCmdGlobaluserstate] = func(t *Twitch, m *Message, c interface{}) {
		if f, ok := c.(*GlobalUserStateCallback); ok {
			(*f)(t, m.Tags)
		}
	}
	callbackEventMap[MsgCmdRoomstate] = func(t *Twitch, m *Message, c interface{}) {
		if f, ok := c.(*RoomStateCallback); ok {
			(*f)(t, m.Tags)
		}
	}

	// on any
	callbackEventMap["*"] = func(t *Twitch, m *Message, c interface{}) {
		if f, ok := c.(*AnyCallback); ok {
			(*f)(t, *m)
		}
	}
}
