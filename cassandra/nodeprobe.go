package cassandra

import "bitbucket.org/crossengage/athena/cassandra/mbean"

// NodeProbe encapsulates JMX interaction with Cassandra in a similar
// way as nodetool does it internally.
type NodeProbe struct {
	StorageService mbean.StorageService
}
