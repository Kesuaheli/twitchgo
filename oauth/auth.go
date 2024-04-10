package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is the data struct for a auth client
type Client struct {
	RequestURL   string    `json:"request_url"`
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	Scope        string    `json:"scope"`
	ExpiryDate   time.Time `json:"expiry_date"`
}

// Token is a data struct to hold a token response from the auth server
type Token struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
}

// New creates a new client to generate a token from
func New(requestURL, clientID, secret, scope string) *Client {
	c := &Client{
		RequestURL:   requestURL,
		ClientID:     clientID,
		ClientSecret: secret,
		Scope:        scope,
	}

	return c
}

// GenerateToken generates and returns a new token for c
func (c *Client) GenerateToken() (string, error) {
	form := url.Values{}
	form.Set("client_id", c.ClientID)
	form.Set("client_secret", c.ClientSecret)
	form.Set("grant_type", "client_credentials")
	if c.Scope != "" {
		form.Set("scope", c.Scope)
	}
	body := strings.NewReader(form.Encode())

	req, err := http.NewRequest(http.MethodPost, c.RequestURL, body)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid status code expected %d but got %d! body: %s", http.StatusOK, resp.StatusCode, string(data))
	}

	var t Token
	err = json.Unmarshal(data, &t)

	return t.Token, err
}
