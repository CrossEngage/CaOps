package jolokia

import "time"

// ReadResponse represents a Jolokia response envelope
type ReadResponse struct {
	Timestamp *time.Time `json:"timestamp,omitempty"`
	Status    uint32     `json:"status,omitempty"`
	Request   struct {
		Type string `json:"type,omitempty"`
	} `json:"request,omitempty"`
	Error string `json:"error,omitempty"`
}
