package server

import (
	"log"

	"bitbucket.org/crossengage/athena/server/agent"
	"bitbucket.org/crossengage/athena/server/api"
)

// Athena encapsulates all the Athena server behaviour
type Athena struct {
	agent     *agent.Agent
	apiServer *api.Server
}

// NewAthena constructs a new Athena server
func NewAthena(apiServerBindAddr, gossipBindAddr, gossipSnapshotPath, jolokiaAddr string) (*Athena, error) {
	agent, err := agent.NewAgent(gossipBindAddr, gossipSnapshotPath, jolokiaAddr)
	if err != nil {
		return nil, err
	}
	apiServer := api.NewServer(apiServerBindAddr)
	return &Athena{agent: agent, apiServer: apiServer}, nil
}

// Run starts the agent and the HTTP API server, and blocks, until it is finished
func (athena *Athena) Run() {
	if err := athena.agent.Start(); err != nil {
		log.Fatal(err)
	}
	if err := athena.apiServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
