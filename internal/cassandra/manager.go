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

// NewManager builds a new Cassandra Manager to encapsulate all interaction with a Cassandra node and Cluster
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
