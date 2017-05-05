package mbean

import (
	"encoding/json"
	"io"

	"bytes"

	"bitbucket.org/crossengage/athena/cassandra/jolokia"
)

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
	resp, err := ss.JolokiaClient.ReadStringList(storageServicePath + "/LiveNodes")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// UnreachableNodes retrieve the list of unreachable nodes in the cluster, as
// determined by this node's failure detector.
func (ss *StorageService) UnreachableNodes() (ips []string, err error) {
	resp, err := ss.JolokiaClient.ReadStringList(storageServicePath + "/UnreachableNodes")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// JoiningNodes retrieve the list of nodes currently bootstrapping into the ring.
func (ss *StorageService) JoiningNodes() (ips []string, err error) {
	resp, err := ss.JolokiaClient.ReadStringList(storageServicePath + "/JoiningNodes")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// LeavingNodes retrieve the list of nodes currently leaving the ring.
func (ss *StorageService) LeavingNodes() (ips []string, err error) {
	resp, err := ss.JolokiaClient.ReadStringList(storageServicePath + "/LeavingNodes")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// MovingNodes retrieve the list of nodes currently moving in the ring.
func (ss *StorageService) MovingNodes() (ips []string, err error) {
	resp, err := ss.JolokiaClient.ReadStringList(storageServicePath + "/MovingNodes")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// Tokens fetch string representations of the tokens for this node.
func (ss *StorageService) Tokens() ([]string, error) {
	resp, err := ss.JolokiaClient.ReadStringList(storageServicePath + "/Tokens")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// ReleaseVersion fetch a string representation of the Cassandra version.
func (ss *StorageService) ReleaseVersion() (version string, err error) {
	resp, err := ss.JolokiaClient.ReadString(storageServicePath + "/ReleaseVersion")
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

// SchemaVersion fetch a string representation of the current Schema version.
func (ss *StorageService) SchemaVersion() (version string, err error) {
	resp, err := ss.JolokiaClient.ReadString(storageServicePath + "/SchemaVersion")
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

// AllDataFileLocations returns the list of all data file locations from conf
func (ss *StorageService) AllDataFileLocations() (paths []string, err error) {
	resp, err := ss.JolokiaClient.ReadStringList(storageServicePath + "/AllDataFileLocations")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// CommitLogLocation returns the location of the commit log
func (ss *StorageService) CommitLogLocation() (string, error) {
	resp, err := ss.JolokiaClient.ReadString(storageServicePath + "/CommitLogLocation")
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

// SavedCachesLocation returns the location of the saved caches dir
func (ss *StorageService) SavedCachesLocation() (string, error) {
	resp, err := ss.JolokiaClient.ReadString(storageServicePath + "/SavedCachesLocation")
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

// // Retrieve a map of range to end points that describe the ring topology
// // of a Cassandra cluster.
// // Map<List<String>, List<String>>
// func (ss *StorageService) RangeToEndpointMap(keyspace string) map[[]string][]string {}

// // Retrieve a map of range to rpc addresses that describe the ring topology
// // of a Cassandra cluster.
// func (ss *StorageService) RangeToRpcaddressMap(keyspace string) map[[]string][]string {}

// // The TokenRange for a given keyspace.
// func (ss *StorageService) DescribeRingJMX(keyspace string) (tokenRange []string) {}

// // Retrieve a map of pending ranges to endpoints that describe the ring topology
// func (ss *StorageService) PendingRangeToEndpointMap(keyspace string) map[[]string][]string {}

// TokenToEndpointMap retrieve a map of tokens to endpoints, including the bootstrapping ones.
func (ss *StorageService) TokenToEndpointMap() (map[string]string, error) {
	resp, err := ss.JolokiaClient.ReadStringMapString(storageServicePath + "/TokenToEndpointMap")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// LocalHostID returns the hosts unique ID
func (ss *StorageService) LocalHostID() (string, error) {
	resp, err := ss.JolokiaClient.ReadString(storageServicePath + "/LocalHostId")
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

// EndpointToHostID retrieve the mapping of endpoint to host ID
func (ss *StorageService) EndpointToHostID() (map[string]string, error) {
	resp, err := ss.JolokiaClient.ReadStringMapString(storageServicePath + "/EndpointToHostId")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// HostIDToEndpoint retrieve the mapping of host ID to endpoint
func (ss *StorageService) HostIDToEndpoint() (map[string]string, error) {
	resp, err := ss.JolokiaClient.ReadStringMapString(storageServicePath + "/HostIdToEndpoint")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// LoadString human-readable load value
func (ss *StorageService) LoadString() (string, error) {
	resp, err := ss.JolokiaClient.ReadString(storageServicePath + "/LoadString")
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

// LoadMap human-readable load value. Keys are IP addresses.
func (ss *StorageService) LoadMap() (map[string]string, error) {
	resp, err := ss.JolokiaClient.ReadStringMapString(storageServicePath + "/LoadMap")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// // This method returns the N endpoints that are responsible for storing the
// // specified key i.e for replication.
// func (ss *StorageService) NaturalEndpoints(keyspaceName, cf, key string) []net.IP {}

// TakeSnapshot is self-explanatory
func (ss *StorageService) TakeSnapshot(tag string, keyspaces ...string) error {
	args := make([]interface{}, 2)
	args[0] = tag
	args[1] = keyspaces
	_, err := ss.JolokiaClient.Exec(storageServicePath, "takeSnapshot", args...)
	if err != nil {
		return err
	}
	return nil
}

// // Takes the snapshot of a specific column family. A snapshot name must be specified.
// func (ss *StorageService) TakeTableSnapshot(keyspaceName, tableName, tag string) error {}

// // Takes the snapshot of a multiple column family from different keyspaces. A snapshot name must be specified.
// //    the tag given to the snapshot; may not be null or empty
// //    list of tables from different keyspace in the form of ks1.cf1 ks2.cf2
// func (ss *StorageService) takeMultipleTableSnapshot(tag string, tableList ...string) error {}

// // Remove the snapshot with the given name from the given keyspaces.
// // If no tag is specified we will remove all snapshots.
// func (ss *StorageService) clearSnapshot(tag string, keyspaceNames ...string) error {}

// SnapshotDetailsResponse encapsulates the weird response when we get while reading this
// property from Cassandra
type SnapshotDetailsResponse struct {
	jolokia.Response
	Value SnapshotDetails `json:"value"`
}

// SnapshotDetails ...
type SnapshotDetails []TableSnapshot

// UnmarshalJSON handles the messy hierarchical json for tabular data that comes from Jolokia
func (sd *SnapshotDetails) UnmarshalJSON(buf []byte) error {
	if string(buf) == "null" {
		return nil
	}

	var rawMaps map[string]map[string]map[string]map[string]map[string]map[string]TableSnapshot
	if err := json.NewDecoder(bytes.NewBuffer(buf)).Decode(&rawMaps); err != nil {
		return err
	}

	// Yes, I know this is ugly, but it was the faster way to decode a mess
	for _, lv1 := range rawMaps {
		for _, lv2 := range lv1 {
			for _, lv3 := range lv2 {
				for _, lv4 := range lv3 {
					for _, lv5 := range lv4 {
						for _, ts := range lv5 {
							*sd = append(*sd, ts)
						}
					}
				}
			}
		}
	}

	return nil
}

// TableSnapshot ...
type TableSnapshot struct {
	SnapshotName string `json:"Snapshot name"`
	Keyspace     string `json:"Keyspace name"`
	Table        string `json:"Column family name"`
	SizeOnDisk   string `json:"Size on disk"`
	TrueSize     string `json:"True size"`
}

// DecodeJSON ...
func (sdr *SnapshotDetailsResponse) DecodeJSON(r io.Reader) error {
	dec := json.NewDecoder(r)
	if err := dec.Decode(sdr); err != nil {
		return err
	}
	return nil
}

// SnapshotDetails get the details of all the snapshots
func (ss *StorageService) SnapshotDetails() (*SnapshotDetailsResponse, error) {
	details := &SnapshotDetailsResponse{}
	err := ss.JolokiaClient.ReadInto(storageServicePath+"/SnapshotDetails", details)
	return details, err
}

// AllSnapshotsSize get the true size taken by all snapshots across all keyspaces.
func (ss *StorageService) AllSnapshotsSize() (uint64, error) {
	response := &jolokia.Uint64ValueResponse{}
	err := ss.JolokiaClient.ExecInto(response, storageServicePath, "trueSnapshotsSize")
	if err != nil {
		return 0, err
	}
	return response.Value, nil
}

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

// Starting returns whether the storage service is starting or not
func (ss *StorageService) Starting() (bool, error) {
	resp, err := ss.JolokiaClient.ReadBool(storageServicePath + "/IsStarting")
	if err != nil {
		return false, err
	}
	return resp.Value, nil
}

// GossipRunning returns whether the gossip is running
func (ss *StorageService) GossipRunning() (bool, error) {
	resp, err := ss.JolokiaClient.ReadBool(storageServicePath + "/GossipRunning")
	if err != nil {
		return false, err
	}
	return resp.Value, nil
}

// // Effective ownership is % of the data each node owns given the keyspace
// // we calculate the percentage using replication factor.
// // If Keyspace == null, this method will try to verify if all the keyspaces
// // in the cluster have the same replication strategies and if yes then we will
// // use the first else a empty Map is returned.
// func (ss *StorageService) effectiveOwnership(keyspace string) (map[net.IP]float, error) {}

// Keyspaces return the list of keyspaces in the cluster
func (ss *StorageService) Keyspaces() ([]string, error) {
	resp, err := ss.JolokiaClient.ReadStringList(storageServicePath + "/Keyspaces")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// NonSystemKeyspaces ...
func (ss *StorageService) NonSystemKeyspaces() ([]string, error) {
	resp, err := ss.JolokiaClient.ReadStringList(storageServicePath + "/NonSystemKeyspaces")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// OperationMode returns the operation mode of the node. STARTING, NORMAL, JOINING,
// LEAVING, DECOMMISSIONED, MOVING, DRAINING, DRAINED.
func (ss *StorageService) OperationMode() (string, error) {
	resp, err := ss.JolokiaClient.ReadString(storageServicePath + "/OperationMode")
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

// IncrementalBackupsEnabled is self explanatory
func (ss *StorageService) IncrementalBackupsEnabled() (bool, error) {
	resp, err := ss.JolokiaClient.ReadBool(storageServicePath + "/IncrementalBackupsEnabled")
	if err != nil {
		return false, err
	}
	return resp.Value, nil
}

// Initialized is self explanatory
func (ss *StorageService) Initialized() (bool, error) {
	resp, err := ss.JolokiaClient.ReadBool(storageServicePath + "/Initialized")
	if err != nil {
		return false, err
	}
	return resp.Value, nil
}

// Joined is self explanatory
func (ss *StorageService) Joined() (bool, error) {
	resp, err := ss.JolokiaClient.ReadBool(storageServicePath + "/Joined")
	if err != nil {
		return false, err
	}
	return resp.Value, nil
}

// ClusterName returns the name of the cluster
func (ss *StorageService) ClusterName() (name string, err error) {
	resp, err := ss.JolokiaClient.ReadString(storageServicePath + "/ClusterName")
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

// PartitionerName returns the cluster partitioner
func (ss *StorageService) PartitionerName() (name string, err error) {
	resp, err := ss.JolokiaClient.ReadString(storageServicePath + "/PartitionerName")
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}
