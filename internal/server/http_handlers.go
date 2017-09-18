package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// BackupPayload ...
type BackupPayload struct {
	KeyspaceGlob string
	Table        string
	TimeMarker   time.Time
}

// NewBackupPayload ...
func NewBackupPayload(payload []byte) (*BackupPayload, error) {
	parts := strings.Split(string(payload), "\x1F")
	p := &BackupPayload{}
	p.KeyspaceGlob = parts[0]
	p.Table = parts[1]
	timeMarker, err := time.Parse(time.RFC3339, parts[2])
	if err != nil {
		return nil, err
	}
	p.TimeMarker = timeMarker
	return p, nil
}

// Encode ...
func (p *BackupPayload) Encode() []byte {
	str := strings.Join([]string{p.KeyspaceGlob, p.Table, p.TimeMarker.Format(time.RFC3339)}, "\x1F")
	return []byte(str)
}

func (caops *CaOps) backupHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	defer r.Body.Close()

	keyspaceGlob, table := "*", "*"
	if val, ok := vars["keyspaceGlob"]; ok {
		keyspaceGlob = val
	}
	if val, ok := vars["table"]; ok {
		table = val
	}

	timeMarker, err := caops.backup(keyspaceGlob, table)
	if err != nil {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "Snapshot of %s.%s at %s was requested", keyspaceGlob, table, timeMarker.Format(time.RFC3339))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error while triggering snapshot: %s", err)
	}
}

func (caops *CaOps) backup(keyspaceGlob, table string) (timeMarker time.Time, err error) {
	// TODO make this time configurable or based on some existing metric (some soft of cluster thrift)
	timeMarker = getNextRoundedTimeWithin(time.Now(), 15*time.Second)
	log.Printf("Backup of %s.%s requested for %s", keyspaceGlob, table, timeMarker.Format(time.RFC3339))
	payload := &BackupPayload{KeyspaceGlob: keyspaceGlob, Table: table, TimeMarker: timeMarker}
	err = caops.gossiper.SendEvent("backup", payload)
	return
}

// EmptyPayload ...
type EmptyPayload struct {
}

// Encode ...
func (p *EmptyPayload) Encode() []byte {
	return []byte{}
}

func (caops *CaOps) clearSnapshotHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if caops.gossiper.SendEvent("clearsnapshot", &EmptyPayload{}) != nil {
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (caops *CaOps) statusHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	jolokiaVersion, err := caops.cassMngr.JolokiaAgentVersion()
	fmt.Fprintln(w, "Jolokia Agent Version:", jolokiaVersion, getErrStr(err))

	version, err := caops.cassMngr.CassandraVersion()
	fmt.Fprintln(w, "C* Version:", version, getErrStr(err))
	schemaVersion, err := caops.cassMngr.SchemaVersion()
	fmt.Fprintln(w, "C* Schema Version:", schemaVersion, getErrStr(err))
	dataFileLocs, err := caops.cassMngr.AllDataFileLocations()
	fmt.Fprintln(w, "C* Data File Locs:", strings.Join(dataFileLocs, ", "), getErrStr(err))
	commitLogLoc, err := caops.cassMngr.CommitLogLocation()
	fmt.Fprintln(w, "C* CommitLog Loc:", commitLogLoc, getErrStr(err))
	savedCacheLoc, err := caops.cassMngr.SavedCachesLocation()
	fmt.Fprintln(w, "C* Saved Caches Loc:", savedCacheLoc, getErrStr(err))
	locHostID, err := caops.cassMngr.LocalHostID()
	fmt.Fprintln(w, "C* Local Host ID:", locHostID, getErrStr(err))
	partitionerName, err := caops.cassMngr.PartitionerName()
	fmt.Fprintln(w, "C* Partitoner Name:", partitionerName, getErrStr(err))
	opMode, err := caops.cassMngr.OperationMode()
	fmt.Fprintln(w, "C* Operation Mode:", opMode, getErrStr(err))
	incrBkpEnabled, err := caops.cassMngr.IncrementalBackupsEnabled()
	fmt.Fprintln(w, "C* Incremental Backups Enabled:", incrBkpEnabled, getErrStr(err))

	clusterName, err := caops.cassMngr.ClusterName()
	fmt.Fprintln(w, "C* Cluster Name:", clusterName, getErrStr(err))
	liveNodes, err := caops.cassMngr.LiveNodes()
	fmt.Fprintln(w, "C* Cluster Live Nodes:", strings.Join(liveNodes, ", "), getErrStr(err))
	joiningNodes, err := caops.cassMngr.JoiningNodes()
	fmt.Fprintln(w, "C* Cluster Joining Nodes:", strings.Join(joiningNodes, ", "), getErrStr(err))
	leavingNodes, err := caops.cassMngr.LeavingNodes()
	fmt.Fprintln(w, "C* Cluster Leaving Nodes:", strings.Join(leavingNodes, ", "), getErrStr(err))
	movingNodes, err := caops.cassMngr.MovingNodes()
	fmt.Fprintln(w, "C* Cluster Moving Nodes:", strings.Join(movingNodes, ", "), getErrStr(err))
	keyspaces, err := caops.cassMngr.Keyspaces()
	fmt.Fprintln(w, "C* Keyspaces:", strings.Join(keyspaces, ", "), getErrStr(err))
	nonSysKeyspaces, err := caops.cassMngr.NonSystemKeyspaces()
	fmt.Fprintln(w, "C* Non-System Keyspaces:", strings.Join(nonSysKeyspaces, ", "), getErrStr(err))
}
