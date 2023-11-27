package twitchgo

type CallbackFunc func(*Twitch, *Message)

func (t *Twitch) SetEventChannelJoin(callback CallbackFunc) {
	c := t.events[MsgCmdJoin]
	c.msg = callback
	t.events[MsgCmdJoin] = c
}

func (t *Twitch) SetEventChannelNotice(callback CallbackFunc) {
	c := t.events[MsgCmdNotice]
	c.msg = callback
	t.events[MsgCmdNotice] = c
}

func (t *Twitch) SetEventChannelLeave(callback CallbackFunc) {
	c := t.events[MsgCmdPart]
	c.msg = callback
	t.events[MsgCmdPart] = c
}

func (t *Twitch) SetEventChannelMessage(callback CallbackFunc) {
	c := t.events[MsgCmdPrivmsg]
	c.msg = callback
	t.events[MsgCmdPrivmsg] = c
}
