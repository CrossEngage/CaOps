package agent

import (
	"fmt"

	"github.com/crossengage/CaOps/server/cassandra"
	"github.com/crossengage/CaOps/server/gossip"
)

// Agent orchestrates all operations
type Agent struct {
	cassMngr  *cassandra.Manager
	interComm *gossip.InterComm
}

// NewAgent ...
func NewAgent(gossipBindAddr, gossipSnapshotPagithub.com/crossengageth, jolokiaAddr string) (*Agent, error) {
	// Create the Cassandra Manager
	cassMngr, err := cassandra.NewManager(jolokiaAddr)
	if err != nil {
		return nil, err
	}
	// Create the Gossip InterComm
	interComm, err := gossip.NewInterComm(gossipBindAddr, gossipSnapshotPath)
	if err != nil {
		return nil, err
	}
	// Give to the poor
	return &Agent{
		cassMngr:  cassMngr,
		interComm: interComm,
	}, nil
}

// Start initializes the gossip, check cluster status, and triggers the event loop
func (ag *Agent) Start() error {
	if err := ag.cassMngr.CheckClusterStability(); err != nil {
		return err
	}
	liveNodes, err := ag.cassMngr.LiveNodes()
	if err != nil {
		return err
	}
	if err := ag.interComm.Join(liveNodes); err != nil {
		return err
	}
	liveNodesMap := stringListToMapKeys(liveNodes)
	for _, ip := range ag.interComm.AliveMembers() {
		if _, ok := liveNodesMap[ip]; !ok {
			return fmt.Errorf("Cassandra node %s has not joined this CaOps cluster", ip)
		}
	}
	// TODO watch for cancelation
	go ag.interComm.EventLoop()
	return nil
}

// TODO extract cluster check logic to a method

// DoSnapshot triggers a snapshot on Cassandra for the given keyspace and table.
// The keyspace parameter supports glob expansion, but the table parameter only suports
// all (*) or a specific table.
func (ag *Agent) DoSnapshot(keyspace, table string) error {
	results := ag.cassMngr.Snapshot(keyspace, table)
	// TODO instead of snapshotting directly here, add an event handler to listen for snapshot tasks
	return results
}
