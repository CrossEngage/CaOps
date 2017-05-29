package cassandra

import (
	"bitbucket.org/crossengage/athena/cassandra/mbean"
	"bitbucket.org/crossengage/athena/jolokia"
)

// NodeProbe encapsulates JMX interaction with Cassandra in a similar
// way as nodetool does it internally.
type NodeProbe struct {
	StorageService mbean.StorageService
}

// NewNodeProbe builds a new Cassandra NodeProbe
func NewNodeProbe(jolokiaClient jolokia.Client) *NodeProbe {
	return &NodeProbe{
		StorageService: mbean.StorageService{JolokiaClient: jolokiaClient},
	}
}
