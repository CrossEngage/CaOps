package agent

import (
	"bitbucket.org/crossengage/athena/server/cassandra"
	"bitbucket.org/crossengage/athena/server/gossip"
	logging "github.com/op/go-logging"
)

// Agent orchestrates all operations
type Agent struct {
	log       *logging.Logger
	cassMngr  *cassandra.Manager
	interComm *gossip.InterComm
}

// NewAgent ...
func NewAgent(log *logging.Logger, gossipBindAddr, gossipSnapshotPath, jolokiaAddr string) (*Agent, error) {
	// Create the Cassandra Manager
	cassMngr, err := cassandra.NewManager(jolokiaAddr)
	if err != nil {
		return nil, err
	}
	// Create the Gossip InterComm
	interComm, err := gossip.NewInterComm(log, gossipBindAddr, gossipSnapshotPath)
	if err != nil {
		return nil, err
	}
	// Give to the poor
	return &Agent{
		cassMngr:  cassMngr,
		interComm: interComm,
	}, nil
}
