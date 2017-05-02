package jolokia

// VersionResponseValue represents the response from the Version() call
type VersionResponseValue struct {
	Protocol string `json:"protocol,omitempty"`
	Agent    string `json:"agent,omitempty"`
	Config   struct {
		AgentDescription string `json:"agentDescription,omitempty"`
		AgentID          string `json:"agentId,omitempty"`
		AgentType        string `json:"agentType,omitempty"`
		MaxDepth string `json:"maxDepth,omitempty"`
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

// {
//     "request": {
//         "type": "version"
//     },
//     "status": 200,
//     "timestamp": 1493718096,
//     "value": {
//         "agent": "1.3.5",
//         "config": {
//             "agentContext": "/jolokia",
//             "agentId": "10.2.2.10-5750-3b9a45b3-jvm",
//             "agentType": "jvm",
//             "debug": "false",
//             "debugMaxEntries": "100",
//             "discoveryEnabled": "false",
//             "historyMaxEntries": "10",
//             "maxCollectionSize": "0",
//             "maxDepth": "15",
//             "maxObjects": "0"
//         },
//         "info": {},
//         "protocol": "7.2"
//     }
// }
