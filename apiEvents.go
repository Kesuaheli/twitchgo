package twitchgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"
)

// Subscription represents a single subscription to an event.
type Subscription struct {
	// ID is the unique identifier for this subscription.
	ID string `json:"id,omitempty"`

	// Status is the status of this subscription e.g. it is set to "enabled" when successfully
	// verified and active.
	Status SubscriptionStatus `json:"status,omitempty"`

	// Type is the acual event that triggers.
	Type SubscriptionType `json:"type"`

	// Version is the version number of the event defined in the Type field.
	Version string `json:"version"`

	// Condition contains a list of key-value-pairs of conditions. The 'key' gives a variable to check
	// and 'value' the value to match. For example "broadcaster_user_id":"12345" requires the event
	// to be triggered at the channel of user "12345".
	Condition map[string]string `json:"condition"`

	// Transport gives information about how this subscription is (or will be) delivered.
	Transport SubscriptionTransport `json:"transport"`

	// CreatedAt is the timestamp of creation of this subscription.
	CreatedAt time.Time `json:"created_at,omitempty"`

	// The amount points this subscription costs. The cost is added to a global count. Each
	// application has a fixed amount of available points to use.
	Cost int `json:"cost,omitempty"`
}

// SubscriptionTransport gives information about how a subscription is (or will be) delivered.
type SubscriptionTransport struct {
	// Method is either set to "webhook" or "websocket".
	Method SubscriptionTransportMethod `json:"method"`

	// WebhookCallbackURI gives the complete URI of the webhook.
	//
	// Only when Method == "webhook"
	WebhookCallbackURI string `json:"callback,omitempty"`

	// WebhookSecret is the secret given with the creation of the subscription to veryfiy its
	// correctness.
	//
	// Only when Method == "webhook"
	WebhookSecret string `json:"secret,omitempty"`

	// WebSocketSessionID is the ID the welcome message returns, when connecting to the twitch
	// websocket. More information needed.
	//
	// Only when Method == "websocket"
	WebSocketSessionID string `json:"session_id,omitempty"`
}

// SubscriptionTransportMethod is the method of delivery for a subscription.
type SubscriptionTransportMethod string

// Available delivery methods for a subscription.
const (
	SubscriptionTransportMethodWebhook   SubscriptionTransportMethod = "webhook"
	SubscriptionTransportMethodWebSocket SubscriptionTransportMethod = "websocket"
)

// SubscriptionType is the type of events that can be subscribed to.
type SubscriptionType string

const (
	// EventChannelUpdate sends notifications when a broadcaster updates the
	// category, title, content classification labels, or broadcast language
	// for their channel.
	EventChannelUpdate SubscriptionType = "channel.update"
	// EventChannelFollow sends a notification when the specified broadcaster
	//starts a stream.
	EventStreamOnline SubscriptionType = "stream.online"
	// EventChannelFollow sends a notification when the specified broadcaster
	// stops a stream.
	EventStreamOffline SubscriptionType = "stream.offline"
)

func (st SubscriptionType) GetVersion() string {
	switch st {
	case EventChannelUpdate:
		return "2"
	case EventStreamOnline:
		return "1"
	case EventStreamOffline:
		return "1"
	default:
		log.Printf("Warning: tried to get version for unknown subscription event type '%s'. Returning \"0\"", st)
		return "0"
	}
}

// SubscriptionStatus is the status of a subscription.
type SubscriptionStatus string

