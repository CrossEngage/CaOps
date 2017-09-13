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
		return false, err
	}
	log.Printf("Going to do snapshot of %s.%s at %s", bp.KeyspaceGlob, bp.TableGlob,
		bp.TimeMarker.Format(time.RFC3339))
	<-time.After(bp.TimeMarker.Sub(time.Now()))
	result := caops.cassMngr.Snapshot(bp.KeyspaceGlob, bp.TableGlob)
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
