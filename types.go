package twitchgo

import "time"

// pagination contains information used to page through the list of results. The object is empty if
// there are no more pages left to page through.
type pagination struct {
	// The cursor used to get the next page of results. Set the request’s after or before query
	// parameter to this value depending on whether you’re paging forwards or backwards.
	Cursor string `json:"cursor"`
}

type rawStreamData struct {
	// The list of streams.
	Data []*Stream `json:"data"`

	pagination
}

// Stream represents a twitch live stream with all its informations.
type Stream struct {
	// An ID that identifies the stream. You can use this ID later to look up the video on demand
	// (VOD).
	ID string `json:"id"`
	// The ID of the user that’s broadcasting the stream.
	UserID string `json:"user_id"`
	// The user’s login name.
	UserLogin string `json:"user_login"`
	// The user’s display name.
	UserName string `json:"user_name"`

	// The ID of the category or game being played.
	GameID string `json:"game_id"`
	// The name of the category or game being played.
	GameName string `json:"game_name"`

	// The type of stream. Possible values are:
	//  "live" //
	// If an error occurs, this field is set to an empty string.
	Type string `json:"Type"`

	// The stream’s title. Is an empty string if not set.
	Title string `json:"title"`
	// The tags applied to the stream.
	Tags []string `json:"tags"`
	// The number of users watching the stream.
	ViewerCount int `json:"viewer_count"`
	// The UTC date and time (in RFC3339 format) of when the broadcast began.
	StartedAt time.Time `json:"started_at"`
	// The language that the stream uses. This is an ISO 639-1 two-letter language code or other if
	// the stream uses a language not in the list of supported stream languages.
	Language string `json:"language"`
	// A URL to an image of a frame from the last 5 minutes of the stream. Replace the width and
	// height placeholders in the URL ({width}x{height}) with the size of the image you want, in
	// pixels.
	ThumbnailURL string `json:"thumbnail_url"`

	// A Boolean value that indicates whether the stream is meant for mature audiences.
	IsMature bool `json:"is_mature"`
}

type rawUserData struct {
	// 	The list of users.
	Data []*User `json:"data"`
}

// User represents a twitch user account with all its informations.
type User struct {
	// An ID that identifies the user.
	ID string `json:"id"`
	// The user’s login name.
	Login string `json:"login"`
	// The user’s display name.
	DisplayName string `json:"display_name"`

	// The type of user. Possible values are:
	//
	//  "admin"      // Twitch administrator
	//  "global_mod"
	//  "staff"      // Twitch staff
	//  ""           // Normal user
	Type string `json:"Type"`
	// The type of broadcaster. Possible values are:
	//
	//  "affiliate" // An affiliate broadcaster
	//  "partner"   // A partner broadcaster
	//  ""          // A normal broadcaster
	BroadcasterType string `json:"broadcaster_type"`

	// The user’s description of their channel.
	Description string `json:"description"`
	// A URL to the user’s profile image.
	ProfileImageURL string `json:"profile_image_url"`
	// A URL to the user’s offline image.
	OfflineImageURL string `json:"offline_image_url"`

	// The user’s verified email address. The object includes this field only if the user access
	// token includes the user:read:email scope.
	//
	// If the request contains more than one user, only the user associated with the access token
	// that provided consent will include an email address - the email address for all other users
	// will be empty.
	Email string `json:"email,omitempty"`

	// The UTC date and time that the user’s account was created. The timestamp is in RFC3339 format.
	CreatedAt time.Time `json:"created_at"`
}
