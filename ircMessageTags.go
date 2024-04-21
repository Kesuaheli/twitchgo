package twitchgo

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// IRCMessageTags https://dev.twitch.tv/docs/irc/tags/
type IRCMessageTags struct {
	// Optional. The message includes this tag if the user was put in a timeout. The tag contains
	// the duration of the timeout, in seconds.
	BanDuration int `json:"ban-duration"`
	// An ID that identifies the chat room (channel).
	RoomID string `json:"room-id"`
	// Optional. The ID of the user that was banned or put in a timeout. The user was banned if the
	// message doesn’t include the ban-duration tag.
	TargetUserID string `json:"target-user-id"`
	// The UNIX timestamp.
	Timestamp time.Time `json:"tmi-sent-ts"`

	// The name of the user who sent the message.
	Login string `json:"login"`
	// A UUID that identifies the message that was removed.
	TargetMsgID string `json:"target-msg-id"`

	// Contains metadata related to the chat badges in the badges tag.
	//
	// Currently, this tag contains metadata only for subscriber badges, to indicate the number of
	// months the user has been a subscriber.
	BadgeInfo []string `json:"badge-info"`
	// Comma-separated list of chat badges in the form, <badge>/<version>. For example, admin/1.
	// There are many possible badge values, but here are few:
	//	"admin"
	//	"bits"
	//	"broadcaster"
	//	"moderator"
	//	"subscriber"
	//	"staff"
	//	"turbo"
	// Most badges have only 1 version, but some badges like subscriber badges offer different
	// versions of the badge depending on how long the user has subscribed.
	//
	// To get the badge, use the Get Global Chat Badges and Get Channel Chat Badges APIs. Match the
	// badge to the set-id field’s value in the response. Then, match the version to the id field in
	// the list of versions.
	Badges []string `json:"badges"`
	// A comma-delimited list of IDs that identify the emote sets that the user has access to. Is
	// always set to at least zero ["0"]. To access the emotes in the set, use the Get Emote Sets
	// API.
	EmoteSets []string `json:"emote-sets"`

	// The user’s display name, escaped as described in the IRCv3 spec. This tag may be empty if it
	// is never set.
	DisplayName string `json:"display-name"`
	// The color of the user’s name in the chat room. This is a hexadecimal RGB color code in the
	// form, #<RGB>. This tag may be empty if it is never set.
	Color string `json:"color"`
	// The user’s ID.
	UserID string `json:"user-id"`
	// The type of user. Possible values are:
	//	"" // A normal user
	//	"admin" // A Twitch administrator
	//	"global_mod" // A global moderator
	//	"staff" // A Twitch employee
	UserType string `json:"user-type"`
	// A Boolean value that indicates whether the user has site-wide commercial free mode enabled.
	// Is true if enabled; otherwise, false.
	Turbo bool `json:"turbo"`
	// A Boolean value that determines whether the user is a subscriber. Is true if the user is a
	// subscriber; otherwise, false.
	Subscriber bool `json:"subscriber"`
	//A Boolean value that determines whether the user is a moderator. Is true if the user is a
	// moderator; otherwise, false.
	Mod bool `json:"mod"`
	// A Boolean value that determines whether the user that sent the chat is a VIP. The message
	// includes this tag if the user is a VIP; otherwise, the message doesn’t include this tag
	// (check for the presence of the tag instead of whether the tag is set to true or false).
	VIP bool `json:"vip"`

	// 	The amount of Bits the user cheered. Only a Bits cheer message includes this tag. To learn
	// more about Bits, see the Extensions Monetization Guide. To get the cheermote, use the Get
	// Cheermotes API. Match the cheer amount to the id field’s value in the response. Then, get the
	// cheermote’s URL based on the cheermote theme, type, and size you want to use.
	Bits int `json:"bits"`
	// The color of the user’s name in the chat room. This is a hexadecimal RGB color code in the
	// form, #<RGB>. This tag may be empty if it is never set.
	// A comma-delimited list of emotes and their positions in the message. Each emote is in the
	// form, <emote ID>:<start position>-<end position>. The position indices are zero-based.
	//
	// To get the actual emote, see the Get Channel Emotes and Get Global Emotes APIs. For
	// information about how to use the information that the APIs return, see Twitch emotes.
	//
	// NOTE It’s possible for the emotes flag’s value to be set to an action instead of identifying
	// an emote. For example, \001ACTION barfs on the floor.\001.
	Emotes []string `json:"emotes"`
	// An ID that uniquely identifies the message.
	//
	// If a privmsg was sent, an ID that uniquely identifies the message.
	ID string `json:"id"`

	// The value of the Hype Chat sent by the user.
	PinnedChatPaidAmount string `json:"pinned-chat-paid-amount"`
	// The ISO 4217 alphabetic currency code the user has sent the Hype Chat in.
	PinnedChatPaidCurrency string `json:"pinned-chat-paid-currency"`
	// Indicates how many decimal points this currency represents partial amounts in. Decimal points
	// start from the right side of the value defined in pinned-chat-paid-amount.
	PinnedChatPaidExponent string `json:"pinned-chat-paid-exponent"`
	// The level of the Hype Chat, in English. Possible values are:
	//	"ONE"
	//	"TWO"
	//	"THREE"
	//	"FOUR"
	//	"FIVE"
	//	"SIX"
	//	"SEVEN"
	//	"EIGHT"
	//	"NINE"
	//	"TEN"
	PinnedChatPaidLevel string `json:"pinned-chat-paid-level"`
	// A Boolean value that determines if the message sent with the Hype Chat was filled in by the
	// system.
	// If true, the user entered no message and the body message was automatically filled in by the
	// system.
	// If false, the user provided their own message to send with the Hype Chat.
	PinnedChatPaidIsSystemMessage bool `json:"pinned-chat-paid-is-system-message"`

	//	An ID that uniquely identifies the direct parent message that this message is replying to.
	// The message does not include this tag if this message is not a reply.
	ReplyParentMsgID string `json:"reply-parent-msg-id"`
	//	An ID that identifies the sender of the direct parent message. The message does not include
	// this tag if this message is not a reply.
	ReplyParentUserID string `json:"reply-parent-user-id"`
	//	The login name of the sender of the direct parent message. The message does not include this
	// tag if this message is not a reply.
	ReplyParentUserLogin string `json:"reply-parent-user-login"`
	//	The display name of the sender of the direct parent message. The message does not include
	// this tag if this message is not a reply.
	ReplyParentDisplayName string `json:"reply-parent-display-name"`
	//	The text of the direct parent message. The message does not include this tag if this message
	// is not a reply.
	ReplyParentMsgBody string `json:"reply-parent-msg-body"`
	//	An ID that uniquely identifies the top-level parent message of the reply thread that this
	// message is replying to. The message does not include this tag if this message is not a reply.
	ReplyThreadParentMsgID string `json:"reply-thread-parent-msg-id"`
	//	The login name of the sender of the top-level parent message. The message does not include
	// this tag if this message is not a reply.
	ReplyThreadParentUserLogin string `json:"reply-thread-parent-user-login"`

	ReplyThreadParentDisplayName string `json:"reply-thread-parent-display-name"`

	// A Boolean value that determines whether the chat room allows only messages with emotes. Is
	// true if only emotes are allowed; otherwise, false.
	EmoteOnly bool `json:"emote-only"`
	// An integer value that determines whether only followers can post messages in the chat room.
	// The value indicates how long, in minutes, the user must have followed the broadcaster before
	// posting chat messages. If the value is -1, the chat room is not restricted to followers only.
	FollowersOnly int `json:"followers-only"`
	//A Boolean value that determines whether a user’s messages must be unique. Applies only to
	// messages with more than 9 characters. Is true if users must post unique messages; otherwise,
	// false.
	R9K bool `json:"r9k"`
	//An integer value that determines how long, in seconds, users must wait between sending
	// messages.
	Slow int `json:"slow"`
	//A Boolean value that determines whether only subscribers and moderators can chat in the chat
	// room. Is true if only subscribers and moderators can chat; otherwise, false.
	SubsOnly bool `json:"subs-only"`

	//The type of notice (not the ID). Possible values are:
	//	"sub"
	//	"resub"
	//	"subgift"
	//	"submysterygift"
	//	"giftpaidupgrade"
	//	"rewardgift"
	//	"anongiftpaidupgrade"
	//	"raid"
	//	"unraid"
	//	"ritual"
	//	"bitsbadgetier"
	MsgType string `json:"msg-id"`
	// The message Twitch shows in the chat room for this notice.
	SystemMsg string `json:"system-msg"`

	// Included only with sub and resub notices.
	//
	// The total number of months the user has subscribed. This is the same as msg-param-months but
	// sent for different types of user notices.
	MsgParamCumulativeMonths string `json:"msg-param-cumulative-months"`
	// Included only with raid notices.
	//
	// The display name of the broadcaster raiding this channel.
	MsgParamDisplayName string `json:"msg-param-displayName"`
	// Included only with raid notices.
	//
	// The login name of the broadcaster raiding this channel.
	MsgParamLogin string `json:"msg-param-login"`
	// Included only with subgift notices.
	//
	// The total number of months the user has subscribed. This is the same as
	// msg-param-cumulative-months but sent for different types of user notices.
	MsgParamMonths string `json:"msg-param-months"`
	// Included only with anongiftpaidupgrade and giftpaidupgrade notices.
	//
	// The number of gifts the gifter has given during the promo indicated by msg-param-promo-name.
	MsgParamPromoGiftTotal string `json:"msg-param-promo-gift-total"`
	// Included only with anongiftpaidupgrade and giftpaidupgrade notices.
	//
	// The subscriptions promo, if any, that is ongoing (for example, Subtember 2018).
	MsgParamPromoName string `json:"msg-param-promo-name"`
	// Included only with subgift notices.
	//
	// The display name of the subscription gift recipient.
	MsgParamRecipientDisplayName string `json:"msg-param-recipient-display-name"`
	// Included only with subgift notices.
	//
	// The user ID of the subscription gift recipient.
	MsgParamRecipientID string `json:"msg-param-recipient-id"`
	// Included only with subgift notices.
	//
	// The user name of the subscription gift recipient.
	MsgParamRecipientUserName string `json:"msg-param-recipient-user-name"`
	// Included only with giftpaidupgrade notices.
	//
	// The login name of the user who gifted the subscription.
	MsgParamSenderLogin string `json:"msg-param-sender-login"`
	// Included only with giftpaidupgrade notices.
	//
	// The display name of the user who gifted the subscription.
	MsgParamSenderName string `json:"msg-param-sender-name"`
	// Included only with sub and resub notices.
	//
	// A Boolean value that indicates whether the user wants their streaks shared.
	MsgParamShouldShareStreak string `json:"msg-param-should-share-streak"`
	// Included only with sub and resub notices.
	//
	// The number of consecutive months the user has subscribed. This is zero (0) if
	// msg-param-should-share-streak is 0.
	MsgParamStreakMonths string `json:"msg-param-streak-months"`
	// Included only with sub, resub and subgift notices.
	//
	// The type of subscription plan being used. Possible values are:
	//	"Prime" // Amazon Prime subscription
	//	"1000" // First level of paid subscription
	//	"2000" // Second level of paid subscription
	//	"3000" // Third level of paid subscription
	MsgParamSubPlan string `json:"msg-param-sub-plan"`
	// Included only with sub, resub, and subgift notices.
	//
	// The display name of the subscription plan. This may be a default name or one created by the
	// channel owner.
	MsgParamSubPlanName string `json:"msg-param-sub-plan-name"`
	// Included only with raid notices.
	//
	// The number of viewers raiding this channel from the broadcaster’s channel.
	MsgParamViewerCount string `json:"msg-param-viewerCount"`
	// Included only with ritual notices.
	//
	// The name of the ritual being celebrated. Possible values are: new_chatter.
	MsgParamRitualName string `json:"msg-param-ritual-name"`
	// Included only with bitsbadgetier notices.
	//
	// The tier of the Bits badge the user just earned. For example, 100, 1000, or 10000.
	MsgParamThreshold string `json:"msg-param-threshold"`
	// Included only with subgift notices.
	//
	// The number of months gifted as part of a single, multi-month gift.
	MsgParamGiftMonths string `json:"msg-param-gift-months"`

	// An ID that uniquely identifies the whisper message.
	MessageID string `json:"message-id"`
	//	An ID that uniquely identifies the whisper thread. The ID is in the form,
	// <smaller-value-user-id>_<larger-value-user-id>.
	ThreadID string `json:"thread-id"`

	// Undocumented
	ClientNonce string `json:"client-nonce"`
	// Undocumented
	Flags string `json:"flags"`

	// This are some tags that are undocumented by twitch, but sent when receiving messages

	// TODO documentation
	CustomRewardID string `json:"custom-reward-id"`

	// TODO documentation
	MessageParamColor string `json:"msg-param-color"`

	// TODO documentation
	MessageParamGoalContribution string `json:"msg-param-goal-contribution-type"`

	// Whether this is the first message by a user in this channel
	FirstMessage bool `json:"first-msg"`

	// Whether this is a message by a returning chatter (more information needed, probably a user
	// who came back to this channel after a long time)
	ReturningChatter bool `json:"returning-chatter"`
}

