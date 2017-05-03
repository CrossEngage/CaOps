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

// ReadAttribute ...
func (c *Client) ReadAttribute(path string) (attrResp *ReadAttributeResponse, err error) {
	resp, err := c.HTTPClient.Get(c.getURL("/read/" + path))
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&attrResp); err != nil {
		return nil, err
	}
	return
}

// ReadStringListAttribute ...
func (c *Client) ReadStringListAttribute(path string) (attrResp *ReadStringListAttributeResponse, err error) {
	resp, err := c.HTTPClient.Get(c.getURL("/read/" + path))
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&attrResp); err != nil {
		return nil, err
	}
	return
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
func (c *Client) Version() (versionResp *VersionResponse, err error) {
	resp, err := c.HTTPClient.Get(c.getURL("/version"))
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&versionResp); err != nil {
		return nil, err
	}
	return
}

// ReadAttributeResponse contains the envelop and the value of the /read/
// attribute response.
type ReadAttributeResponse struct {
	Response
	Value interface{} `json:"value,omitempty"`
}

// ReadStringListAttributeResponse contains the response envelop and the string
// list value of the /read/ attribute response.
type ReadStringListAttributeResponse struct {
	Response
	Value []string `json:"value,omitempty"`
}

// ListResponseValue represents the response from the List() call
type ListResponseValue map[string]MBeanPackage

// ListResponse contains the envelop and the value of any /list response
type ListResponse struct {
	Response
	Value ListResponseValue `json:"value,omitempty"`
}

// Response represents a Jolokia response envelope
type Response struct {
	Timestamp int64 `json:"timestamp,omitempty"`
	Status    int64 `json:"status,omitempty"`
	Request   struct {
		Type string `json:"type,omitempty"`
	} `json:"request,omitempty"`
	Error      string `json:"error,omitempty"`
	ErrorType  string `json:"error_type,omitempty"`
	Stacktrace string `json:"stacktrace,omitempty"`
}

// VersionResponseValue represents the response from the Version() call
type VersionResponseValue struct {
	Protocol string `json:"protocol,omitempty"`
	Agent    string `json:"agent,omitempty"`
	Config   struct {
		AgentDescription string `json:"agentDescription,omitempty"`
		AgentID          string `json:"agentId,omitempty"`
		AgentType        string `json:"agentType,omitempty"`
		MaxDepth         string `json:"maxDepth,omitempty"`
	} `json:"config,omitempty"`
	Info struct {
		Product   string `json:"product,omitempty"`
		Vendor    string `json:"vendor,omitempty"`
		Version   string `json:"version,omitempty"`
		ExtraInfo struct {
			AMXBooted bool `json:"amxBooted,omitempty"`
		} `json:"extraInfo,omitempty"`
	} `json:"info,omitempty"`
}

// VersionResponse contains the envelop and the value of the /version response
type VersionResponse struct {
	Response
	Value VersionResponseValue `json:"value,omitempty"`
}
