package cmd

import "github.com/hashicorp/serf/serf"
import "github.com/spf13/viper"
import "net"

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
