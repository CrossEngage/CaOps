package server

import (
	"log"

	"github.com/CrossEngage/CaOps/server/agent"
	"github.com/CrossEngage/CaOps/server/api"
)

// CaOps encapsulates all the CaOps server behaviour
type CaOps struct {
	agent     *agent.Agent
	apiServer *api.Server
}

// NewCaOps constructs a new CaOps server
func NewCaOps(apiServerBindAddr, gossipBindAddr, gossipSnapshotPath, jolokiaAddr string) (*CaOps, error) {
	agent, err := agent.NewAgent(gossipBindAddr, gossipSnapshotPath, jolokiaAddr)
	if err != nil {
		return nil, err
	}
	apiServer := api.NewServer(apiServerBindAddr, agent)
	return &CaOps{agent: agent, apiServer: apiServer}, nil
}

// Run starts the agent and the HTTP API server, and blocks, until it is finished
func (CaOps *CaOps) Run() {
	if err := CaOps.agent.Start(); err != nil {
		log.Fatal(err)
	}
	if err := CaOps.apiServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
