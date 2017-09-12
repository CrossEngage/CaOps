package cassandra

import (
	"bytes"
	"encoding/json"
	"io"
	"log"

	"github.com/CrossEngage/CaOps/internal/jolokia"
)

// storageService is analogous to the original storageService on Cassandra,
// except that all JMX calls are made through a Jolokia agent.
type storageService struct {
	jolokiaClient jolokia.Client
}

const (
	storageServicePath = "org.apache.cassandra.db:type=StorageService"
)

// LiveNodes retrieve the list of live nodes in the cluster, where "liveness"
// is determined by the failure detector of the node being queried.
func (ss storageService) LiveNodes() (ips []string, err error) {
	resp, err := ss.jolokiaClient.ReadStringList(storageServicePath + "/LiveNodes")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// UnreachableNodes retrieve the list of unreachable nodes in the cluster, as
// determined by this node's failure detector.
func (ss storageService) UnreachableNodes() (ips []string, err error) {
	resp, err := ss.jolokiaClient.ReadStringList(storageServicePath + "/UnreachableNodes")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// JoiningNodes retrieve the list of nodes currently bootstrapping into the ring.
func (ss storageService) JoiningNodes() (ips []string, err error) {
	resp, err := ss.jolokiaClient.ReadStringList(storageServicePath + "/JoiningNodes")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// LeavingNodes retrieve the list of nodes currently leaving the ring.
func (ss storageService) LeavingNodes() (ips []string, err error) {
	resp, err := ss.jolokiaClient.ReadStringList(storageServicePath + "/LeavingNodes")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// MovingNodes retrieve the list of nodes currently moving in the ring.
func (ss storageService) MovingNodes() (ips []string, err error) {
	resp, err := ss.jolokiaClient.ReadStringList(storageServicePath + "/MovingNodes")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// Tokens fetch string representations of the tokens for this node.
func (ss storageService) Tokens() ([]string, error) {
	resp, err := ss.jolokiaClient.ReadStringList(storageServicePath + "/Tokens")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// ReleaseVersion fetch a string representation of the Cassandra version.
func (ss storageService) ReleaseVersion() (version string, err error) {
	resp, err := ss.jolokiaClient.ReadString(storageServicePath + "/ReleaseVersion")
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

// SchemaVersion fetch a string representation of the current Schema version.
func (ss storageService) SchemaVersion() (version string, err error) {
	resp, err := ss.jolokiaClient.ReadString(storageServicePath + "/SchemaVersion")
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

// AllDataFileLocations returns the list of all data file locations from conf
func (ss storageService) AllDataFileLocations() (paths []string, err error) {
	resp, err := ss.jolokiaClient.ReadStringList(storageServicePath + "/AllDataFileLocations")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// CommitLogLocation returns the location of the commit log
func (ss storageService) CommitLogLocation() (string, error) {
	resp, err := ss.jolokiaClient.ReadString(storageServicePath + "/CommitLogLocation")
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

// SavedCachesLocation returns the location of the saved caches dir
func (ss storageService) SavedCachesLocation() (string, error) {
	resp, err := ss.jolokiaClient.ReadString(storageServicePath + "/SavedCachesLocation")
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

// // Retrieve a map of range to end points that describe the ring topology
// // of a Cassandra cluster.
// // Map<List<String>, List<String>>
// func (ss storageService) RangeToEndpointMap(keyspace string) map[[]string][]string {}

// // Retrieve a map of range to rpc addresses that describe the ring topology
// // of a Cassandra cluster.
// func (ss storageService) RangeToRpcaddressMap(keyspace string) map[[]string][]string {}

// // The TokenRange for a given keyspace.
// func (ss storageService) DescribeRingJMX(keyspace string) (tokenRange []string) {}

// // Retrieve a map of pending ranges to endpoints that describe the ring topology
// func (ss storageService) PendingRangeToEndpointMap(keyspace string) map[[]string][]string {}

// TokenToEndpointMap retrieve a map of tokens to endpoints, including the bootstrapping ones.
func (ss storageService) TokenToEndpointMap() (map[string]string, error) {
	resp, err := ss.jolokiaClient.ReadStringMapString(storageServicePath + "/TokenToEndpointMap")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// LocalHostID returns the hosts unique ID
func (ss storageService) LocalHostID() (string, error) {
	resp, err := ss.jolokiaClient.ReadString(storageServicePath + "/LocalHostId")
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

// EndpointToHostID retrieve the mapping of endpoint to host ID
func (ss storageService) EndpointToHostID() (map[string]string, error) {
	resp, err := ss.jolokiaClient.ReadStringMapString(storageServicePath + "/EndpointToHostId")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// HostIDToEndpoint retrieve the mapping of host ID to endpoint
func (ss storageService) HostIDToEndpoint() (map[string]string, error) {
	resp, err := ss.jolokiaClient.ReadStringMapString(storageServicePath + "/HostIdToEndpoint")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// LoadString human-readable load value
func (ss storageService) LoadString() (string, error) {
	resp, err := ss.jolokiaClient.ReadString(storageServicePath + "/LoadString")
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

// LoadMap human-readable load value. Keys are IP addresses.
func (ss storageService) LoadMap() (map[string]string, error) {
	resp, err := ss.jolokiaClient.ReadStringMapString(storageServicePath + "/LoadMap")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// // This method returns the N endpoints that are responsible for storing the
// // specified key i.e for replication.
// func (ss storageService) NaturalEndpoints(keyspaceName, cf, key string) []net.IP {}

// TakeSnapshot is self-explanatory
func (ss storageService) TakeSnapshot(tag string, keyspaces ...string) error {
	args := make([]interface{}, 2)
	args[0] = tag
	args[1] = keyspaces
	r, err := ss.jolokiaClient.Exec(storageServicePath, "takeSnapshot", args...)
	if err != nil {
		return err
	}
	log.Println(r)
	return nil
}

// TakeTableSnapshot is self-explanatory
func (ss storageService) TakeTableSnapshot(tag, keyspace, table string) error {
	args := make([]interface{}, 3)
	args[0] = keyspace
	args[1] = table
	args[2] = tag
	r, err := ss.jolokiaClient.Exec(storageServicePath, "takeTableSnapshot", args...)
	if err != nil {
		return err
	}
	log.Println(r)
	return nil
}

// TakeMultipleTableSnapshot takes the snapshot of a multiple column family from different
// keyspaces. A snapshot name must be specified.
func (ss storageService) TakeMultipleTableSnapshot(tag string, tableList ...string) error {
	args := make([]interface{}, 2)
	args[0] = tag
	args[1] = tableList
	r, err := ss.jolokiaClient.Exec(storageServicePath, "takeMultipleTableSnapshot", args...)
	if err != nil {
		return err
	}
	log.Println(r)
	return nil
}

// ClearSnapshot remove the snapshot with the given name from the given keyspaces.
// If no tag is specified we will remove all snapshots.
func (ss storageService) ClearSnapshot(tag string, keyspaces ...string) error {
	args := make([]interface{}, 2)
	args[0] = tag
	args[1] = keyspaces
	r, err := ss.jolokiaClient.Exec(storageServicePath, "clearSnapshot", args...)
	if err != nil {
		return err
	}
	log.Println(r)
	return nil
}

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

	// Yes, I know this is ugly, but it was the faster (quick-and-dirty) way to decode a mess
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
func (ss storageService) SnapshotDetails() (*SnapshotDetailsResponse, error) {
	details := &SnapshotDetailsResponse{}
	err := ss.jolokiaClient.ReadInto(storageServicePath+"/SnapshotDetails", details)
	return details, err
}

// AllSnapshotsSize get the true size taken by all snapshots across all keyspaces.
func (ss storageService) AllSnapshotsSize() (uint64, error) {
	response := &jolokia.Uint64ValueResponse{}
	err := ss.jolokiaClient.ExecInto(response, storageServicePath, "trueSnapshotsSize")
	if err != nil {
		return 0, err
	}
	log.Printf("\n%+v\n\n", response)
	return response.Value, nil
}

// RefreshSizeEstimates forces refresh of values stored in system.size_estimates of all
// column families.
func (ss storageService) refreshSizeEstimates() error {
	r, err := ss.jolokiaClient.Exec(storageServicePath, "refreshSizeEstimates")
	if err != nil {
		return err
	}
	log.Println(r)
	return nil
}

// // Verify (checksums of) the given keyspace.
// // If tableNames array is empty, all CFs are verified.
// // The entire sstable will be read to ensure each cell validates if extendedVerify is true
// func (ss storageService) verify(extendedVerify bool, keyspaceName string, tableNames ...string) (int, error) {
// }

// ForceKeyspaceFlush flush all memtables for the given column families, or all columnfamilies for
// the given keyspace if none are explicitly listed.
func (ss storageService) ForceKeyspaceFlush(keyspace string, tables ...string) error {
	args := make([]interface{}, 2)
	args[0] = keyspace
	args[1] = tables
	r, err := ss.jolokiaClient.Exec(storageServicePath, "forceKeyspaceFlush", args...)
	if err != nil {
		return err
	}
	log.Println(r)
	return nil
}

// Starting returns whether the storage service is starting or not
func (ss storageService) Starting() (bool, error) {
	resp, err := ss.jolokiaClient.ReadBool(storageServicePath + "/IsStarting")
	if err != nil {
		return false, err
	}
	return resp.Value, nil
}

// GossipRunning returns whether the gossip is running
func (ss storageService) GossipRunning() (bool, error) {
	resp, err := ss.jolokiaClient.ReadBool(storageServicePath + "/GossipRunning")
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
// func (ss storageService) effectiveOwnership(keyspace string) (map[net.IP]float, error) {}

// Keyspaces return the list of keyspaces in the cluster
func (ss storageService) Keyspaces() ([]string, error) {
	resp, err := ss.jolokiaClient.ReadStringList(storageServicePath + "/Keyspaces")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// NonSystemKeyspaces ...
func (ss storageService) NonSystemKeyspaces() ([]string, error) {
	resp, err := ss.jolokiaClient.ReadStringList(storageServicePath + "/NonSystemKeyspaces")
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// OperationMode returns the operation mode of the node. STARTING, NORMAL, JOINING,
// LEAVING, DECOMMISSIONED, MOVING, DRAINING, DRAINED.
func (ss storageService) OperationMode() (string, error) {
	resp, err := ss.jolokiaClient.ReadString(storageServicePath + "/OperationMode")
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

// IncrementalBackupsEnabled is self explanatory
func (ss storageService) IncrementalBackupsEnabled() (bool, error) {
	resp, err := ss.jolokiaClient.ReadBool(storageServicePath + "/IncrementalBackupsEnabled")
	if err != nil {
		return false, err
	}
	return resp.Value, nil
}

// Initialized is self explanatory
func (ss storageService) Initialized() (bool, error) {
	resp, err := ss.jolokiaClient.ReadBool(storageServicePath + "/Initialized")
	if err != nil {
		return false, err
	}
	return resp.Value, nil
}

// Joined is self explanatory
func (ss storageService) Joined() (bool, error) {
	resp, err := ss.jolokiaClient.ReadBool(storageServicePath + "/Joined")
	if err != nil {
		return false, err
	}
	return resp.Value, nil
}

// ClusterName returns the name of the cluster
func (ss storageService) ClusterName() (name string, err error) {
	resp, err := ss.jolokiaClient.ReadString(storageServicePath + "/ClusterName")
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}

// PartitionerName returns the cluster partitioner
func (ss storageService) PartitionerName() (name string, err error) {
	resp, err := ss.jolokiaClient.ReadString(storageServicePath + "/PartitionerName")
	if err != nil {
		return "", err
	}
	return resp.Value, nil
}
