package twitchgo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// pagination contains information used to page through the list of results. The object is empty if
// there are no more pages left to page through.
type pagination struct {
	// The cursor used to get the next page of results. Set the request’s after or before query
	// parameter to this value depending on whether you’re paging forwards or backwards.
	Cursor string `json:"cursor"`
}

func (s *Session) requestHelper(method, endpoint string, queryParams map[string][]string, body io.Reader, result any) error {
	req, err := s.buildRequest(method, endpoint, queryParams, body)
	if err != nil {
		return err
	}

	t, err := s.oauth.GenerateToken()
	if err != nil {
		return fmt.Errorf("generate token: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t)
	req.Header.Set("Client-Id", s.clientID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("expected a 2xx status code, but got '%s': %s", resp.Status, respData)
	}

	if result == nil {
		return nil
	}
	return json.Unmarshal(respData, result)
}

func (s *Session) buildRequest(method, endpoint string, queryParams map[string][]string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, baseURL+endpoint, body)
	if err != nil {
		return
	}

	var rawQueries []string
	for k, v := range queryParams {
		for _, v := range v {
			rawQueries = append(rawQueries, fmt.Sprintf("%s=%s", k, v))
		}
	}
	req.URL.RawQuery = strings.Join(rawQueries, "&")
	return
}
