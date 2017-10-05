package server

import (
	"time"

	"github.com/Sirupsen/logrus"
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
		logrus.Error(err)
		return false, err
	}
	logrus.Infof("Going to do snapshot of %s.%s at %s", bp.KeyspaceGlob, bp.Table, bp.TimeMarker.Format(time.RFC3339))

	keyspaces, err := caops.cassMngr.MatchKeyspaces(bp.KeyspaceGlob)
	if err != nil {
		logrus.Error(err)
		return false, err
	}

	<-time.After(bp.TimeMarker.Sub(time.Now()))

	if bp.Table == "" || bp.Table == "*" {
		_, tag, err := caops.cassMngr.SnapshotKeyspaces(keyspaces)
		if err != nil {
			logrus.Error(err)
			return false, err
		}
		logrus.Infof("Snapshot of keyspaces (%#v) is done and tagged as %s ", keyspaces, tag)
	} else {
		for _, keyspace := range keyspaces {
			tag, err := caops.cassMngr.SnapshotTable(keyspace, bp.Table)
			if err != nil {
				logrus.Error(err)
				return false, err
			}
			logrus.Infof("Snapshot of %s.%s is done and tagged as %s ", keyspace, bp.Table, tag)
		}
	}

	return false, nil
}

func (caops *CaOps) clearSnapshotEventHandler(event serf.UserEvent) (breakLoop bool, err error) {
	logrus.Info("Clearing snapshots...")
	if err := caops.cassMngr.ClearSnapshot(); err != nil {
		logrus.Error(err)
		return false, err
	}
	return false, nil
}
