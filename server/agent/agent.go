package agent

import (
	"fmt"

	"bitbucket.org/crossengage/athena/server/cassandra"
	"bitbucket.org/crossengage/athena/server/gossip"
)

// Agent orchestrates all operations
type Agent struct {
	cassMngr  *cassandra.Manager
	interComm *gossip.InterComm
}

// NewAgent ...
func NewAgent(gossipBindAddr, gossipSnapshotPath, jolokiaAddr string) (*Agent, error) {
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
			return fmt.Errorf("Cassandra node %s has not joined this Athena cluster", ip)
		}
	}
	// TODO watch for cancelation
	go ag.interComm.EventLoop()
	return nil
}

func stringListToMapKeys(list []string) map[string]bool {
	ret := make(map[string]bool)
	for _, item := range list {
		ret[item] = true
	}
	return ret
}
