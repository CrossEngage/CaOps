package server

import (
	"log"
	"strings"

	"github.com/hashicorp/serf/serf"
)

func (caops *CaOps) backupEventHandler(event serf.UserEvent) (breakLoop bool, err error) {
	pl := strings.Split(string(event.Payload), ":")
	keyspace := pl[0]
	table := pl[1]
	log.Printf("Triggered snapshot of %#v:%#v", keyspace, table)
	result := caops.cassMngr.Snapshot(keyspace, table)
	log.Println(result.String())
	return false, nil
}

func (caops *CaOps) clearSnapshotEventHandler(event serf.UserEvent) (breakLoop bool, err error) {
	log.Println("Clearing snapshots...")
	if err := caops.cassMngr.ClearSnapshot(); err != nil {
		log.Println(err)
	}
	return false, nil
}
