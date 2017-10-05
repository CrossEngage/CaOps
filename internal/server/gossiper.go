package server

import (
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/serf/serf"
)

// EventHandler receives the event name, a payload (up to 512 bytes), and must return
// if the event handler processing must break or not, and if necessary, return an error
type EventHandler func(event serf.UserEvent) (breakLoop bool, err error)

// EventHandlersMap is a map of event IDs to its collection of handlers
type EventHandlersMap map[string][]EventHandler

// Gossiper handles CaOps cluster-wide communication. It is used to send cluster-wide commands,
type Gossiper struct {
	eventCh          chan serf.Event
	serf             *serf.Serf
	eventHandlers    EventHandlersMap
	eventHandlersMtx sync.Mutex
	shutdownCh       chan struct{}
}

// NewGossiper constructs a new Gossiper object
func NewGossiper(bindTo, snapshotPath string) (*Gossiper, error) {
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

	gossiper := &Gossiper{
		eventCh:       eventCh,
		serf:          serfCli,
		eventHandlers: make(EventHandlersMap),
		shutdownCh:    make(chan struct{}), // TODO handle shutdowns
	}
	return gossiper, nil
}

// Join the cluster formed by nodes
func (g *Gossiper) Join(nodes []string) error {
	logrus.Debug("Joining gossiper to nodes: ", nodes)
	contacted, err := g.serf.Join(nodes, false)
	if err != nil {
		logrus.Errorf("Contacted %d nodes, but %s", contacted, err)
		return err
	}
	logrus.Infof("Contacted %d nodes", contacted)
	return nil
}

// AliveMembers return the IPs of all CaOps agents that are alive
func (g *Gossiper) AliveMembers() []string {
	ips := make([]string, 0)
	for _, member := range g.serf.Members() {
		if member.Status == serf.StatusAlive {
			ips = append(ips, member.Addr.String())
		}
	}
	return ips
}

// EventLoop watches for events and calls the proper triggers
func (g *Gossiper) EventLoop() {
	for {
		select {

		case e := <-g.eventCh:
			switch ev := e.(type) {
			case serf.MemberEvent:
				logrus.Debug("[84] Event member", ev.EventType())
			case *serf.Query:
				logrus.Debug("[86] Event query", ev.EventType())
			case serf.UserEvent:
				logrus.Debug("[88] Event user", ev.String())
				g.handleUserEvent(ev)
			default:
				logrus.Debugf("[91] Unknown event: %#v", e)
			}
		case <-g.shutdownCh:
			return
		}
	}
}

func (g *Gossiper) handleUserEvent(event serf.UserEvent) {
	g.eventHandlersMtx.Lock() // to avoid data races
	defer g.eventHandlersMtx.Unlock()
	if !g.isEventHandlerNameRegistered(event.Name) {
		logrus.Errorf("Unknown event type '%s'", event.Name)
		return
	}
	total := len(g.eventHandlers[event.Name])
	if total == 0 {
		logrus.Warnf("No event handlers for '%s'", event.Name)
	}
	for i, handler := range g.eventHandlers[event.Name] {
		logrus.Debugf("Running handler %d of %d for event '%s'", i+1, total, event.Name)
		if stop, err := handler(event); err != nil {
			logrus.Errorf("Error when running handler %d for event '%s': %s", i, event.String(), err)
		} else if stop && i < total-1 {
			logrus.Warnf("'%s' event handler %d broke the handler loop, before %d handlers", event.Name, i, total-i-1)
			break
		}
	}
}

// RegisterEventHandler adds a new event handler
func (g *Gossiper) RegisterEventHandler(name string, handler EventHandler) {
	g.eventHandlersMtx.Lock()
	defer g.eventHandlersMtx.Unlock()
	if !g.isEventHandlerNameRegistered(name) {
		g.eventHandlers[name] = make([]EventHandler, 0)
	}
	g.eventHandlers[name] = append(g.eventHandlers[name], handler)
}

// SendEvent sends user events
func (g *Gossiper) SendEvent(name string, payload EventPayload) error {
	return g.serf.UserEvent(name, payload.Encode(), true)
}

func (g *Gossiper) isEventHandlerNameRegistered(name string) bool {
	_, ok := g.eventHandlers[name]
	return ok
}

// EventPayload ...
type EventPayload interface {
	Encode() []byte
}
