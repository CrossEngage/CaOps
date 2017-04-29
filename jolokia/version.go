package jolokia

// VersionResponseValue represents the response from the Version() call
type VersionResponseValue struct {
	Protocol string `json:"protocol,omitempty"`
	Agent    string `json:"agent,omitempty"`
	Config   struct {
		AgentDescription   string `json:"agentDescription,omitempty"`
		AgentID            string `json:"agentId,omitempty"`
		AgentType          string `json:"agentType,omitempty"`
		SerializeException string `json:"serializeException,omitempty"`
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
	ReadResponse
	Value VersionResponseValue `json:"value,omitempty"`
}
