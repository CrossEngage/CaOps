package cassandra

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gobwas/glob"
)

// Snapshot takes snapshots of specific keyspaces that matches the keyspace glob
// and all tables, or any specific table that matches
func (m *Manager) Snapshot(keyspaceGlob, tableName string) SnapshotResult {
	if strings.TrimSpace(keyspaceGlob) == "" {
		return newSnapshotResultFor(ErrRequiredKeyspaceOrAsterisk)
	}
	if strings.TrimSpace(tableName) == "" {
		return newSnapshotResultFor(ErrRequiredTableOrAsterisk)
	}
	keyspaces, err := m.matchKeyspaces(keyspaceGlob)
	if err != nil {
		return newSnapshotResultFor(err)
	}
	name := m.genSnapshotName(keyspaceGlob, tableName)
	results := newSnapshotResult(name)
	for _, keyspace := range keyspaces {
		log.Printf("Doing snapshot %s of %s.%s", name, keyspace, tableName)
		if tableName == "*" {
			results.Set(keyspace, m.storageService.TakeSnapshot(name, keyspace))
		} else {
			results.Set(fmt.Sprintf("%s.%s", keyspace, tableName), m.storageService.TakeTableSnapshot(name, keyspace, tableName))
		}
	}
	return results
}

// SnapshotResult is a map between "keyspaceGlob.tableName" or
// just a keyspace name, and errors
type SnapshotResult struct {
	name    string
	results map[string]error
	single  error
}

func newSnapshotResultFor(err error) SnapshotResult {
	return SnapshotResult{single: err}
}

func newSnapshotResult(name string) SnapshotResult {
	return SnapshotResult{name: name, results: make(map[string]error)}
}

// Set a result error for a given keyspaceGlob.tableName
// The error can be nil.
func (sr *SnapshotResult) Set(key string, err error) {
	sr.results[key] = err
}

// HasError returns true if there are any errors
func (sr SnapshotResult) HasError() bool {
	if sr.single != nil {
		return true
	}
	for _, v := range sr.results {
		if v != nil {
			return true
		}
	}
	return false
}

// Error returns a single representation of all errors found
func (sr SnapshotResult) Error() string {
	allErrors := make([]string, 0)
	allErrors = append(allErrors, sr.single.Error())
	for k, v := range sr.results {
		allErrors = append(allErrors, fmt.Sprintf("Error while snapshotting %s: %s", k, v))
	}
	return strings.Join(allErrors, "; ")
}

// String returns an overview of the snapshot
func (sr SnapshotResult) String() string {
	all := make([]string, 0)
	all = append(all)
	for k := range sr.results {
		all = append(all, k)
	}
	return fmt.Sprintf("Snapshot '%s' of %s", sr.name, strings.Join(all, ", "))
}

// Name returns the snapshot name
func (sr SnapshotResult) Name() string {
	return sr.name
}

func (m *Manager) matchKeyspaces(keyspaceGlob string) ([]string, error) {
	kg := glob.MustCompile(keyspaceGlob)
	nonSysKeyspaces, err := m.storageService.NonSystemKeyspaces()
	if err != nil {
		return nil, err
	}
	keyspaces := make([]string, 0, len(nonSysKeyspaces))
	for _, keyspace := range nonSysKeyspaces {
		if kg.Match(keyspace) {
			keyspaces = append(keyspaces, keyspace)
		}
	}
	log.Printf("Keyspace glob matching of %s found: %#v", keyspaceGlob, keyspaces)
	return keyspaces, nil
}

func (m *Manager) genSnapshotName(keyspaceName, tableName string) string {
	return fmt.Sprintf(
		"athena:%s:%s:%s",
		strings.Replace(keyspaceName, "*", "", -1),
		strings.Replace(tableName, "*", "", -1),
		time.Now().Format(time.RFC3339))
}
