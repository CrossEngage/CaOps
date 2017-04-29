package jolokia

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// Client implements a Jolokia client inspired by the official Java client
type Client struct {
	HTTPClient *http.Client
	BaseURL    url.URL
}

func (c *Client) getURL(path string) string {
	url := c.BaseURL
	url.Path += path
	return url.String()
}

// Read ...
func (c *Client) Read() {

}

// Write ...
func (c *Client) Write() {

}

// Exec ...
func (c *Client) Exec() {

}

// Search ...
func (c *Client) Search() {

}

// List ...
func (c *Client) List() {

}

// Version ...
func (c *Client) Version() (readResp *ReadResponse, err error) {
	resp, err := c.HTTPClient.Get(c.getURL("/version"))
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&readResp); err != nil {
		return nil, err
	}
	return
}
