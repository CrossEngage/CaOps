package jolokia

// MBeanPackage ...
type MBeanPackage map[string]MBean

// MBean ...
type MBean struct {
	Attributes  map[string]MBeanAttribute `json:"attr"`
	Description string                    `json:"desc"`
	Operations  map[string]MBeanOperation `json:"op"`
}

// MBeanAttribute ...
type MBeanAttribute struct {
	Description string `json:"desc"`
	ReadWrite   bool   `json:"rw"`
	Type        string `json:"type"`
}

// MBeanOperationArgument ...
type MBeanOperationArgument struct {
	Description string `json:"desc"`
	Name        string `json:"name"`
	Type        string `json:"type"`
}

// MBeanOperation ...
type MBeanOperation struct {
	Arguments   []MBeanOperationArgument `json:"args"`
	Description string                   `json:"desc"`
	ReturnType  string                   `json:"ret"`
}
