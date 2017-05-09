package cmd

import (
	"time"

	"github.com/gocql/gocql"
)

var (
	schema = []string{
		`CREATE KEYSPACE IF NOT EXISTS athena
		 WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };`,

		`CREATE TABLE IF NOT EXISTS athena.recurring_job (
			id       uuid,
			interval bigint,
			unit     varchar,
			at_time  varchar,
			kspace   varchar,	
			cftable  varchar,
			PRIMARY KEY (id),
		 ) WITH comment='Recurring schedule'`,

		`CREATE TABLE IF NOT EXISTS athena.one_time_job (
			id      uuid,
			kspace  varchar,
			cftable varchar,
			at      timeuuid,
			PRIMARY KEY ((id), at),
		 ) WITH CLUSTERING ORDER BY (at DESC) 
		   AND default_time_to_live=2592000
		   AND comment='Scheduled one time backups'`,
	}
)

func bootstrapDatabase(session *gocql.Session) error {
	for _, query := range schema {
		if err := session.Query(query).RetryPolicy(nil).Exec(); err != nil {
			log.Debug(query)
			return err
		}
	}
	return nil
}

type recurringJob struct {
	Interval uint64
	Unit     string
	AtTime   string
	Keyspace string
	Table    string
}

type oneTimeJob struct {
	Keyspace string
	Table    string
	At       time.Time
}

func getAllRecurringJobs(session *gocql.Session) (map[string]recurringJob, error) {
	return nil, nil
}
