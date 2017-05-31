package gossip

import (
	"log"
	"net"
	"path/filepath"

	"os"

	"github.com/hashicorp/serf/serf"
)

// InterComm handles inter-agent communication
type InterComm struct {
	eventCh chan serf.Event
	serf    *serf.Serf
}

// NewInterComm constructs a new InterComm object that handles inter-agent communication
func NewInterComm(bindTo, snapshotPath string) (*InterComm, error) {
	serfBindAddr, err := net.ResolveTCPAddr("tcp", bindTo)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(snapshotPath), os.FileMode(0770)); err != nil {
		return nil, err
	}

	eventCh := make(chan serf.Event, 256)
	config := serf.DefaultConfig()
	config.Init()
	config.MemberlistConfig.BindAddr = serfBindAddr.IP.String()
	config.MemberlistConfig.BindPort = serfBindAddr.Port
	config.EventCh = eventCh
	config.SnapshotPath = snapshotPath
	serfCli, err := serf.Create(config)
	if err != nil {
		return nil, err
	}

	interComm := &InterComm{
		eventCh: eventCh,
		serf:    serfCli,
	}
	return interComm, nil
}

// Join ...
func (ic *InterComm) Join(nodes []string) error {
	contacted, err := ic.serf.Join(nodes, false)
	if err != nil {
		log.Printf("Contacted %d nodes, but %s", contacted, err)
		return err
	}
	log.Printf("Contacted %d nodes", contacted)
	return nil
}

// AliveMembers return the IPs of all Athena agents that are alive
func (ic *InterComm) AliveMembers() []string {
	ips := make([]string, 0)
	for _, member := range ic.serf.Members() {
		if member.Status == serf.StatusAlive {
			ips = append(ips, member.Addr.String())
		}
	}
	return ips
}

// EventLoop watches for events and calls the proper triggers
func (ic *InterComm) EventLoop() {
	for {
		select {
		case ev := <-ic.eventCh:
			log.Printf("Event Type: %s, Event: %s\n", ev.EventType().String(), ev.String())
			// case <-ticker.C:
			// log.Debugf("Num Nodes: %d, Members: %+v\n", serfCli.NumNodes(), serfCli.Members())
		}
	}
}
