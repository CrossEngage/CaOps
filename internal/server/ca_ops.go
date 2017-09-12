package server

import (
	"fmt"
	"log"
	"time"

	"github.com/CrossEngage/CaOps/internal/cassandra"
)

// CaOps encapsulates all the CaOps server behavior
type CaOps struct {
	httpService *HTTPService
	cassMngr    *cassandra.Manager
	gossiper    *Gossiper
}

// NewCaOps constructs a new CaOps server
func NewCaOps(httpBindAddr, gossipBindAddr, gossipSnapshotPath, jolokiaAddr string) (*CaOps, error) {

	// Create the Cassandra Manager
	cassMngr, err := cassandra.NewManager(jolokiaAddr)
	if err != nil {
		return nil, err
	}

	// Create the Gossiper
	gossiper, err := NewGossiper(gossipBindAddr, gossipSnapshotPath)
	if err != nil {
		return nil, err
	}

	// Create the HTTP Service
	httpService := NewHTTPService(httpBindAddr, gossiper)

	return &CaOps{
		httpService: httpService,
		cassMngr:    cassMngr,
		gossiper:    gossiper,
	}, nil
}

// Run starts the agent and the HTTP API server, and blocks, until it is finished
func (caops *CaOps) Run() {
	for {
		if err := caops.Init(); err != nil {
			log.Println(err)
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	if err := caops.httpService.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// Init starts gossiper, check cluster status, and triggers the event loop
func (caops *CaOps) Init() error {
	if err := caops.cassMngr.CheckClusterStability(); err != nil {
		return err
	}
	liveNodes, err := caops.cassMngr.LiveNodes()
	if err != nil {
		return err
	}
	if err := caops.gossiper.Join(liveNodes); err != nil {
		return err
	}
	liveNodesMap := stringListToMapKeys(liveNodes)
	for _, ip := range caops.gossiper.AliveMembers() {
		if _, ok := liveNodesMap[ip]; !ok {
			return fmt.Errorf("Cassandra node %s has not joined this CaOps cluster", ip)
		}
	}
	go caops.gossiper.EventLoop()
	return nil
}

// CheckClustersConsistency compares the live nodes of Cassandra with the live nodes of CaOps to
// determine if both are consistent with each other, and returns error if not.
func (caops *CaOps) CheckClustersConsistency() error {
	liveNodes, err := caops.cassMngr.LiveNodes()
	if err != nil {
		return err
	}
	liveNodesMap := stringListToMapKeys(liveNodes)
	for _, ip := range caops.gossiper.AliveMembers() {
		if _, ok := liveNodesMap[ip]; !ok {
			return fmt.Errorf("Cassandra node %s has not joined the CaOps cluster", ip)
		}
	}
	return nil
}
