package mbean

import "bitbucket.org/crossengage/athena/cassandra/jolokia"

// StorageService is analogous to the original StorageService on Cassandra,
// except that all JMX calls are made through a Jolokia agent.
type StorageService struct {
	JolokiaClient jolokia.Client
}

const (
	storageServicePath = "org.apache.cassandra.db:type=StorageService"
)

// LiveNodes retrieve the list of live nodes in the cluster, where "liveness"
// is determined by the failure detector of the node being queried.
func (ss *StorageService) LiveNodes() (ips []string, err error) {
	resp, err := ss.JolokiaClient.ReadStringListAttribute(storageServicePath + "/LiveNodes")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// UnreachableNodes retrieve the list of unreachable nodes in the cluster, as
// determined by this node's failure detector.
func (ss *StorageService) UnreachableNodes() (ips []string, err error) {
	resp, err := ss.JolokiaClient.ReadStringListAttribute(storageServicePath + "/UnreachableNodes")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// JoiningNodes retrieve the list of nodes currently bootstrapping into the ring.
func (ss *StorageService) JoiningNodes() (ips []string, err error) {
	resp, err := ss.JolokiaClient.ReadStringListAttribute(storageServicePath + "/JoiningNodes")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// LeavingNodes retrieve the list of nodes currently leaving the ring.
func (ss *StorageService) LeavingNodes() (ips []string, err error) {
	resp, err := ss.JolokiaClient.ReadStringListAttribute(storageServicePath + "/LeavingNodes")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// MovingNodes retrieve the list of nodes currently moving in the ring.
func (ss *StorageService) MovingNodes() (ips []string, err error) {
	resp, err := ss.JolokiaClient.ReadStringListAttribute(storageServicePath + "/MovingNodes")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// // Fetch string representations of the tokens for this node.
// func (ss *StorageService) Tokens() (tokens []string) {}

// // Fetch string representations of the tokens for a specified node.
// func (ss *StorageService) TokensAt(endpoint string) (tokens []string) {}

// // Fetch a string representation of the Cassandra version.
// func (ss *StorageService) ReleaseVersion() (version string) {}

// // Fetch a string representation of the current Schema version.
// func (ss *StorageService) SchemaVersion() string {}

// // Get the list of all data file locations from conf
// func (ss *StorageService) AllDataFileLocations() []string {}

// // Get location of the commit log
// func (ss *StorageService) CommitLogLocation() string {}

// // Get location of the saved caches dir
// func (ss *StorageService) SavedCachesLocation() string {}

// // Retrieve a map of range to end points that describe the ring topology
// // of a Cassandra cluster.
// // Map<List<String>, List<String>>
// func (ss *StorageService) RangeToEndpointMap(keyspace string) map[[]string][]string {}

// // Retrieve a map of range to rpc addresses that describe the ring topology
// // of a Cassandra cluster.
// func (ss *StorageService) RangeToRpcaddressMap(keyspace string) map[[]string][]string {}

// // The TokenRange for a given keyspace.
// func (ss *StorageService) describeRingJMX(keyspace string) (tokenRange []string) {}

// // Retrieve a map of pending ranges to endpoints that describe the ring topology
// func (ss *StorageService) PendingRangeToEndpointMap(keyspace string) map[[]string][]string {}

// // Retrieve a map of tokens to endpoints, including the bootstrapping ones.
// func (ss *StorageService) getTokenToEndpointMap() map[string]string {}

// //  Retrieve this hosts unique ID
// func (ss *StorageService) getLocalHostId() string {}

// //  Retrieve the mapping of endpoint to host ID
// func (ss *StorageService) getEndpointToHostId() map[string]string {}

// //  Retrieve the mapping of host ID to endpoint
// func (ss *StorageService) getHostIdToEndpoint() map[string]string {}

// //  Human-readable load value
// func (ss *StorageService) getLoadString() string {}

// //  Human-readable load value.  Keys are IP addresses.
// func (ss *StorageService) getLoadMap() map[string]string {}

// // Return the generation value for this node.
// func (ss *StorageService) getCurrentGenerationNumber() int {}

// // This method returns the N endpoints that are responsible for storing the
// // specified key i.e for replication.
// func (ss *StorageService) getNaturalEndpoints(keyspaceName, cf, key string) []net.IP {}

// // Takes the snapshot for the given keyspaces. A snapshot name must be specified.
// func (ss *StorageService) takeSnapshot(tag string, keyspaceNames ...string) error {}

// // Takes the snapshot of a specific column family. A snapshot name must be specified.
// func (ss *StorageService) takeTableSnapshot(keyspaceName, tableName, tag string) error {}

// // Takes the snapshot of a multiple column family from different keyspaces. A snapshot name must be specified.
// //    the tag given to the snapshot; may not be null or empty
// //    list of tables from different keyspace in the form of ks1.cf1 ks2.cf2
// func (ss *StorageService) takeMultipleTableSnapshot(tag string, tableList ...string) error {}

// // Remove the snapshot with the given name from the given keyspaces.
// // If no tag is specified we will remove all snapshots.
// func (ss *StorageService) clearSnapshot(tag string, keyspaceNames ...string) error {}

// //  Get the details of all the snapshot
// // @return A map of snapshotName to all its details in Tabular form.  Map<String, TabularData>
// func (ss *StorageService) getSnapshotDetails() map[string]interface{} {}

// // Get the true size taken by all snapshots across all keyspaces.
// // @return True size taken by all the snapshots.
// func (ss *StorageService) trueSnapshotsSize() uint64 {}

// // Forces refresh of values stored in system.size_estimates of all column families.
// func (ss *StorageService) refreshSizeEstimates() error {}

// // Verify (checksums of) the given keyspace.
// // If tableNames array is empty, all CFs are verified.
// // The entire sstable will be read to ensure each cell validates if extendedVerify is true
// func (ss *StorageService) verify(extendedVerify bool, keyspaceName string, tableNames ...string) (int, error) {
// }

// // Flush all memtables for the given column families, or all columnfamilies for the given keyspace
// // if none are explicitly listed.
// func (ss *StorageService) forceKeyspaceFlush(keyspaceName string, tableNames ...string) error {}

// //  get the operational mode (leaving, joining, normal, decommissioned, client)
// func (ss *StorageService) getOperationMode() string {}

// //  Returns whether the storage service is starting or not
// func (ss *StorageService) isStarting() bool {}

// // Effective ownership is % of the data each node owns given the keyspace
// // we calculate the percentage using replication factor.
// // If Keyspace == null, this method will try to verify if all the keyspaces
// // in the cluster have the same replication strategies and if yes then we will
// // use the first else a empty Map is returned.
// func (ss *StorageService) effectiveOwnership(keyspace string) (map[net.IP]float, error) {}

// func (ss *StorageService) Keyspaces() []string {}

// func (ss *StorageService) NonSystemKeyspaces() []string {}

// func (ss *StorageService) NonLocalStrategyKeyspaces() []string {}

// // allows a user to see whether gossip is running or not
// func (ss *StorageService) isGossipRunning() bool {}

// // to determine if gossip is disabled
// func (ss *StorageService) isInitialized() bool {}

// func (ss *StorageService) isJoined() bool   {}
// func (ss *StorageService) isDrained() bool  {}
// func (ss *StorageService) isDraining() bool {}

// func (ss *StorageService) isIncrementalBackupsEnabled() bool       {}
// func (ss *StorageService) setIncrementalBackupsEnabled(value bool) {}

// //  Returns the name of the cluster
// func (ss *StorageService) getClusterName() string {}

// //  Returns the cluster partitioner
// func (ss *StorageService) getPartitionerName() string {}
