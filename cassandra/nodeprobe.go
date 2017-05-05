package cassandra

import (
	"net/http"
	"net/url"

	"bitbucket.org/crossengage/athena/cassandra/jolokia"
	"bitbucket.org/crossengage/athena/cassandra/mbean"
)

// NodeProbe encapsulates JMX interaction with Cassandra in a similar
// way as nodetool does it internally.
type NodeProbe struct {
	StorageService mbean.StorageService
}

// NewNodeProbe builds a new Cassandra NodeProbe
func NewNodeProbe(httpClient *http.Client, baseJolokiaURL url.URL) *NodeProbe {
	JolokiaClient := jolokia.Client{
		HTTPClient: http.DefaultClient,
		BaseURL:    baseJolokiaURL,
	}
	return &NodeProbe{
		StorageService: mbean.StorageService{JolokiaClient: JolokiaClient},
	}
}
