package gossip

import (
	"net"

	"github.com/hashicorp/serf/serf"
	logging "github.com/op/go-logging"
)

// InterComm handles inter-agent communication
type InterComm struct {
	log     *logging.Logger
	eventCh chan serf.Event
	serf    *serf.Serf
}

// NewInterComm constructs a new InterComm object that handles inter-agent communication
func NewInterComm(log *logging.Logger, bindTo, snapshotPath string) (*InterComm, error) {
	serfBindAddr, err := net.ResolveTCPAddr("tcp", bindTo)
	if err != nil {
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
		log:     log,
		eventCh: eventCh,
		serf:    serfCli,
	}
	return interComm, nil
}

// Join ...
func (ic *InterComm) Join(nodes []string) error {
	contacted, err := ic.serf.Join(nodes, false)
	if err != nil {
		ic.log.Errorf("Contacted %d nodes, but %s", contacted, err)
		return err
	}
	ic.log.Debugf("Contacted %d nodes", contacted)
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

// func checkClusterStatus() error {

// 	livenodes, err := nodeprobe.StorageService.LiveNodes()
// 	if err != nil {
// 		return err
// 	}

// 	liveNodesMap := stringListToMapKeys(livenodes)
// 	for _, member := range serfCli.Members() {
// 		ip := member.Addr.String()
// 		if _, ok := liveNodesMap[ip]; !ok {
// 			return fmt.Errorf("Cassandra node %s has not joined this Athena cluster", ip)
// 		}
// 		if member.Status != serf.StatusAlive {
// 			return fmt.Errorf("The Athena node %s is not alive", ip)
// 		}
// 	}

// 	return nil
// }

// func stringListToMapKeys(list []string) map[string]bool {
// 	ret := make(map[string]bool)
// 	for _, item := range list {
// 		ret[item] = true
// 	}
// 	return ret
// }
