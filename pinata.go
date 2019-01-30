package pinata

import (
	"net/http"
	"time"
)

const (
	// DefaultAPIURL is the base url specified by the pinata documentation
	defaultBaseURL = "https://api.pinata.cloud"
	defaultTimeout = 30 * time.Second
)

// Client is a pinata specific wrapper around a standard http.Client
// it handles the required header creation, key storage, and api url
// management.
type Client struct {
	http.Client
	key     string
	secret  string
	BaseURL string
}

// NewClient takes an API Key and API Secret and returns a Client
func NewClient(key, secret string) *Client {
	client := Client{
		key:     key,
		secret:  secret,
		BaseURL: defaultBaseURL,
	}
	client.Timeout = defaultTimeout
	return &client
}

// newRequest returns a http.NewRequest with the required headers for Pinata
func (c Client) newRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("pinata_api_key", c.key)
	req.Header.Set("pinata_secret_api_key", c.secret)
	return req, nil
}

// TestAuthentication is used to test authentication credentials and returns an http Response and error
func (c Client) TestAuthentication() (*http.Response, error) {
	req, err := c.newRequest("GET", c.BaseURL+"/data/testAuthentication")
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
