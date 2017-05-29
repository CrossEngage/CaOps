package cassandra

import (
	"net/http"
	"net/url"

	"bitbucket.org/crossengage/athena/jolokia"
)

// Manager handles all interaction with a Cassandra node and cluster
type Manager struct {
	storageService *storageService
}

// NewManager builds a new Cassandra Manager to encapsulate all interaction with a Cassandra node and Cluster
func NewManager(jolokiaAddr string) (*Manager, error) {
	jolokiaURL, err := url.Parse(jolokiaAddr)
	if err != nil {
		return nil, err
	}
	manager := &Manager{
		storageService: &storageService{jolokia.Client{HTTPClient: http.DefaultClient, BaseURL: *jolokiaURL}},
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

// LiveNodes return a list of IPs of the Cassandra nodes
// with status=LIVE
func (m *Manager) LiveNodes() ([]string, error) {
	return m.storageService.LiveNodes()
}
