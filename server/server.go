package server

import (
	"log"

	"bitbucket.org/crossengage/athena/server/agent"
	"bitbucket.org/crossengage/athena/server/api"
	logging "github.com/op/go-logging"
)

// Athena encapsulates all the Athena server behaviour
type Athena struct {
	log       *logging.Logger
	agent     *agent.Agent
	apiServer *api.Server
}

// NewAthena constructs a new Athena server
func NewAthena(log *logging.Logger, apiServerBindAddr, gossipBindAddr, gossipSnapshotPath, jolokiaAddr string) (*Athena, error) {
	agent, err := agent.NewAgent(log, gossipBindAddr, gossipSnapshotPath, jolokiaAddr)
	if err != nil {
		return nil, err
	}
	apiServer := api.NewServer(log, apiServerBindAddr)
	return &Athena{log: log, agent: agent, apiServer: apiServer}, nil
}

// Run starts the agent and the HTTP API server, and blocks, until it is finished
func (athena *Athena) Run() {
	if err := athena.agent.Start(); err != nil {
		log.Fatal(err)
	}
	if err := athena.apiServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	// for {
	// 	select {
	// 	case ev := <-eventCh:
	// 		log.Debugf("Event Type: %s, Event: %s\n", ev.EventType().String(), ev.String())
	// 	// case <-ticker.C:
	// 	// log.Debugf("Num Nodes: %d, Members: %+v\n", serfCli.NumNodes(), serfCli.Members())
	// 	case <-tickerNewSched.C:
	// 		serfCli.UserEvent(fmt.Sprintf("Evento:%d:%s", counter, serfCli.LocalMember().Name), nil, true)
	// 		counter++
	// 	}
	// }
}
