package pinata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

// Metadata is a general purpose data structure used for submitted
// metadata in api requests
type Metadata struct {
	Name      string                 `json:"name,omitempty"`
	Keyvalues map[string]interface{} `json:"keyvalues,omitempty"`
}

// HashToPinRequest is the request structure as outlined in
// https://pinata.cloud/documentation#PinHashToIPFS
type HashToPinRequest struct {
	Metadata  *Metadata `json:"pinataMetadata,omitempty"`
	HashToPin string    `json:"hashToPin"`
}

// NewMetadata returns a Metadata struct with initialized keyvalues map
func NewMetadata() Metadata {
	return Metadata{Keyvalues: make(map[string]interface{})}
}

// NewMetadataWithName initializes a Meatadata struct and sets the name
func NewMetadataWithName(name string) Metadata {
	m := NewMetadata()
	m.Name = name
	return m
}

// SetKeyValue safely sets the KeyValue store with one of the 3 supported types
// outlined in the documentation
// 		JSON             GO
// ----------------------------------------------
// - Strings             string
// - Numbers             int, int64, float32, float64
// - Dates               time.Time
// if time.Time is passed in as a value, this method formats it with `.UTC().Format(time.RFC3339)`
// which is a stricter implementation of ISO 8601. If you do not want it formatted in UTC
// pass in a string encoded timestamp
func (m *Metadata) SetKeyValue(key string, value interface{}) error {
	switch t := value.(type) {
	case string, int, int64, float32, float64:
		m.Keyvalues[key] = value
	case time.Time:
		m.Keyvalues[key] = value.(time.Time).UTC().Format(time.RFC3339)
	default:
		return fmt.Errorf("unsupported type %T for Metadata key/value store", t)
	}
	return nil
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

// NewRequestWithHeaders returns a http.NewRequest with the required headers for Pinata
func (c Client) NewRequestWithHeaders(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("pinata_api_key", c.key)
	req.Header.Set("pinata_secret_api_key", c.secret)
	return req, nil
}

// TestAuthentication is used to test authentication credentials and returns an http Response and error
func (c Client) TestAuthentication() (*http.Response, error) {
	req, err := c.NewRequestWithHeaders("GET", c.BaseURL+"/data/testAuthentication", nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// PinHashToIPFSWithMetadata takes a hash and metadata struct and returns an http response
func (c Client) PinHashToIPFSWithMetadata(hash string, metadata Metadata) (*http.Response, error) {

	r := HashToPinRequest{
		HashToPin: hash,
	}

	if metadata.Name != "" || len(metadata.Keyvalues) > 0 {
		r.Metadata = &metadata
	}

	// if metadata != Metadata{} {
	// 	r.Metadata = metadata
	// }

	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	req, err := c.NewRequestWithHeaders("POST", c.BaseURL+"/pinning/pinHashToIPFS", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.Do(req)
}

// PinHashToIPFS is used to test authentication credentials and returns an http Response and error
func (c Client) PinHashToIPFS(hash string) (*http.Response, error) {
	return c.PinHashToIPFSWithMetadata(hash, Metadata{})
}
