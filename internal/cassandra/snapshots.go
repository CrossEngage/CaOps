package cassandra

import (
	"fmt"
	"time"

	"github.com/gobwas/glob"
)

// SnapshotKeyspaces triggers a snapshot for the given list of keyspaces, and returns a generated tag
func (m *Manager) SnapshotKeyspaces(keyspaces []string) (tag string, err error) {
	tag = m.genSnapshotName()
	return tag, m.storageService.TakeSnapshot(tag, keyspaces...)
}

// SnapshotTable triggers a snapshot for the specified keyspace and table, and returns a generated tag
func (m *Manager) SnapshotTable(keyspace, table string) (tag string, err error) {
	tag = m.genSnapshotName()
	return tag, m.storageService.TakeTableSnapshot(tag, keyspace, table)
}

// SnapshotTables triggers a snapshot for the specified keyspace.table combinations, and returns a generated tag
func (m *Manager) SnapshotTables(tables []string) (tag string, err error) {
	tag = m.genSnapshotName()
	return tag, m.storageService.TakeMultipleTableSnapshot(tag, tables...)
}

// MatchKeyspaces returns a list of keyspace names that matches the glob
func (m *Manager) MatchKeyspaces(keyspaceGlob string) ([]string, error) {
	kg := glob.MustCompile(keyspaceGlob)
	allKeyspaces, err := m.storageService.Keyspaces()
	if err != nil {
		return nil, err
	}
	keyspaces := make([]string, 0, len(allKeyspaces))
	for _, keyspace := range allKeyspaces {
		if kg.Match(keyspace) {
			keyspaces = append(keyspaces, keyspace)
		}
	}
	return keyspaces, nil
}

func (m *Manager) genSnapshotName() string {
	return fmt.Sprintf("%s-CaOps", time.Now().Format("20060102T150405.000000"))
}

// ClearSnapshot is similar to nodetool clearsnapshot
func (m *Manager) ClearSnapshot() error {
	return m.storageService.ClearSnapshot("")
}
