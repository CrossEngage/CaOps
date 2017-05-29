package jolokia

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
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

// ReadInto ...
func (c *Client) ReadInto(mbean string, response ValueResponse) (err error) {
	resp, err := c.HTTPClient.Get(c.getURL("/read/" + mbean))
	if err != nil {
		return err
	}
	if err := response.DecodeJSON(resp.Body); err != nil {
		return err
	}
	if err := response.Error(); err != nil {
		return err
	}
	return
}

// ReadStringList ...
func (c *Client) ReadStringList(mbean string) (vr *StringListValueResponse, err error) {
	resp, err := c.HTTPClient.Get(c.getURL("/read/" + mbean))
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&vr); err != nil {
		return nil, err
	}
	if err := vr.Error(); err != nil {
		return nil, err
	}
	return
}

// ReadString ...
func (c *Client) ReadString(mbean string) (vr *StringValueResponse, err error) {
	resp, err := c.HTTPClient.Get(c.getURL("/read/" + mbean))
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&vr); err != nil {
		return nil, err
	}
	if err := vr.Error(); err != nil {
		return nil, err
	}
	return
}

// ReadStringMapString ...
func (c *Client) ReadStringMapString(mbean string) (vr *StringMapStringValueResponse, err error) {
	resp, err := c.HTTPClient.Get(c.getURL("/read/" + mbean))
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&vr); err != nil {
		return nil, err
	}
	if err := vr.Error(); err != nil {
		return nil, err
	}
	return
}

// ReadBool ...
func (c *Client) ReadBool(mbean string) (vr *BoolValueResponse, err error) {
	resp, err := c.HTTPClient.Get(c.getURL("/read/" + mbean))
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&vr); err != nil {
		return nil, err
	}
	if err := vr.Error(); err != nil {
		return nil, err
	}
	return
}

// Exec ...
func (c *Client) Exec(mbean, operation string, args ...interface{}) (r *Response, err error) {
	request := &Request{
		Type:      "exec",
		MBean:     mbean,
		Operation: operation,
		Arguments: args,
	}

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(jsonBytes)
	resp, err := c.HTTPClient.Post(c.getURL("/"), "application/json", buffer)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if err := r.Error(); err != nil {
		return nil, err
	}
	return
}

// ExecInto ...
func (c *Client) ExecInto(response ValueResponse, mbean, operation string, args ...interface{}) error {
	request := &Request{
		Type:      "exec",
		MBean:     mbean,
		Operation: operation,
		Arguments: args,
	}

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(jsonBytes)
	resp, err := c.HTTPClient.Post(c.getURL("/"), "application/json", buffer)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := response.DecodeJSON(resp.Body); err != nil {
		return err
	}
	if err := response.Error(); err != nil {
		return err
	}
	return nil
}

// StringMapStringValueResponse contains the response envelop and a string list value
type StringMapStringValueResponse struct {
	Response
	Value map[string]string `json:"value,omitempty"`
}

// StringListValueResponse contains the response envelop and a string list value
type StringListValueResponse struct {
	Response
	Value []string `json:"value,omitempty"`
}

// StringValueResponse contains the response envelop and a string value
type StringValueResponse struct {
	Response
	Value string `json:"value,omitempty"`
}

// Uint64ValueResponse contains the response envelop and int64/long value
type Uint64ValueResponse struct {
	Response
	Value uint64 `json:"value,omitempty"`
}

// DecodeJSON ...
func (vr *Uint64ValueResponse) DecodeJSON(r io.Reader) error {
	dec := json.NewDecoder(r)
	if err := dec.Decode(vr); err != nil {
		return err
	}
	return nil
}

// BoolValueResponse contains the response envelop and a boolean value
type BoolValueResponse struct {
	Response
	Value bool `json:"value,omitempty"`
}

// ListResponseValue represents the response from the List() call
type ListResponseValue map[string]MBeanPackage

// ListResponse contains the envelop and the value of any /list response
type ListResponse struct {
	Response
	Value ListResponseValue `json:"value,omitempty"`
}

// Request ...
type Request struct {
	Type      string        `json:"type,omitempty"`
	MBean     string        `json:"mbean,omitempty"`
	Operation string        `json:"operation,omitempty"`
	Arguments []interface{} `json:"arguments,omitempty"`
}

// Response represents a Jolokia response envelope
type Response struct {
	Timestamp  int64   `json:"timestamp,omitempty"`
	Status     int64   `json:"status,omitempty"`
	Request    Request `json:"request,omitempty"`
	ErrorText  string  `json:"error,omitempty"`
	ErrorType  string  `json:"error_type,omitempty"`
	Stacktrace string  `json:"stacktrace,omitempty"`
}

// AnyResponse ...
type AnyResponse interface {
	Error() error
}

// Errored returns true of the Response got errors
func (r *Response) Error() error {
	if len(r.ErrorText) > 0 {
		return errors.New(r.ErrorText)
	}
	return nil
}

// ValueResponse ...
type ValueResponse interface {
	DecodeJSON(r io.Reader) error
	Error() error
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
