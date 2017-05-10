package cmd

import (
	"fmt"
	"time"

	"bitbucket.org/crossengage/athena/cassandra"
	"github.com/hashicorp/serf/serf"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the main daemon",
	Run:   runServeCmd,
}

func init() {
	RootCmd.AddCommand(serveCmd)
}

var serfCli *serf.Serf

func runServeCmd(cmd *cobra.Command, args []string) {
	serfCli = setupSerf()
	log.Info(serfCli.Stats())

	nodeprobe := cassandra.NewNodeProbe(getJolokiaClient())
	livenodes, err := nodeprobe.StorageService.LiveNodes()
	if err != nil {
		log.Fatal(err)
	}

	contactedNodes, err := serfCli.Join(livenodes, false)
	if err != nil {
		log.Errorf("Contacted %d nodes, but %s", contactedNodes, err)
	} else {
		log.Infof("Contacted %d nodes", contactedNodes)
	}

	// ticker := time.NewTicker(1 * time.Second)
	tickerNewSched := time.NewTicker(30 * time.Second)

	var counter int

	go httpServer()

	for {
		select {
		case ev := <-eventCh:
			log.Debugf("Event Type: %s, Event: %s\n", ev.EventType().String(), ev.String())
		// case <-ticker.C:
		// log.Debugf("Num Nodes: %d, Members: %+v\n", serfCli.NumNodes(), serfCli.Members())
		case <-tickerNewSched.C:
			serfCli.UserEvent(fmt.Sprintf("Evento:%d:%s", counter, serfCli.LocalMember().Name), nil, true)
			counter++
		}
	}
}
