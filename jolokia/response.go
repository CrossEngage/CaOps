package jolokia

// Response represents a Jolokia response envelope
type Response struct {
	Timestamp int64 `json:"timestamp,omitempty"`
	Status    int64 `json:"status,omitempty"`
	Request   struct {
		Type string `json:"type,omitempty"`
	} `json:"request,omitempty"`
	Error string `json:"error,omitempty"`
}
