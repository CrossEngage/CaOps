package server

import (
	"log"
	"time"

	"github.com/hashicorp/serf/serf"
)

// TODO backup
// Wait until event time is reached
// Check cluster status
// Check remote storage connection
// Flush
// Check available disk space
// Get schema and upload to remote storage
// Trigger snapshotting
// Check amount of data of snapshots
// Upload to remote storage
// Cleanup snapshot

func (caops *CaOps) backupEventHandler(event serf.UserEvent) (breakLoop bool, err error) {
	bp, err := NewBackupPayload(event.Payload)
	if err != nil {
		log.Println(err)
		return false, err
	}
	log.Printf("Going to do snapshot of %s.%s at %s", bp.KeyspaceGlob, bp.Table, bp.TimeMarker.Format(time.RFC3339))

	keyspaces, err := caops.cassMngr.MatchKeyspaces(bp.KeyspaceGlob)
	if err != nil {
		log.Println(err)
		return false, err
	}

	<-time.After(bp.TimeMarker.Sub(time.Now()))

	if bp.Table == "" || bp.Table == "*" {
		tag, err := caops.cassMngr.SnapshotKeyspaces(keyspaces)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Snapshot of keyspaces (%#v) is done and tagged as %s ", keyspaces, tag)
	} else {
		for _, keyspace := range keyspaces {
			tag, err := caops.cassMngr.SnapshotTable(keyspace, bp.Table)
			if err != nil {
				log.Println(err)
			}
			log.Printf("Snapshot of %s.%s is done and tagged as %s ", keyspace, bp.Table, tag)
		}
	}

	return false, nil
}

func (caops *CaOps) clearSnapshotEventHandler(event serf.UserEvent) (breakLoop bool, err error) {
	log.Println("Clearing snapshots...")
	if err := caops.cassMngr.ClearSnapshot(); err != nil {
		log.Println(err)
	}
	return false, nil
}
