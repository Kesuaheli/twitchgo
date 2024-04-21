package oauth

import (
	"bytes"
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

	lastToken Token
}

// Token is a data struct to hold a token response from the auth server
type Token struct {
	Token        string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	Scopes       []string `json:"scope"`
	ExpiresIn    int      `json:"expires_in"`

	expiresAt time.Time
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

// SetRefreshToken lets set a new custom refresh token to use for generating the next token. It also
// clears the latest saved token so the next call to [Client.GenerateToken] uses the new
// refreshToken.
func (c *Client) SetRefreshToken(refreshToken string) {
	c.lastToken = Token{RefreshToken: refreshToken}
}

// GenerateToken generates and returns a new token for c
func (c *Client) GenerateToken() (string, error) {
	if c.lastToken.expiresAt.After(time.Now()) {
		return c.lastToken.Token, nil
	}

	if c.lastToken.RefreshToken != "" {
		return c.generateFromRefreshToken()
	}

	return c.generateFromCredentials()
}

func (c *Client) generateFromCredentials() (string, error) {
	form := url.Values{}
	form.Set("client_id", c.ClientID)
	form.Set("client_secret", c.ClientSecret)
	form.Set("grant_type", "client_credentials")
	if c.Scope != "" {
		form.Set("scope", c.Scope)
	}

	body := strings.NewReader(form.Encode())
	return c.tokenRequest(body)
}

func (c *Client) generateFromRefreshToken() (string, error) {
	form := url.Values{}
	form.Set("client_id", c.ClientID)
	form.Set("client_secret", c.ClientSecret)
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", c.lastToken.RefreshToken)

	body := strings.NewReader(form.Encode())
	return c.tokenRequest(body)
}

func (c *Client) generateFromAuthorizationCode(code string) (string, error) {
	authCodeBody := struct {
		ClientID          string `json:"client_id"`
		ClientSecret      string `json:"client_secret"`
		AuthorizationCode string `json:"code"`
		GrantType         string `json:"grant_type"`   // always "authorization_code"
		RedirectURI       string `json:"redirect_uri"` // always "https://webhook.cake4everyone.de/auth/twitch"
	}{
		ClientID:          c.ClientID,
		ClientSecret:      c.ClientSecret,
		AuthorizationCode: code,
		GrantType:         "authorization_code",
		RedirectURI:       "https://webhook.cake4everyone.de/auth/twitch",
	}

	rawBody, err := json.Marshal(authCodeBody)
	if err != nil {
		panic(err)
	}
	body := bytes.NewReader(rawBody)

	return c.tokenRequest(body)
}

func (c *Client) tokenRequest(body io.Reader) (string, error) {
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
	t.expiresAt = time.Now().Add(time.Duration(t.ExpiresIn) * time.Second)
	c.lastToken = t

	return t.Token, err
}
