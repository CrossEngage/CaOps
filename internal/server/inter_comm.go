package server

import (
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/hashicorp/serf/serf"
)

// InterComm handles inter-agent communication
type InterComm struct {
	eventCh           chan serf.Event
	serf              *serf.Serf
	userEventHandlers map[string][]UserEventHandler
	shutdownCh        chan struct{}
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
		eventCh:    eventCh,
		serf:       serfCli,
		shutdownCh: make(chan struct{}), // TODO handle shutdowns
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

// AliveMembers return the IPs of all CaOps agents that are alive
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

		case e := <-ic.eventCh:
			switch ev := e.(type) {
			case serf.MemberEvent:
				log.Println("Event member", ev.EventType())
			case *serf.Query:
				log.Println("Event query", ev.EventType())
			case serf.UserEvent:
				log.Println("Event user", ev.String())
				ic.handleUserEvent(ev)
			default:
				log.Printf("Unknown event: %#v", e)
			}
		case <-ic.shutdownCh:
			return
		}
	}
}

func (ic *InterComm) handleUserEvent(event serf.UserEvent) {
	if !ic.isEventHandlerNameRegistered(event.Name) {
		return
	}
	total := len(ic.userEventHandlers[event.Name])
	for i, handler := range ic.userEventHandlers[event.Name] {
		log.Printf("Running handler %d of %d for event '%s'", i, total, event.Name)
		if stop, err := handler(event); err != nil {
			log.Printf("Error when running handler %d for event '%s': %s", i, event.String(), err)
		} else if stop && i < total-1 {
			log.Printf("'%s' event handler %d broke the handler loop, before %d handlers", event.Name, i, total-i-1)
			break
		}
	}
}

func (ic *InterComm) isEventHandlerNameRegistered(name string) bool {
	_, ok := ic.userEventHandlers[name]
	return ok
}

// RegisterEventHandler ...
func (ic *InterComm) RegisterEventHandler(name string, handler UserEventHandler) {
	if !ic.isEventHandlerNameRegistered(name) {
		ic.userEventHandlers[name] = make([]UserEventHandler, 0)
	}
	ic.userEventHandlers[name] = append(ic.userEventHandlers[name], handler)
}