const (
	// SubscriptionStatusEnabled is the status of a subscription that is
	//verified and able to send notifications.
	SubscriptionStatusEnabled SubscriptionStatus = "enabled"

	// SubscriptionStatusWebhookCallbackVerificationPending is the status of a
	// webhook subscription that is waiting for Twitch to verify the callback.
	SubscriptionStatusWebhookCallbackVerificationPending SubscriptionStatus = "webhook_callback_verification_pending"
	// SubscriptionStatusWebhookCallbackVerificationFailed is the status of a
	// webhook subscription that failed to verify the callback.
	SubscriptionStatusWebhookCallbackVerificationFailed SubscriptionStatus = "webhook_callback_verification_failed"
	// SubscriptionStatusNotificationFailuresExceeded is the status of a webhhok
	// subscription that was revoked because the notification delivery failure
	// rate was too high.
	SubscriptionStatusNotificationFailuresExceeded SubscriptionStatus = "notification_failures_exceeded"

	// SubscriptionStatusAuthorizationRevoked is the status of a subscription
	// that was revoked because the user revoked their authorization.
	SubscriptionStatusAuthorizationRevoked SubscriptionStatus = "authorization_revoked"
	// SubscriptionStatusModeratorRemoved is the status of a subscription that
	// was revoked because the moderator was removed.
	SubscriptionStatusModeratorRemoved SubscriptionStatus = "moderator_removed"
	// SubscriptionStatusUserRemoved is the status of a subscription that was
	// revoked because the users in the condition object are no longer Twitch
	// users.
	SubscriptionStatusUserRemoved SubscriptionStatus = "user_removed"
	// SubscriptionStatusVersionRemoved is the status of a subscription that was
	// revoked because the subscription to subscription type and version is no
	// longer supported.
	SubscriptionStatusVersionRemoved SubscriptionStatus = "version_removed"
	// SubscriptionStatusBetaMaintenance is the status of a subscription that
	// was revoked because the beta subscription type was undergoing maintenance.
	SubscriptionStatusBetaMaintenance SubscriptionStatus = "beta_maintenance"

	// SubscriptionStatusWebSocketDisconnected is the status of a subscription
	// that was revoked because the client closed the connection.
	SubscriptionStatusWebSocketDisconnected SubscriptionStatus = "websocket_disconnected"
	// SubscriptionStatusWebSocketFailedPingPong is the status of a websocket
	// subscription that was closed because the client failed to respond to a
	// ping message.
	SubscriptionStatusWebSocketFailedPingPong SubscriptionStatus = "websocket_failed_ping_pong"
	// SubscriptionStatusWebSocketReceivedInboundTraffic is the status of a
	// websocket subscription that was closed because the client sent a non-pong
	// message.
	SubscriptionStatusWebSocketReceivedInboundTraffic SubscriptionStatus = "websocket_received_inbound_traffic"
	// SubscriptionStatusWebSocketConnectionUnused is the status of a websocket
	// subscription that was closed because the client failed to subscribe to
	// events within the required time.
	SubscriptionStatusWebSocketConnectionUnused SubscriptionStatus = "websocket_connection_unused"
	// SubscriptionStatusWebSocketInternalError is the status of a websocket
	// subscription that was closed because the Twitch WebSocket server
	// experienced an unexpected error.
	SubscriptionStatusWebSocketInternalError SubscriptionStatus = "websocket_internal_error"
	// SubscriptionStatusWebSocketNetworkTimeout is the status of a websocket
	// subscription that was closed because the Twitch WebSocket server timed
	// out writing the message to the client.
	SubscriptionStatusWebSocketNetworkTimeout SubscriptionStatus = "websocket_network_timeout"
	// SubscriptionStatusWebSocketNetworkError is the status of a websocket
	// subscription that was closed because the Twitch WebSocket server
	// experienced a network error writing the message to the client.
	SubscriptionStatusWebSocketNetworkError SubscriptionStatus = "websocket_network_error"
)

// GetSubscriptions returns all subscriptions for the authenticated user.
//
// If onlyEnabled is set to true, only enabled subscriptions are returned.
func (s *Session) GetSubscriptions(onlyEnabled bool) (subscriptions []*Subscription, err error) {
	subscriptionsResult := struct {
		Data       []*Subscription `json:"data"`
		Pagination pagination      `json:"pagination"`
	}{}

	queryParams := make(url.Values)
	if onlyEnabled {
		queryParams.Set("status", "enabled")
	}
	for {
		err = s.requestHelper("GET", "/eventsub/subscriptions", queryParams, nil, &subscriptionsResult)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscriptionsResult.Data...)
		if subscriptionsResult.Pagination.Cursor == "" {
			break
		}
		queryParams.Set("after", subscriptionsResult.Pagination.Cursor)
	}
	return subscriptions, nil
}

// DeleteSubscription deletes the subscription with the specified ID.
func (s *Session) DeleteSubscription(id string) (err error) {
	queryParams := make(url.Values)
	queryParams.Set("id", id)
	return s.requestHelper("DELETE", "/eventsub/subscriptions", queryParams, nil, nil)
}

// SubscribeToEvent is a helper function to subscribe to the specified event.
func (s *Session) SubscribeToEvent(broadcasterID, callbackURL string, event SubscriptionType) (err error) {
	subData := &Subscription{
		Type:    event,
		Version: event.GetVersion(),
		Condition: map[string]string{
			"broadcaster_user_id": broadcasterID,
		},
		Transport: SubscriptionTransport{
			Method:             SubscriptionTransportMethodWebhook,
			WebhookCallbackURI: callbackURL,
			WebhookSecret:      s.webhookSecret,
		},
	}
	body := &bytes.Buffer{}
	err = json.NewEncoder(body).Encode(subData)
	if err != nil {
		return fmt.Errorf("encode subscription data: %v", err)
	}

	return s.requestHelper("POST", "/eventsub/subscriptions", nil, body, nil)
}

// SubscribeChannelUpdate subscribes to the channel update event.
//
// This event is triggered when the specified broadcaster updates the category,
// title, content classification labels, or broadcast language for their
// channel.
func (s *Session) SubscribeChannelUpdate(broadcasterID, callbackURL string) (err error) {
	return s.SubscribeToEvent(broadcasterID, callbackURL, EventChannelUpdate)
}

// SubscribeStreamOnline subscribes to the stream online event.
//
// This event is triggered when the specified broadcaster starts a stream.
func (s *Session) SubscribeStreamOnline(broadcasterID, callbackURL string) (err error) {
	return s.SubscribeToEvent(broadcasterID, callbackURL, EventStreamOnline)
}

// SubscribeStreamOffline subscribes to the stream offline event.
//
// This event is triggered when the specified broadcaster stops a stream.
func (s *Session) SubscribeStreamOffline(broadcasterID, callbackURL string) (err error) {
	return s.SubscribeToEvent(broadcasterID, callbackURL, EventStreamOffline)
}
