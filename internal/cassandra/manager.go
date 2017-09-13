package cassandra

import (
	"net/http"
	"net/url"

	"github.com/CrossEngage/CaOps/internal/jolokia"
)

// Manager handles all interaction with a Cassandra node and cluster
type Manager struct {
	jolokiaClient  jolokia.Client
	storageService storageService
}

// NewManager builds a new Cassandra Manager to encapsulate all interaction with a Cassandra node
func NewManager(jolokiaAddr string) (*Manager, error) {
	jolokiaURL, err := url.Parse(jolokiaAddr)
	if err != nil {
		return nil, err
	}
	jolokiaClient := jolokia.NewClient(*http.DefaultClient, *jolokiaURL)
	manager := &Manager{
		storageService: storageService{jolokiaClient},
		jolokiaClient:  jolokiaClient,
	}
	return manager, nil
}

// CheckClusterStability checks if the cluster is stable, or if it have no
// unreachable, joining, leaving, or moving nodes
func (m *Manager) CheckClusterStability() error {
	if err := mustGetEmptyStringList(m.storageService.UnreachableNodes,
		ErrUnreachableCassandraNodes); err != nil {
		return err
	}
	if err := mustGetEmptyStringList(m.storageService.JoiningNodes,
		ErrJoiningCassandraNodes); err != nil {
		return err
	}
	if err := mustGetEmptyStringList(m.storageService.LeavingNodes,
		ErrLeavingCassandraNodes); err != nil {
		return err
	}
	if err := mustGetEmptyStringList(m.storageService.MovingNodes,
		ErrMovingCassandraNodes); err != nil {
		return err
	}
	return nil
}

// LiveNodes return a list of IPs of the Cassandra nodes with status=LIVE
func (m *Manager) LiveNodes() ([]string, error) {
	return m.storageService.LiveNodes()
}

// JolokiaAgentVersion returns the version of the Jolokia Agent running
// into the Cassandra's JVM
func (m *Manager) JolokiaAgentVersion() (string, error) {
	verResp, err := m.jolokiaClient.Version()
	if err != nil {
		return "", err
	}
	return verResp.Value.Agent, nil
}

// CassandraVersion returns the version of the Cassandra node this manager
// is connected to
func (m *Manager) CassandraVersion() (string, error) {
	cver, err := m.storageService.ReleaseVersion()
	if err != nil {
		return "", err
	}
	return cver, nil
}

// UnreachableNodes retrieve the list of unreachable nodes in the cluster, as
// determined by this node's failure detector.
func (m *Manager) UnreachableNodes() (ips []string, err error) {
	return m.storageService.UnreachableNodes()
}

// JoiningNodes retrieve the list of nodes currently bootstrapping into the ring.
func (m *Manager) JoiningNodes() (ips []string, err error) {
	return m.storageService.JoiningNodes()
}

// LeavingNodes retrieve the list of nodes currently leaving the ring.
func (m *Manager) LeavingNodes() (ips []string, err error) {
	return m.storageService.LeavingNodes()
}

// MovingNodes retrieve the list of nodes currently moving in the ring.
func (m *Manager) MovingNodes() (ips []string, err error) {
	return m.storageService.MovingNodes()
}

// Tokens fetch string representations of the tokens for this node.
func (m *Manager) Tokens() ([]string, error) {
	return m.storageService.Tokens()
}

// SchemaVersion fetch a string representation of the current Schema version.
func (m *Manager) SchemaVersion() (version string, err error) {
	return m.storageService.SchemaVersion()
}

// AllDataFileLocations returns the list of all data file locations from conf
func (m *Manager) AllDataFileLocations() (paths []string, err error) {
	return m.storageService.AllDataFileLocations()
}

// CommitLogLocation returns the location of the commit log
func (m *Manager) CommitLogLocation() (string, error) {
	return m.storageService.CommitLogLocation()
}

// SavedCachesLocation returns the location of the saved caches dir
func (m *Manager) SavedCachesLocation() (string, error) {
	return m.storageService.SavedCachesLocation()
}

// TokenToEndpointMap retrieve a map of tokens to endpoints, including the bootstrapping ones.
func (m *Manager) TokenToEndpointMap() (map[string]string, error) {
	return m.storageService.TokenToEndpointMap()
}

// LocalHostID returns the hosts unique ID
func (m *Manager) LocalHostID() (string, error) {
	return m.storageService.LocalHostID()
}

// EndpointToHostID retrieve the mapping of endpoint to host ID
func (m *Manager) EndpointToHostID() (map[string]string, error) {
	return m.storageService.EndpointToHostID()
}

// HostIDToEndpoint retrieve the mapping of host ID to endpoint
func (m *Manager) HostIDToEndpoint() (map[string]string, error) {
	return m.storageService.HostIDToEndpoint()
}

// Starting returns whether the storage service is starting or not
func (m *Manager) Starting() (bool, error) {
	return m.storageService.Starting()
}

// GossipRunning returns whether the gossip is running
func (m *Manager) GossipRunning() (bool, error) {
	return m.storageService.GossipRunning()
}

// Keyspaces return the list of keyspaces in the cluster
func (m *Manager) Keyspaces() ([]string, error) {
	return m.storageService.Keyspaces()
}

// NonSystemKeyspaces ...
func (m *Manager) NonSystemKeyspaces() ([]string, error) {
	return m.storageService.NonSystemKeyspaces()
}

// OperationMode returns the operation mode of the node. STARTING, NORMAL, JOINING,
// LEAVING, DECOMMISSIONED, MOVING, DRAINING, DRAINED.
func (m *Manager) OperationMode() (string, error) {
	return m.storageService.OperationMode()
}

// IncrementalBackupsEnabled is self explanatory
func (m *Manager) IncrementalBackupsEnabled() (bool, error) {
	return m.storageService.IncrementalBackupsEnabled()
}

// Initialized is self explanatory
func (m *Manager) Initialized() (bool, error) {
	return m.storageService.Initialized()
}

// Joined is self explanatory
func (m *Manager) Joined() (bool, error) {
	return m.storageService.Joined()
}

// ClusterName returns the name of the cluster
func (m *Manager) ClusterName() (name string, err error) {
	return m.storageService.ClusterName()
}

// PartitionerName returns the cluster partitioner
func (m *Manager) PartitionerName() (name string, err error) {
	return m.storageService.PartitionerName()
}
