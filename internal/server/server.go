package server

import (
	"log"

	"github.com/CrossEngage/CaOps/internal/agent"
)

// CaOps encapsulates all the CaOps server behaviour
type CaOps struct {
	agent       *agent.Agent
	httpService *HTTPService
}

// NewCaOps constructs a new CaOps server
func NewCaOps(httpServiceBindAddr, gossipBindAddr, gossipSnapshotPath, jolokiaAddr string) (*CaOps, error) {
	agent, err := agent.NewAgent(gossipBindAddr, gossipSnapshotPath, jolokiaAddr)
	if err != nil {
		return nil, err
	}
	httpService := NewHTTPService(httpServiceBindAddr, agent)
	return &CaOps{agent: agent, httpService: httpService}, nil
}

// Run starts the agent and the HTTP API server, and blocks, until it is finished
func (CaOps *CaOps) Run() {
	if err := CaOps.agent.Start(); err != nil {
		log.Fatal(err)
	}
	if err := CaOps.httpService.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
