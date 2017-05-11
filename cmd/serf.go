package cmd

import (
	"errors"
	"fmt"
	"net"

	"bitbucket.org/crossengage/athena/cassandra"
	"github.com/hashicorp/serf/serf"
	"github.com/spf13/viper"
)

var (
	eventCh chan serf.Event
)

func setupSerf() *serf.Serf {

	serfBindAddr, err := net.ResolveTCPAddr("tcp", viper.GetString("serf.bind"))
	if err != nil {
		log.Fatal(err)
	}

	eventCh = make(chan serf.Event, 256)
	config := serf.DefaultConfig()
	config.Init()
	config.MemberlistConfig.BindAddr = serfBindAddr.IP.String()
	config.MemberlistConfig.BindPort = serfBindAddr.Port
	config.EventCh = eventCh
	config.SnapshotPath = viper.GetString("serf.snapshot_path")
	log.Debugf("%+v\n", config)
	serfCli, err := serf.Create(config)
	if err != nil {
		log.Fatal(err)
	}
	return serfCli
}

func mustGetEmptyStringList(listErrFunc func() ([]string, error), nonEmptyListError error) error {
	list, err := listErrFunc()
	if err != nil {
		return err
	}
	if len(list) > 0 {
		return nonEmptyListError
	}
	return nil
}

var (
	errUnreachableCassandraNodes = errors.New("Unreachable Cassandra Nodes")
	errJoiningCassandraNodes     = errors.New("Joining Cassandra Nodes")
	errLeavingCassandraNodes     = errors.New("Leaving Cassandra Nodes")
	errMovingCassandraNodes      = errors.New("Moving Cassandra Nodes")
)

func stringListToMapKeys(list []string) map[string]bool {
	ret := make(map[string]bool)
	for _, item := range list {
		ret[item] = true
	}
	return ret
}

func checkClusterStatus() error {
	nodeprobe := cassandra.NewNodeProbe(getJolokiaClient())
	if err := mustGetEmptyStringList(nodeprobe.StorageService.UnreachableNodes, errUnreachableCassandraNodes); err != nil {
		return err
	}
	if err := mustGetEmptyStringList(nodeprobe.StorageService.JoiningNodes, errJoiningCassandraNodes); err != nil {
		return err
	}
	if err := mustGetEmptyStringList(nodeprobe.StorageService.LeavingNodes, errLeavingCassandraNodes); err != nil {
		return err
	}
	if err := mustGetEmptyStringList(nodeprobe.StorageService.MovingNodes, errMovingCassandraNodes); err != nil {
		return err
	}

	livenodes, err := nodeprobe.StorageService.LiveNodes()
	if err != nil {
		return err
	}

	liveNodesMap := stringListToMapKeys(livenodes)
	for _, member := range serfCli.Members() {
		ip := member.Addr.String()
		if _, ok := liveNodesMap[ip]; !ok {
			return fmt.Errorf("Cassandra node %s has not joined this Athena cluster", ip)
		}
		if member.Status != serf.StatusAlive {
			return fmt.Errorf("The Athena node %s is not alive", ip)
		}
	}

	return nil
}
