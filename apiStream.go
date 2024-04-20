package twitchgo

import (
	"fmt"
	"net/http"
	"time"
)

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

// GetStreamsByID gets all the streams matching the given user IDs of the streamers.
// Returns only the streams of those users that are broadcasting.
func (s *Session) GetStreamsByID(userIDs ...string) ([]*Stream, error) {
	if len(userIDs) == 0 {
		return []*Stream{}, nil
	}
	queryParams := map[string][]string{
		"user_id": userIDs,
		"first":   {"100"},
	}

	var streamData rawStreamData
	err := s.requestHelper(http.MethodGet, "/streams", queryParams, nil, &streamData)
	if err != nil {
		return []*Stream{}, fmt.Errorf("get streams by id: %v", err)
	}

	return streamData.Data, nil
}

// GetStreamsByName gets all the streams matching the given user login names of the streamers.
// Returns only the streams of those users that are broadcasting.
func (s *Session) GetStreamsByName(userLoginNames ...string) ([]*Stream, error) {
	if len(userLoginNames) == 0 {
		return []*Stream{}, nil
	}
	queryParams := map[string][]string{
		"user_login": userLoginNames,
		"first":      {"100"},
	}

	var streamData rawStreamData
	err := s.requestHelper(http.MethodGet, "/streams", queryParams, nil, &streamData)
	if err != nil {
		return []*Stream{}, fmt.Errorf("get streams by name: %v", err)
	}

	return streamData.Data, nil
}
