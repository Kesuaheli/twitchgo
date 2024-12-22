package twitchgo

// IRCUser represents the source user of a IRCMessage.
type IRCUser struct {
	Nickname string
	Host     string
}

// String implements the [fmt.Stringer].
func (u IRCUser) String() string {
	if u.Nickname == "" {
		return u.Host
	}
	return u.Nickname
}
