package twitchgo

import (
	"fmt"
	"net/http"
	"time"
)

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

// GetUsersByID gets all the Twitch users matching the given user IDs.
func (s *Session) GetUsersByID(userIDs ...string) ([]*User, error) {
	if len(userIDs) == 0 {
		return []*User{}, nil
	}
	queryParams := map[string][]string{
		"id": userIDs,
	}

	var streamData rawUserData
	err := s.requestHelper(http.MethodGet, "/users", queryParams, nil, &streamData)
	if err != nil {
		return []*User{}, fmt.Errorf("get users by id: %v", err)
	}

	return streamData.Data, nil
}

// GetUsersByName gets all the Twitch users matching the given user login names.
func (s *Session) GetUsersByName(userLoginNames ...string) ([]*User, error) {
	if len(userLoginNames) == 0 {
		return []*User{}, nil
	}
	queryParams := map[string][]string{
		"login": userLoginNames,
	}

	var streamData rawUserData
	err := s.requestHelper(http.MethodGet, "/users", queryParams, nil, &streamData)
	if err != nil {
		return []*User{}, fmt.Errorf("get users by name: %v", err)
	}

	return streamData.Data, nil
}
