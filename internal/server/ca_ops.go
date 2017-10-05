package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/CrossEngage/CaOps/internal/cassandra"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

// CaOps encapsulates all the CaOps server behavior
type CaOps struct {
	cassMngr *cassandra.Manager
	gossiper *Gossiper
	stopChan chan os.Signal
	server   *http.Server
	router   *mux.Router
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

	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	router := mux.NewRouter()

	caops := &CaOps{
		stopChan: stopChan,
		server:   &http.Server{Addr: httpBindAddr, Handler: router},
		router:   router,
		cassMngr: cassMngr,
		gossiper: gossiper,
	}

	router.Methods("GET").
		Path("/status").
		HandlerFunc(caops.statusHandler)
	router.Methods("GET").
		Path("/backup-keyspaces/{keyspaceGlob}").
		HandlerFunc(caops.backupHandler)
	router.Methods("GET").
		Path("/backup-tables/{keyspaceGlob}/{table}").
		HandlerFunc(caops.backupHandler)
	router.Methods("DELETE").
		Path("/snapshots").
		HandlerFunc(caops.clearSnapshotHandler)

	return caops, nil
}

func (caops *CaOps) waitForShutdown() {
	<-caops.stopChan
	logrus.Info("Shutting down HTTP server...")
	// shut down gracefully, but wait no longer than 5 seconds before halting
	// TODO make this configurable - maybe increase it for when there are uploads happening
	ctx, cancelFun := context.WithTimeout(context.Background(), 5*time.Second)
	cancelFun()
	caops.server.Shutdown(ctx)
	logrus.Info("HTTP Server gracefully stopped")
}

// Run starts the agent and the HTTP API server, and blocks, until it is finished
func (caops *CaOps) Run() {
	for { // TODO add a timeout here
		if err := caops.Init(); err != nil {
			logrus.Error(err)
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	caops.gossiper.RegisterEventHandler("backup", caops.backupEventHandler)
	caops.gossiper.RegisterEventHandler("clearsnapshot", caops.clearSnapshotEventHandler)

	go caops.waitForShutdown()

	if err := caops.server.ListenAndServe(); err != nil {
		logrus.Fatal(err)
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
