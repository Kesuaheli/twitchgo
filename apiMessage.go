package twitchgo

import "net/http"

// DeleteNessage tries to delete the given message from the broadcaster's chat. The current session
// has to have the "moderator:manage:chat_message" permission.
func (s *Session) DeleteMessage(broadcasterID, msgID string) (err error) {
	user, err := s.GetUser()
	if err != nil {
		return err
	}
	if broadcasterID == "" {
		broadcasterID = user.ID
	}

	queryParams := map[string][]string{
		"broadcaster_id": {broadcasterID},
		"moderator_id":   {user.ID},
		"message_id":     {msgID},
	}

	return s.requestHelper(http.MethodDelete, "/moderation/chat", queryParams, nil, nil)
}
