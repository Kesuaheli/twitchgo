package twitchgo

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

type ChannelJoinCallback func(t *Twitch, channel string, source *User)
type ChannelLeaveCallback func(t *Twitch, channel string, source *User)
type ChannelMessageCallback func(t *Twitch, channel string, source *User, msg string)

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
}