func ParseRawIRCTags(raw string) IRCMessageTags {
	var b []byte
	b = append(b, '{')
	for i, t := range strings.Split(raw, `;`) {
		if i != 0 {
			b = append(b, ',')
		}
		b = append(b, formatRawIRCTag(t)...)
	}
	b = append(b, '}')
	t := IRCMessageTags{}
	err := json.Unmarshal(b, &t)
	if err != nil {
		log.Printf("Failed to parse Tags err: %+v\nraw: %s\nformated: %s", err, raw, string(b))
		return IRCMessageTags{}
	}
	return t
}

func formatRawIRCTag(raw string) []byte {
	var b []byte
	tagPair := strings.Split(raw, "=")
	if len(tagPair) != 2 {
		return []byte(fmt.Sprintf("\"%s\":\"\"", raw))
	}

	var i IRCMessageTags
	t := reflect.TypeOf(i)
	found := false
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		jsonTag := strings.Split(f.Tag.Get("json"), ",")[0]
		if jsonTag != tagPair[0] {
			continue
		}
		found = true

		tagPair[1] = strings.ReplaceAll(tagPair[1], "\\s", " ")
		tagPair[1] = strings.ReplaceAll(tagPair[1], "\\", "\\\\")

		switch f.Type.Kind() {
		case reflect.Slice:
			tagPair[1] = fmt.Sprintf("[\"%s\"]", strings.ReplaceAll(tagPair[1], ",", "\",\""))
		case reflect.Int:
			tagPair[1] = fmt.Sprintf("%s", tagPair[1])
		case reflect.Bool:
			if tagPair[1] == "1" || tagPair[1] == "true" {
				tagPair[1] = "true"
			} else {
				tagPair[1] = "false"
			}
		case reflect.String:
			tagPair[1] = fmt.Sprintf("\"%s\"", tagPair[1])
		case reflect.Struct:
			if f.Type == reflect.TypeOf(time.Time{}) {
				ts, err := strconv.Atoi(tagPair[1])
				if err != nil {
					log.Printf("Could not parse int from '%s' in %+v (json:'%s'): %+v", tagPair[1], err, f.Type, jsonTag)
					tagPair[1] = fmt.Sprintf("\"%s\"", tagPair[1])
					break
				}
				tagPair[1] = "\"" + time.Unix(0, int64(ts)).Format(time.RFC3339) + "\""
			}
		default:
			tagPair[1] = fmt.Sprintf("\"%s\"", tagPair[1])
			log.Printf("formated %+v '%d' (json:'%s') as string", f.Type, f.Type.Kind(), jsonTag)
		}
		break
	}
	if !found {
		tagPair[1] = fmt.Sprintf("\"%s\"", tagPair[1])
		log.Printf("WARN: unknown key '%s', formatted '%s' as string", tagPair[0], tagPair[1])
	}

	formated := fmt.Sprintf("\"%s\":%s", tagPair[0], tagPair[1])
	b = append(b, formated...)
	return b
}
