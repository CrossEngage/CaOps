package jolokia

// ListResponseValue represents the response from the List() call
type ListResponseValue map[string]MBeanPackage

// ListResponse contains the envelop and the value of any /list response
type ListResponse struct {
	Response
	Value ListResponseValue `json:"value,omitempty"`
}
